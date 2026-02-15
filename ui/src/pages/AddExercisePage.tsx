import React, { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { Loader, ArrowLeft } from "lucide-react";
import Button from "../components/Button";

import {
    getExerciseGroups,
    getExerciseTypesByGroup,
    addExercise,
} from "../api/exercises";

type Group = {
    code: string;
    name: string;
};

interface ExerciseType {
    id: number;
    name: string;
    description: string;
    accent: string;
}

const AddExercisePage = () => {
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();

    const [groups, setGroups] = useState<Group[]>([]);
    const [exerciseTypes, setExerciseTypes] = useState<ExerciseType[]>([]);
    const [selectedGroup, setSelectedGroup] = useState<Group | null>(null);

    const [loading, setLoading] = useState(true);
    const [adding, setAdding] = useState(false);

    // ================= LOAD GROUPS =================
    useEffect(() => {
        const loadGroups = async () => {
            try {
                const data = await getExerciseGroups();
                setGroups(data);
            } catch (e) {
                alert("Ошибка загрузки групп");
            } finally {
                setLoading(false);
            }
        };

        loadGroups();
    }, []);

    // ================= LOAD EXERCISES =================
    const handleSelectGroup = async (group: Group) => {
        setSelectedGroup(group);
        setLoading(true);

        try {
            const data = await getExerciseTypesByGroup(group.code);
            setExerciseTypes(data);
        } catch (e) {
            alert("Ошибка загрузки упражнений");
        } finally {
            setLoading(false);
        }
    };

    // ================= ADD EXERCISE =================
    const handleAddExercise = async (exerciseId: number) => {
        if (!id) return;

        try {
            setAdding(true);
            await addExercise(Number(id), exerciseId);

            // возвращаемся назад
            navigate(`/workouts/${id}`, { replace: true });
        } catch (e) {
            alert("Ошибка добавления упражнения");
        } finally {
            setAdding(false);
        }
    };

    // ================= UI =================
    if (loading) {
        return (
            <div className="page center">
                <Loader />
            </div>
        );
    }

    return (
        <div className="page stack">
            {/* HEADER */}
            <div style={{ display: "flex", alignItems: "center", gap: 10 }}>
                <Button variant="ghost" onClick={() => navigate(-1)}>
                    <ArrowLeft size={18} />
                </Button>
                <h2 style={{ margin: 0 }}>
                    {selectedGroup ? selectedGroup.name : "Выбор группы"}
                </h2>
            </div>

            {/* GROUP LIST */}
            {!selectedGroup && (
                <div className="stack">
                    {groups.map((group) => (
                        <div
                            key={group.code}
                            className="card"
                            onClick={() => handleSelectGroup(group)}
                            style={{
                                padding: "1rem",
                                borderRadius: 12,
                                border: "1px solid #eee",
                                cursor: "pointer",
                            }}
                        >
                            {group.name}
                        </div>
                    ))}
                </div>
            )}

            {/* EXERCISE LIST */}
            {selectedGroup && (
                <div className="stack">
                    {exerciseTypes.map((exercise) => (
                        <div
                            key={exercise.id}
                            className="card"
                            onClick={() =>
                                !adding &&
                                handleAddExercise(exercise.id)
                            }
                            style={{
                                padding: "1rem",
                                borderRadius: 12,
                                border: "1px solid #eee",
                                cursor: "pointer",
                                opacity: adding ? 0.6 : 1,
                            }}
                        >
                            <div style={{ fontWeight: 600 }}>
                                {exercise.name}
                            </div>

                            {/*{exercise.description && (*/}
                            {/*    <div*/}
                            {/*        style={{*/}
                            {/*            fontSize: 14,*/}
                            {/*            color: "#666",*/}
                            {/*            marginTop: 6,*/}
                            {/*        }}*/}
                            {/*    >*/}
                            {/*        {exercise.description}*/}
                            {/*    </div>*/}
                            {/*)}*/}
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
};

export default AddExercisePage;
