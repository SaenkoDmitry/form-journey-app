import {useEffect, useRef, useState} from "react";
import {useParams} from "react-router-dom";
import {getProgram} from "../api/days";
import {parsePreset, savePreset} from "../api/presets";
import {getExerciseGroups, getExerciseTypesByGroup} from "../api/exercises";
import Button from "../components/Button";
import "../styles/DayDetailsPage.css";

export default function DayDetailsPage() {
    const {programId, dayId} = useParams();

    const [exercises, setExercises] = useState<any[]>([]);
    const [dayName, setDayName] = useState("");

    const [groups, setGroups] = useState<Group[]>([]);
    const [types, setTypes] = useState<any[]>([]);
    const [selectedGroup, setSelectedGroup] = useState("");
    const [selectedType, setSelectedType] = useState("");

    const [toast, setToast] = useState<string | null>(null);

    const autosaveTimer = useRef<any>(null);
    const firstLoad = useRef(true);

    // ---------------- LOAD ----------------
    const load = async () => {
        const program = await getProgram(Number(programId));
        const day = program.day_types.find((d) => d.id === Number(dayId));
        if (!day) return;

        setDayName(day.name);

        if (day.preset) {
            const parsed = await parsePreset(day.preset);

            setExercises(
                parsed.exercises.map((ex: any) => ({
                    id: ex.id,
                    name: ex.name,
                    sets: ex.sets,
                }))
            );
        }

        firstLoad.current = false;
    };

    useEffect(() => {
        load();
        getExerciseGroups().then((groups: Group[]) => setGroups(groups));
    }, []);

    // ---------------- AUTOSAVE (debounce 500ms) ----------------
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

    // ---------------- DRAG ----------------
    const onDragStart = (e: any, i: number) => {
        e.dataTransfer.setData("index", i);
    };

    const onDrop = (e: any, i: number) => {
        const from = Number(e.dataTransfer.getData("index"));
        const copy = [...exercises];
        const [moved] = copy.splice(from, 1);
        copy.splice(i, 0, moved);
        setExercises(copy);
    };

    // ---------------- EXERCISE ADD ----------------
    const loadTypes = async (code: string) => {
        setSelectedGroup(code);
        const exerciseTypes = await getExerciseTypesByGroup(code);
        setTypes(exerciseTypes);
    };

    const addExercise = () => {
        const ex: ExerciseType = types.find((t) => t.id === Number(selectedType));
        if (!ex) return;

        let reps = ex.units.includes("reps") ? 10 : 0
        let weight = ex.units.includes("weight") ? 10 : 0
        let minutes = ex.units.includes("minutes") ? 10 : 0
        let meters = ex.units.includes("meters") ? 10 : 0
        setExercises([
            ...exercises,
            {
                id: ex.id,
                name: ex.name,
                sets: [{reps: reps, weight: weight, minutes: minutes, meters: meters}],
            },
        ]);
    };

    const removeExercise = (i: number) => {
        const confirmed = window.confirm("Вы уверены, что хотите удалить упражнение?");
        if (!confirmed) return;

        setExercises(exercises.filter((_, idx) => idx !== i));
    }

    // ---------------- SET OPS ----------------
    const updateSet = (ei: number, si: number, field: string, value: number) => {
        const copy = [...exercises];
        copy[ei].sets[si][field] = value;
        setExercises(copy);
    };

    const addSet = (ei: number, sets: SetDTO[]) => {
        const copy = [...exercises];
        let reps = 0
        let weight = 0
        let minutes = 0
        let meters = 0
        if (sets.length > 0) {
            reps = sets[sets.length-1].reps
            weight = sets[sets.length-1].weight
            minutes = sets[sets.length-1].minutes
            meters = sets[sets.length-1].meters
        }
        copy[ei].sets.push({reps: reps, weight: weight, minutes: minutes, meters: meters});
        setExercises(copy);
    };

    const removeSet = (ei: number, si: number) => {
        if (si == 0) {
            return
        }
        const copy = [...exercises];
        copy[ei].sets.splice(si, 1);
        setExercises(copy);
    };

    // ---------------- PRESET ----------------
    const buildPreset = () =>
        exercises
            .map((ex) => {
                const sets = ex.sets
                    .map((s: any) =>
                        s.weight
                            ? `${s.reps}*${s.weight}`
                            : `${s.reps || s.minutes || s.meters}`
                    )
                    .join(",");
                return `${ex.id}:[${sets}]`;
            })
            .join(";");

    // ---------------- UI ----------------
    return (
        <div className="page stack">
            <h2>{dayName}</h2>

            {/* selector */}
            <div className="selector">
                <select style={{minHeight: "var(--tap-size)"}} onChange={(e) => loadTypes(e.target.value)}>
                    <option>Группа</option>
                    {groups.map((g) => (
                        <option key={g.code} value={g.code}>
                            {g.name}
                        </option>
                    ))}
                </select>

                <select style={{minHeight: "var(--tap-size)"}} onChange={(e) => setSelectedType(e.target.value)}>
                    <option>Упражнение</option>
                    {types.map((t) => (
                        <option key={t.id} value={t.id}>
                            {t.name}
                        </option>
                    ))}
                </select>

                <Button variant={"primary"} onClick={addExercise}>➕️ Добавить</Button>
            </div>

            {exercises.map((ex, ei) => (
                <div
                    key={ex.id}
                    className="card exercise-card animate"
                    onDragOver={(e) => e.preventDefault()}
                    onDrop={(e) => onDrop(e, ei)}
                >
                    <div className="exercise-header">
                        <div
                            className="drag-handle"
                            draggable
                            onDragStart={(e) => onDragStart(e, ei)}
                        >
                            ☰
                        </div>

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
                                {s.meters > 0 && <div>
                                    <input
                                        type="number"
                                        value={s.meters}
                                        onChange={(e) =>
                                            updateSet(
                                                ei,
                                                si,
                                                "meters",
                                                +e.target.value
                                            )
                                        }
                                    />
                                    <span> метр.</span>
                                </div>}
                                {s.minutes > 0 && <div>
                                    <input
                                        type="number"
                                        value={s.minutes}
                                        onChange={(e) =>
                                            updateSet(
                                                ei,
                                                si,
                                                "minutes",
                                                +e.target.value
                                            )
                                        }
                                    />
                                    <span> мин.</span>
                                </div>}
                                {s.weight > 0 && <div>
                                    <input
                                        type="number"
                                        value={s.reps}
                                        onChange={(e) =>
                                            updateSet(
                                                ei,
                                                si,
                                                "reps",
                                                +e.target.value
                                            )
                                        }
                                    />
                                    <span> × </span>
                                    <input
                                        type="number"
                                        value={s.weight}
                                        onChange={(e) =>
                                            updateSet(
                                                ei,
                                                si,
                                                "weight",
                                                +e.target.value
                                            )
                                        }
                                    />
                                    <span> кг</span>
                                </div>}
                                <button
                                    className="minus"
                                    onClick={() => removeSet(ei, si)}
                                >
                                    🗑
                                </button>
                            </div>
                        ))}

                        <Button onClick={() => addSet(ei, ex.sets)}>
                            + сет
                        </Button>
                    </div>
                </div>
            ))}

            {toast && <div className="toast">{toast}</div>}
        </div>
    );
}
