import React, {useEffect, useRef, useState} from "react";
import {useParams} from "react-router-dom";
import {getProgram} from "../api/days";
import {parsePreset, savePreset} from "../api/presets";
import {getExerciseGroups, getExerciseTypesByGroup} from "../api/exercises";
import Button from "../components/Button";
import "../styles/DayDetailsPage.css";
import {Loader} from "lucide-react";

import {closestCenter, DndContext} from "@dnd-kit/core";
import {arrayMove, SortableContext, useSortable, verticalListSortingStrategy,} from "@dnd-kit/sortable";
import {CSS} from "@dnd-kit/utilities";

// ================= SORTABLE ITEM =================
function SortableExercise({id, children}: any) {
    const {attributes, listeners, setNodeRef, transform, transition} =
        useSortable({id});

    const style = {
        transform: CSS.Transform.toString(transform),
        transition,
    };

    return (
        <div ref={setNodeRef} style={style}>
            {children({listeners, attributes})}
        </div>
    );
}

// ================= PAGE =================
export default function DayDetailsPage() {
    const {programId, dayId} = useParams();

    const [loading, setLoading] = useState(false);
    const [exercises, setExercises] = useState<any[]>([]);
    const [dayName, setDayName] = useState("");

    const [groups, setGroups] = useState<Group[]>([]);
    const [types, setTypes] = useState<any[]>([]);
    const [selectedGroup, setSelectedGroup] = useState("");
    const [selectedType, setSelectedType] = useState("");

    const [toast, setToast] = useState<string | null>(null);

    const autosaveTimer = useRef<any>(null);
    const firstLoad = useRef(true);

    // ================= LOAD =================
    const fetchDay = async () => {
        const program = await getProgram(Number(programId));
        const day = program.day_types.find((d: any) => d.id === Number(dayId));
        if (!day) return;

        setDayName(day.name);

        if (day.preset) {
            setLoading(true);
            const parsed = await parsePreset(day.preset).finally(() =>
                setLoading(false)
            );

            setExercises(
                parsed.exercises.map((ex: any, i: number) => ({
                    uid: `${ex.id}-${i}`,
                    id: ex.id,
                    name: ex.name,
                    sets: ex.sets,
                    units: ex.units,
                }))
            );
        }

        firstLoad.current = false;
    };

    useEffect(() => {
        fetchDay();
        getExerciseGroups().then(setGroups);
    }, []);

    // ================= AUTOSAVE =================
    useEffect(() => {
        if (firstLoad.current) return;

        clearTimeout(autosaveTimer.current);

        autosaveTimer.current = setTimeout(async () => {
            try {
                await savePreset(Number(dayId), buildPreset());
                showToast("💾 Сохранено");
            } catch {
                showToast("❌ Ошибка сохранения");
            }
        }, 500);
    }, [exercises]);

    const showToast = (text: string) => {
        setToast(text);
        setTimeout(() => setToast(null), 1500);
    };

    // ================= DRAG END =================
    const handleDragEnd = (event: any) => {
        const {active, over} = event;
        if (!over || active.id === over.id) return;

        const oldIndex = exercises.findIndex((e) => e.uid === active.id);
        const newIndex = exercises.findIndex((e) => e.uid === over.id);

        setExercises(arrayMove(exercises, oldIndex, newIndex));
    };

    // ================= EXERCISE ADD =================
    const loadTypes = async (code: string) => {
        setSelectedGroup(code);
        setTypes(await getExerciseTypesByGroup(code));
    };

    const addExercise = () => {
        const ex: any = types.find((t) => t.id === Number(selectedType));
        if (!ex) return;

        let reps = ex.units.includes("reps") ? 10 : 0;
        let weight = ex.units.includes("weight") ? 10 : 0;
        let minutes = ex.units.includes("minutes") ? 10 : 0;
        let meters = ex.units.includes("meters") ? 10 : 0;

        setExercises([
            ...exercises,
            {
                uid: crypto.randomUUID(),
                id: ex.id,
                name: ex.name,
                sets: [{reps, weight, minutes, meters}],
            },
        ]);
    };

    const removeExercise = (i: number) => {
        if (!window.confirm("Удалить упражнение?")) return;
        setExercises(exercises.filter((_, idx) => idx !== i));
    };

    const matchInt = (input: string): boolean => {
        return /^\d*$/.test(input)
    }

    const matchFloat = (input: string): boolean => {
        return /^\d*\.?\d?$/.test(input)
    }

    // ================= SET OPS =================
    const updateSet = (ei: number, si: number, field: string, value: string, typeParam: string) => {
        if (value == "") {
            return;
        }

        if (typeParam == 'int' && (!matchInt(value) || parseInt(value) <= 0)) {
            return;
        } else if (typeParam == 'float' && (!matchFloat(value) || parseFloat(value) <= 0)) {
            return;
        }
        const copy = [...exercises];
        copy[ei].sets[si][field] = +value;
        setExercises(copy);
    };

    const addSet = (ei: number, sets: SetDTO[], units: string) => {
        const copy = [...exercises];
        let reps = units.includes('reps') ? 10 : 0;
        let weight = units.includes('weight') ? 10 : 0;
        let minutes = units.includes('minutes') ? 10 : 0;
        let meters = units.includes('meters') ? 10 : 0;
        if (sets.length > 0) {
            reps = sets[sets.length - 1].reps;
            weight = sets[sets.length - 1].weight;
            minutes = sets[sets.length - 1].minutes;
            meters = sets[sets.length - 1].meters;
        }
        copy[ei].sets.push({reps, weight, minutes, meters});
        setExercises(copy);
    };

    const removeSet = (ei: number, si: number) => {
        const copy = [...exercises];
        copy[ei].sets.splice(si, 1);
        setExercises(copy);
    };

    // ================= PRESET =================
    const buildPreset = () =>
        exercises
            .map((ex) => {
                const sets = ex.sets
                    .map((s: any) =>
                        s.weight ? `${s.reps}*${s.weight}` : `${s.reps || s.minutes || s.meters}`
                    )
                    .join(",");
                return `${ex.id}:[${sets}]`;
            })
            .join(";");

    // ================= UI =================
    return (
        <div className="page stack">
            <h2>{dayName}</h2>

            <div className="selector">
                <select onChange={(e) => loadTypes(e.target.value)}>
                    <option>Группа</option>
                    {groups.map((g) => (
                        <option key={g.code} value={g.code}>{g.name}</option>
                    ))}
                </select>

                <select onChange={(e) => setSelectedType(e.target.value)}>
                    <option>Упражнение</option>
                    {types.map((t) => (
                        <option key={t.id} value={t.id}>{t.name}</option>
                    ))}
                </select>

                <Button variant="active" onClick={addExercise}>
                    ➕️ Добавить
                </Button>
            </div>

            {loading && <Loader/>}

            <DndContext collisionDetection={closestCenter} onDragEnd={handleDragEnd}>
                <SortableContext
                    items={exercises.map((e) => e.uid)}
                    strategy={verticalListSortingStrategy}
                >
                    {exercises.map((ex, ei) => (
                        <SortableExercise key={ex.uid} id={ex.uid}>
                            {({listeners, attributes}: any) => (
                                <div className="card exercise-card-edit animate">

                                    {/* 🔥 drag handle только на ручку */}
                                    <div className="drag-handle" {...listeners} {...attributes}>
                                        ☰
                                    </div>

                                    <div className="exercise-card-edit-header">
                                        <h3>{ex.name}</h3>
                                        <Button
                                            variant="danger"
                                            onClick={() => removeExercise(ei)}
                                        >
                                            ✕
                                        </Button>
                                    </div>

                                    <div className="sets">
                                        {ex.sets.map((s: any, si: number) => (
                                            <div key={si} className="set-row">
                                                <>{si+1}</>
                                                {s.reps > 0 && (
                                                    <>
                                                        <input
                                                            type="number"
                                                            value={s.reps}
                                                            onChange={(e) => updateSet(ei, si, "reps", e.target.value, 'int')}
                                                            onPointerDown={(e) => e.stopPropagation()}
                                                        />
                                                        <span>повт.</span>
                                                    </>
                                                )}
                                                {s.weight > 0 && (
                                                    <>
                                                        <input
                                                            type="number"
                                                            value={s.weight}
                                                            onChange={(e) => updateSet(ei, si, "weight", e.target.value, 'float')}
                                                            onPointerDown={(e) => e.stopPropagation()}
                                                        />
                                                        <span>кг</span>
                                                    </>
                                                )}
                                                {s.meters > 0 && (
                                                    <>
                                                        <input
                                                            type="number"
                                                            value={s.meters}
                                                            onChange={(e) => updateSet(ei, si, "meters", e.target.value, 'int')}
                                                            onPointerDown={(e) => e.stopPropagation()}
                                                        />
                                                        <span>м</span>
                                                    </>
                                                )}
                                                {s.minutes > 0 && (
                                                    <>
                                                        <input
                                                            type="number"
                                                            value={s.minutes}
                                                            onChange={(e) => updateSet(ei, si, "minutes", e.target.value, 'int')}
                                                            onPointerDown={(e) => e.stopPropagation()}
                                                        />
                                                        <span>мин</span>
                                                    </>
                                                )}

                                                <button
                                                    className="minus"
                                                    onClick={(e) => {
                                                        e.stopPropagation();
                                                        removeSet(ei, si);
                                                    }}
                                                >
                                                    🗑
                                                </button>
                                            </div>
                                        ))}

                                        <Button
                                            onClick={(e) => {
                                                e.stopPropagation();
                                                addSet(ei, ex.sets, ex.units);
                                            }}
                                        >
                                            + сет
                                        </Button>
                                    </div>
                                </div>
                            )}
                        </SortableExercise>
                    ))}
                </SortableContext>
            </DndContext>

            {toast && <div className="toast">{toast}</div>}
        </div>
    );
}
