import {useCallback, useEffect, useState} from "react";
import {useNavigate, useParams} from "react-router-dom";
import {moveToExerciseSession, showCurrentExerciseSession,} from "../api/sessions";
import {finishWorkout} from "../api/workouts";

import ExerciseView from "../components/ExerciseView";
import WorkoutControls from "../components/WorkoutControls";

import {Flag, Loader} from "lucide-react";
import Button from "../components/Button.tsx";

export default function WorkoutSession() {
    const {id} = useParams();
    const workoutID = Number(id);

    const navigate = useNavigate();

    const [session, setSession] = useState<CurrentExerciseSession | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(false);

    // ---------------- LOAD ----------------
    const load = useCallback(async () => {
        try {
            setLoading(true);
            setError(false);

            const data = await showCurrentExerciseSession(workoutID);

            setSession(data);
        } catch (e) {
            console.error("load session failed:", e);
            setError(true);
        } finally {
            setLoading(false);
        }
    }, [workoutID]);

    // ---------------- FIRST LOAD ----------------
    useEffect(() => {
        let cancelled = false;

        const run = async () => {
            try {
                setLoading(true);
                const data = await showCurrentExerciseSession(workoutID);
                if (!cancelled) setSession(data);
            } catch (e) {
                if (!cancelled) setError(true);
            } finally {
                if (!cancelled) setLoading(false);
            }
        };

        run();

        return () => {
            cancelled = true;
        };
    }, [workoutID]);

    // ---------------- MOVE ----------------
    const move = async (next: boolean) => {
        if (isLast && next) return;
        try {
            setLoading(true);
            await moveToExerciseSession(workoutID, next);
            await load();
        } catch (e) {
            console.error(e);
            setError(true);
        }
    };

    // ---------------- FINISH ----------------
    const finish = async () => {
        if (!confirm("Вы уверены, что хотите завершить тренировку?")) return;

        try {
            await finishWorkout(workoutID);
            navigate("/");
        } catch (e) {
            console.error(e);
            setError(true);
        }
    };

    // ================= RENDER =================

    const isFirst = session?.exercise_index === 0;
    const isLast = session?.exercise_index === session?.workout.exercises.length - 1;

    return (
        <div className="page stack">

            {loading && (
                <div className="center">
                    <Loader/>
                </div>
            )}

            {!loading && session && (
                <div>
                    <b>
                        {session.exercise_index + 1} /{" "}
                        {session.workout.exercises.length}
                    </b>
                    <ExerciseView
                        session={session}
                        onAllSetsCompleted={() => move(true)}
                        onReload={load}
                    />
                </div>
            )}
            <Button variant="danger" onClick={finish}>
                <Flag size={18} />
                Завершить всё
            </Button>

            <WorkoutControls
                onPrev={() => move(false)}
                onNext={() => move(true)}
                workoutID={workoutID}
                disablePrev={isFirst}
                disableNext={isLast}
            />
        </div>
    );
}
