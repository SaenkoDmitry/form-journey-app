import React, {useEffect, useState} from 'react';
import {useAuth} from '../context/AuthContext';
import SafeTextRenderer from "../components/SafeTextRenderer.tsx";
import {getExerciseGroups, getExerciseTypesByGroup} from "../api/exercises.ts";
import Button from "../components/Button.tsx";
import {ChevronDown, Loader} from "lucide-react";
import VideoPlayer from "../components/VideoPlayer.tsx";

const LibraryPage: React.FC = () => {
    const {user, loading: authLoading} = useAuth();

    const [groups, setGroups] = useState<Group[]>([]);
    const [selectedGroup, setSelectedGroup] = useState<string | null>(null);
    const [exercises, setExercises] = useState<ExerciseType[]>([]);
    const [loading, setLoading] = useState(false);
    const [openedId, setOpenedId] = useState<number | null>(null);

    // -------- load groups --------
    useEffect(() => {
        if (!user) return;

        getExerciseGroups().then((groups: Group[]) => {
            setGroups(groups);
            if (groups.length > 0) setSelectedGroup(groups[0].code);
        });
    }, [user]);

    // -------- load exercises --------
    useEffect(() => {
        if (!selectedGroup) return;

        setLoading(true);

        getExerciseTypesByGroup(selectedGroup).then((exerciseTypes: ExerciseType[]) => {
            setExercises(exerciseTypes);
            setLoading(false);
        })
    }, [selectedGroup]);

    return <div className={"page stack"}>
        <h1>Библиотека упражнений</h1>

        {/* ---------- GROUP TABS ---------- */}
        <div
            style={{
                display: 'grid',
                gridTemplateColumns: 'repeat(auto-fill, minmax(120px, 1fr))',
                gap: 12,          // ← больше расстояние
                padding: '4px 2px',
                marginBottom: 20
            }}
        >
            {groups.map(g => (
                <Button
                    variant={selectedGroup === g.code ? "active" : "ghost"}
                    key={g.code}
                    onClick={() => setSelectedGroup(g.code)}
                >
                    {g.name}
                </Button>
            ))}
        </div>

        {/* ---------- EXERCISES ---------- */
        }
        {loading && <Loader/>}

        <div style={{display: 'flex', flexDirection: 'column', gap: 12}}>
            {!loading && exercises.map((ex, index) => {
                const isOpen = openedId === ex.id;
                const softBg = index % 2 === 0 ? 'var(--color-card)' : 'var(--color-card-alt)'; // мягкое чередование

                return (
                    <div
                        key={ex.id}
                        style={{
                            border: '1px solid var(--color-border)',
                            borderRadius: 12,
                            padding: 12,
                            transition: 'all 0.2s ease',
                            boxShadow: isOpen ? '0 6px 12px rgba(0,0,0,0.15)' : '0 2px 4px rgba(0,0,0,0.05)',
                            backgroundColor: softBg,
                            cursor: 'pointer',
                        }}
                        onClick={() => setOpenedId(prev => (prev === ex.id ? null : ex.id))}
                    >
                        <div style={{display: 'flex', justifyContent: 'space-between', alignItems: 'center'}}>
                            <strong>{ex.name}</strong>
                            {/* стрелка */}
                            <span style={{
                                display: 'inline-block',
                                transition: 'transform 0.3s ease',
                                transform: isOpen ? 'rotate(90deg)' : 'rotate(0deg)',
                            }}><ChevronDown size={22}/></span>
                        </div>

                        {/* раскрытие с плавным эффектом */}
                        <div
                            style={{
                                maxHeight: isOpen ? 1000 : 0,
                                overflow: 'hidden',
                                transition: 'max-height 0.3s ease',
                                marginTop: isOpen ? 10 : 0,
                            }}
                        >
                            {ex.description && (
                                <p>
                                    <SafeTextRenderer html={ex.description}/>
                                </p>
                            )}

                            <div>
                                {ex.rest_in_seconds > 0 &&
                                    <div style={{marginBottom: 10}}><b>Отдых: </b>{ex.rest_in_seconds} секунд
                                    </div>}
                                <div><b>Единицы
                                    измерения:</b> {ex.units.split(',').map(field => unitTypes[field]).join(", ")}
                                </div>
                            </div>

                            {ex.url && (
                                <div
                                    onClick={(e) => e.stopPropagation()} // ← остановка всплытия клика
                                >
                                    <VideoPlayer url={ex.url}/>
                                </div>
                            )}
                        </div>
                    </div>
                );
            })}
        </div>
    </div>
        ;
};

export default LibraryPage;

const unitTypes = {
    'reps': 'Повторения',
    'weight': 'Вес',
    'minutes': 'Минуты',
    'meters': 'Метры',
};
