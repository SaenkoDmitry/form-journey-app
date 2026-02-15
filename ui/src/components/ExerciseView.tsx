import { useEffect, useState } from "react";
import { addSet, changeSet, completeSet, deleteSet } from "../api/sets";
import SetRow from "./SetRow";
import Button from "./Button";
import Toast from "./Toast";
import "../styles/workout.css";
import { deleteExercise } from "../api/exercises";
import { Plus, X } from "lucide-react";
import { useGlobalTimer } from "../context/TimerContext";

export default function ExerciseView({ session, onAllSetsCompleted, onReload }) {
    const [sets, setSets] = useState(session.exercise.sets);
    const [toast, setToast] = useState<string | null>(null);
    const [justCompletedSet, setJustCompletedSet] = useState<number | null>(null);
    const { start } = useGlobalTimer();

    const ex = session.exercise;

    useEffect(() => {
        setSets(session.exercise.sets);
    }, [session.exercise.sets]);

    const showError = () => setToast("Ошибка сервера 😢");

    const handleAdd = async (exerciseID: number, lastSet: any) => {
        const temp = {
            id: Date.now(),
            reps: lastSet?.fact_reps > 0 ? lastSet.fact_reps : lastSet?.reps ?? 0,
            weight: lastSet?.fact_weight > 0 ? lastSet.fact_weight : lastSet?.weight ?? 0,
            fact_reps: 0,
            fact_weight: 0,
            completed: false,
            index: sets.length,
        };

        setSets(prev => [...prev, temp]);

        try {
            await addSet(exerciseID);
            await onReload();
        } catch {
            showError();
            setSets(prev => prev.slice(0, -1));
        }
    };

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

    const handleCompleteSet = async (id: number) => {
        const old = sets;
        let updatedSets: any[] = [];

        setSets(prev => {
            updatedSets = prev.map(s => s.id === id ? { ...s, completed: !s.completed } : s);

            const setNowCompleted = updatedSets.find(s => s.id === id)?.completed;
            if (setNowCompleted) setJustCompletedSet(id);

            const allDone = updatedSets.every(s => s.completed);
            if (allDone) onAllSetsCompleted?.();

            return updatedSets;
        });

        const currentSet = sets.find(s => s.id === id);
        if (currentSet) {
            let reps = currentSet.fact_reps > 0 ? currentSet.fact_reps : currentSet.reps;
            let weight = currentSet.fact_weight > 0 ? currentSet.fact_weight : currentSet.weight;
            await handleChange(id, reps, weight);
        }

        try {
            await completeSet(id);
        } catch {
            showError();
            setSets(old);
        }
    };

    const handleChange = async (id, reps, weight) => {
        setSets(prev =>
            prev.map(s =>
                s.id === id ? { ...s, fact_reps: reps, fact_weight: weight } : s
            )
        );

        try {
            await changeSet(id, reps, weight, 0, 0);
        } catch {
            showError();
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

    // ---- Автостарт глобального таймера ----
    useEffect(() => {
        if (justCompletedSet !== null) {
            start(ex.rest_in_seconds); // старт таймера
            setJustCompletedSet(null);
        }
    }, [justCompletedSet, start, ex.rest_in_seconds]);

    return (
        <div className="exercise-card-view">
            <div className="exercise-card-view-header">
                <div className="exercise-card-view-title">{ex.name}</div>
                {ex.url && <a className="exercise-card-view-link" href={ex.url}>Техника упражнения ↗</a>}
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

            <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: "8px" }}>
                <Button variant={"ghost"} onClick={() => handleAdd(ex.id, sets.length > 0 ? sets[sets.length - 1] : null)}>
                    Добавить подход
                </Button>
                <Button variant={"danger"} onClick={() => handleDeleteExercise(ex.id)}>
                    Убрать упражнение
                </Button>
            </div>

            {toast && <Toast message={toast} onClose={() => setToast(null)}/>}
        </div>
    );
}
