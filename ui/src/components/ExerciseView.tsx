import {useEffect, useState} from "react";
import {addSet, changeSet, completeSet, deleteSet,} from "../api/sets";
import SetRow from "./SetRow.tsx";
import RestTimer from "./RestTimer.tsx";
import Button from "./Button.tsx";
import Toast from "./Toast.tsx";
import "../styles/workout.css";
import {deleteExercise} from "../api/exercises.ts";
import {ArrowDown, ArrowUp, Loader, Plus, Trash2} from "lucide-react";
import {Link} from "react-router-dom";
import {api} from "../api/client.ts";

export default function ExerciseView({session, onAllSetsCompleted, onReload}) {
    const [sets, setSets] = useState(session.exercise.sets);
    const [toast, setToast] = useState<string | null>(null);

    const [restTrigger, setRestTrigger] = useState(0);

    // --- inline видео ---
    const [videoUrl, setVideoUrl] = useState<string | null>(null);
    const [videoLoading, setVideoLoading] = useState(false);
    const [videoError, setVideoError] = useState<string | null>(null);
    const [videoOpen, setVideoOpen] = useState(false);

    const showError = () => setToast("Ошибка сервера 😢");

    const fetchVideoUrl = async (originalUrl: string) => {
        const CACHE_KEY = `video-${originalUrl}`;
        const cached = localStorage.getItem(CACHE_KEY);
        if (cached) {
            const { url, expires } = JSON.parse(cached);
            if (Date.now() < expires) {
                setVideoUrl(url);
                return;
            }
        }

        try {
            setVideoLoading(true);
            const data = await api<{ url: string }>(
                `/api/video/link?url=${encodeURIComponent(originalUrl)}`
            );
            setVideoUrl(data.url);

            localStorage.setItem(CACHE_KEY, JSON.stringify({
                url: data.url,
                expires: Date.now() + 4 * 60 * 1000 // кэш 4 минуты
            }));
        } catch (e: any) {
            setVideoError(e.message || "Ошибка загрузки видео");
            setToast("Ошибка загрузки видео 😢");
        } finally {
            setVideoLoading(false);
        }
    };

    // --- useEffect на первоначальные подходы ---
    useEffect(() => {
        setSets(session.exercise.sets);
    }, [session.exercise.sets]);

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

            const justCompleted = updatedSets.find(s => s.id === id)?.completed;

            // 🔥 если подход завершён — запускаем отдых
            if (justCompleted) {
                setRestTrigger(Date.now());
            }

            const allDone = updatedSets.every(s => s.completed);
            if (allDone) onAllSetsCompleted?.();

            return updatedSets;
        });

        const currentSet = sets.find(s => s.id === id);
        if (currentSet) {
            let reps = currentSet.fact_reps > 0 ? currentSet.fact_reps : currentSet.reps;
            let weight = currentSet.fact_weight > 0 ? currentSet.fact_weight : currentSet.weight;
            let minutes = currentSet.fact_minutes > 0 ? currentSet.fact_minutes : currentSet.minutes;
            let meters = currentSet.fact_meters > 0 ? currentSet.fact_meters : currentSet.meters;
            await handleChange(id, reps, weight, minutes, meters);
        }

        try {
            await completeSet(id);
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
                <div className="exercise-card-view-title">
                    {ex.name}
                </div>

                <Button
                    variant="danger"
                    style={{position: "absolute", right: 0}}
                    onClick={() => handleDeleteExercise(ex.id)}
                >
                    <Trash2 size={14}/>
                </Button>

                {ex.url && (
                    <Button
                        variant="ghost"
                        style={{ marginTop: 8 }}
                        onClick={() => {
                            setVideoOpen(!videoOpen);
                            if (!videoUrl) fetchVideoUrl(ex.url);
                        }}
                    >
                        {videoOpen ? <ArrowUp /> : <ArrowDown />} Техника упражнения
                    </Button>
                )}
            </div>

            {videoOpen && (
                <div style={{ marginTop: 8, padding: 8, borderRadius: 8, border: "1px solid #eee" }}>
                    {videoLoading && <Loader />}
                    {videoError && <div>{videoError}</div>}
                    {videoUrl && !videoLoading && (
                        <video
                            src={videoUrl}
                            controls
                            playsInline
                            style={{ width: "100%", borderRadius: 12 }}
                        />
                    )}
                </div>
            )}

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
                autoStartTrigger={restTrigger}
                workoutID={session.workout.id}
            />

            <div style={{display: "grid", gridTemplateColumns: "1fr", gap: "8px"}}>
                <Button variant={"primary"}
                        onClick={() => handleAdd(ex.id, sets.length > 0 ? sets[sets.length - 1] : null)}
                ><Plus size={14}/>Добавить подход</Button>
            </div>

            {toast && <Toast message={toast} onClose={() => setToast(null)}/>}
        </div>
    );
}
