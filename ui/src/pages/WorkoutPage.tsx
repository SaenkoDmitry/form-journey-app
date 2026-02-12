import {useParams} from 'react-router-dom';
import React, {useEffect, useState} from 'react';
import SafeTextRenderer from "../components/SafeTextRenderer.tsx";

const WorkoutPage = () => {
    const {id} = useParams<{ id: string }>();
    const [data, setData] = useState<ReadWorkoutDTO | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchWorkout = async () => {
            try {
                const res = await fetch(`/api/workouts/${id}`);
                if (!res.ok) throw new Error(`Ошибка: ${res.status}`);
                const json: ReadWorkoutDTO = await res.json();
                setData(json);
            } catch (err: any) {
                setError(err.message || 'Не удалось загрузить данные');
            } finally {
                setLoading(false);
            }
        };

        fetchWorkout();
    }, [id]);

    if (loading) return <p>Загрузка...</p>;
    if (error) return <p style={{color: 'red'}}>{error}</p>;
    if (!data) return <p>Данные тренировки не найдены</p>;

    const {progress, Stats} = data;
    const {workout, ProgressPercent, RemainingMin, SessionStarted, CompletedExercises, TotalExercises} = progress;

    return (
        <div style={{maxWidth: '700px', margin: '0 auto', padding: '1rem'}}>
            <h2>{workout.day_type_name || `Тренировка ${workout.id}`}</h2>
            <p>
                Статус: {workout.status} {progress?.workout?.duration && <span><span>~ </span>{progress.workout.duration}</span>}
            </p>
            <p>{workout.started_at}</p>
            {RemainingMin !== undefined && RemainingMin > 0 && <p>Оставшееся время: {RemainingMin} мин</p>}

            {/* Прогресс тренировки */}
            <div style={{margin: '1rem 0'}}>
                <div style={{background: '#eee', borderRadius: '8px', overflow: 'hidden', height: '20px'}}>
                    <div
                        style={{
                            width: `${ProgressPercent}%`,
                            background: '#4caf50',
                            height: '100%',
                            transition: 'width 0.3s',
                        }}
                    />
                </div>
                <p>{ProgressPercent}% выполнено</p>
            </div>

            {/* Упражнения */}
            <h3>Упражнения ({CompletedExercises}/{TotalExercises})</h3>
            <ul style={{listStyle: "none", padding: 0}}>
                {workout.exercises?.map((ex: FormattedExercise) => (
                    <li
                        key={ex.id}
                        style={{
                            border: "1px solid #ddd",
                            borderRadius: "8px",
                            padding: "0.5rem",
                            marginBottom: "0.5rem",
                        }}
                    >
                        <strong>{ex.name}</strong>
                        <ul style={{paddingLeft: "1rem"}}>
                            {ex.sets?.map((set: any) => {
                                    return <li key={set.ID} style={{marginBottom: "0.5rem"}}>
                                        <SafeTextRenderer html={set.formatted_string}/>
                                    </li>
                                }
                            )}
                            <div
                                style={{
                                    background: "#eee",
                                    height: "8px",
                                    borderRadius: "4px",
                                    overflow: "hidden",
                                    marginTop: "2px"
                                }}>
                                <div style={{
                                    width: `${ex.sets?.filter((set: FormattedSet) => set.completed).length / ex.sets?.length * 100}%`,
                                    height: "100%",
                                    background: "#4caf50",
                                    transition: "width 0.3s",
                                }}/>
                            </div>
                        </ul>
                    </li>
                ))}
            </ul>

            {(Stats.CardioTime > 0 || Stats.TotalWeight > 0) && <h3>Статистика</h3>}

            <div>
                {Stats.CardioTime > 0 && <p><strong>🫀 Время кардио:</strong> {Stats.CardioTime} мин</p>}
                {Stats.TotalWeight > 0 && <p><strong>🏋 Общий вес:</strong> {Stats.TotalWeight} кг</p>}
            </div>
        </div>
    );
};

export default WorkoutPage;
