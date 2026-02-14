import {useEffect, useState} from "react";
import {addSet, changeSet, completeSet, deleteSet,} from "../api/sets";
import SetRow from "./SetRow.tsx";
import RestTimer from "./RestTimer.tsx";
import Button from "./Button.tsx";
import Toast from "./Toast.tsx";
import "../styles/workout.css";
import {deleteExercise} from "../api/exercises.ts";

export default function ExerciseView({session, onAllSetsCompleted, onReload}) {
    const [sets, setSets] = useState(session.exercise.sets);
    const [toast, setToast] = useState<string | null>(null);
    console.log('session', session)

    useEffect(() => {
        setSets(session.exercise.sets);
    }, [session.exercise.sets]);

    const showError = () => setToast("Ошибка сервера 😢");

    // ---------- ADD ----------
    const handleAdd = async (exerciseID: number, lastSet: FormattedSet | null) => {
        const temp: FormattedSet = {
            id: Date.now(),

            reps: (lastSet?.fact_reps > 0 ? lastSet?.fact_reps : lastSet?.reps) ?? 0,
            weight: (lastSet?.fact_weight > 0 ? lastSet?.fact_weight : lastSet?.weight) ?? 0,
            minutes: (lastSet?.fact_minutes > 0 ? lastSet?.fact_minutes : lastSet?.minutes) ?? 0,
            meters: (lastSet?.fact_meters > 0 ? lastSet?.fact_meters : lastSet?.meters) ?? 0,

            fact_reps: 0,
            fact_weight: 0,
            fact_minutes: 0,
            fact_meters: 0,

            formatted_string: "",
            completed: false,
            completed_at: "",
            index: sets.length,
        };

        setSets(prev => [...prev, temp]);

        try {
            await addSet(exerciseID);
            await onReload(); // сервер даст правильный id
        } catch {
            showError();
            setSets(prev => prev.slice(0, -1));
        }
    };

    const handleDeleteExercise = async (id: number) => {
        if (!window.confirm("Вы уверены, что хотите удалить упражнение из тренировки?")) return;

        try {
            await deleteExercise(id);
            onReload();
        } catch {
            showError();
        }
    };

    // ---------- DELETE ----------
    const handleDeleteSet = async (id: number) => {
        const old = sets;

        setSets(prev => prev.filter(s => s.id !== id));

        try {
            await deleteSet(id);
        } catch {
            showError();
            setSets(old);
        }
    };

    // ---------- COMPLETE ----------
    const handleCompleteSet = async (id: number) => {
        const old = sets; // для rollback

        let updatedSets: FormattedSet[] = [];

        // optimistic
        setSets(prev => {
            updatedSets = prev.map(s =>
                s.id === id ? {...s, completed: !s.completed} : s
            );

            const allDone = updatedSets.every(s => s.completed);
            if (allDone) onAllSetsCompleted?.();

            return updatedSets;
        });

        try {
            await completeSet(id); // 🔥 ВОТ ЭТОГО НЕ ХВАТАЛО
        } catch {
            showError();
            setSets(old); // rollback
        }
    };


    // ---------- CHANGE ----------
    const handleChange = async (id, reps, weight, minutes, meters) => {
        setSets(prev =>
            prev.map(s =>
                s.id === id ? {
                    ...s,
                    fact_reps: reps,
                    fact_weight: weight,
                    fact_minutes: minutes,
                    fact_meters: meters
                } : s
            )
        );

        try {
            await changeSet(id, reps, weight, minutes, meters);
        } catch {
            showError();
        }
    };

    const ex = session.exercise;

    return (
        <div className="exercise-card-view">

            <div className="exercise-card-view-header">
                <div className="exercise-card-view-title">{ex.name}</div>

                {ex.url && <a className="exercise-card-view-link" href={ex.url}>
                    Техника упражнения ↗
                </a>}
            </div>

            <div className="sets">
                {sets.map((s, i) => (
                    <SetRow
                        key={s.id}
                        set={s}
                        index={i}
                        onDelete={() => handleDeleteSet(s.id)}
                        onComplete={() => handleCompleteSet(s.id)}
                        onChange={handleChange}
                    />
                ))}
            </div>

            <RestTimer
                seconds={ex.rest_in_seconds}
                onFinish={() => setToast("Отдых закончен 💪")}
            />

            <div style={{display: "grid", gridTemplateColumns: "1fr 1fr", gap: "8px"}}>
                <Button variant={"ghost"}
                        onClick={() => handleAdd(ex.id, sets.length > 0 ? sets[sets.length - 1] : null)}
                >+ Добавить подход</Button>
                <Button variant={"danger"} onClick={() => handleDeleteExercise(ex.id)}>✕ Убрать упражнение</Button>
            </div>

            {toast && <Toast message={toast} onClose={() => setToast(null)}/>}
        </div>
    );
}
