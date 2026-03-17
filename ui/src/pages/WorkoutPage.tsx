import {useNavigate, useParams} from 'react-router-dom';
import React, {useEffect, useState} from 'react';
import SafeTextRenderer from "../components/SafeTextRenderer.tsx";
import {Loader, Play, Plus} from "lucide-react";
import Button from "../components/Button.tsx";
import {getWorkout} from "../api/workouts.ts";
import {moveToCertainExerciseSession} from "../api/sessions.ts";
import {getExerciseGroups} from "../api/exercises.ts";

const WorkoutPage = () => {
    const {id} = useParams<{ id: number }>();
    const [data, setData] = useState<ReadWorkoutDTO | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const navigate = useNavigate();
    const [groupsMap, setGroupsMap] = useState<Record<string, Group>>({});

    useEffect(() => {
        const fetchWorkout = async () => {
            try {
                setLoading(true);

                const [workoutData, groups]: [ReadWorkoutDTO, Group[]] = await Promise.all([
                    getWorkout(Number(id)),
                    getExerciseGroups()
                ]);

                setData(workoutData);

                const map = groups.reduce<Record<string, Group>>((acc, group) => {
                    acc[group.code] = group;
                    return acc;
                }, {});

                setGroupsMap(map);

            } catch (err: any) {
                setError(err.message || 'Не удалось загрузить данные');
            } finally {
                setLoading(false);
            }
        };

        fetchWorkout();
    }, [id]);

    if (loading) return <Loader/>;
    if (error) return <p style={{color: 'red'}}>{error}</p>;
    if (!data) return <p>Данные тренировки не найдены</p>;

    const {progress, stats} = data;
    const {workout, ProgressPercent, RemainingMin, SessionStarted, CompletedExercises, TotalExercises} = progress;

    return <div className={"page stack"}>
        <h2>{workout.day_type_name || `Тренировка ${workout.id}`}</h2>
        <span>
            Статус: {workout.status} {progress?.workout?.duration &&
            <span><span>~ </span>{progress.workout.duration}</span>}
        </span>
        <span>{workout.started_at}</span>
        {RemainingMin !== undefined && RemainingMin > 0 && <span>Оставшееся время: {RemainingMin} мин</span>}

        {/* Прогресс тренировки */}
        <div>
            <div style={{background: '#eee', borderRadius: '8px', overflow: 'hidden', height: '20px'}}>
                <div
                    style={{
                        width: `${ProgressPercent}%`,
                        background: ProgressPercent >= 85 ? 'var(--color-active)' : 50 < ProgressPercent && ProgressPercent < 85 ? 'var(--color-attention)' : 'var(--color-danger)',
                        height: '100%',
                    }}
                />
            </div>
            <div style={{marginTop: 10}}>{ProgressPercent}% выполнено</div>
        </div>

        {!data.progress.SessionStarted && <div>
            {(stats.cardio_time > 0 || stats.total_weight > 0) && <h3>Статистика</h3>}
            <div>
                {stats.cardio_time > 0 && <p><strong>🫀 Время кардио:</strong> {stats.cardio_time} мин</p>}
                {stats.total_weight > 0 && <p><strong>🏋 Общий вес:</strong> {stats.total_weight} кг</p>}
                {stats.exercise_map && <strong>Группы мышц:</strong>}
                {stats.exercise_map && [...new Set(
                    Object.values(stats.exercise_map).map(ex => ex.group_name)
                )].map(groupName => {
                    // Фильтруем упражнения этой группы
                    const exercisesInGroup = Object.values(stats.exercise_map).filter(
                        ex => ex.group_name === groupName
                    );

                    console.log('exercisesInGroup', exercisesInGroup)
                    console.log('stats.exercise_weight_map', stats.exercise_weight_map)
                    console.log('stats.exercise_time_map', stats.exercise_time_map)

                    // Кол-во упражнений
                    const exerciseCount = exercisesInGroup.length;

                    // Общий вес или время
                    let totalWeight = 0;
                    let totalTime = 0;

                    exercisesInGroup.forEach(ex => {
                        const exId = ex.id;
                        if (stats.exercise_weight_map[exId]) {
                            totalWeight += stats.exercise_weight_map[exId];
                        }
                        if (stats.exercise_time_map[exId]) {
                            totalTime += stats.exercise_time_map[exId];
                        }
                    });

                    console.log('totalWeight', totalWeight, 'totalTime', totalTime)

                    return (
                        <p key={groupName}>
                            • {groupName} — {exerciseCount} упражн{exerciseCount > 1 ? 'ений' : 'ение'},
                            {totalWeight > 0 ? ` общий вес: ${totalWeight} кг` : ''}
                            {totalTime > 0 ? ` общее время: ${totalTime} мин` : ''}
                        </p>
                    );
                })}
            </div>
        </div>}

        {data.progress.SessionStarted &&
            <Button variant={"active"} onClick={() => navigate(`/sessions/${data?.progress.workout.id}`)}>
                <Play size={14}/>К тренировке</Button>
        }

        {data.progress.SessionStarted && (
            <Button
                variant="primary"
                onClick={() => navigate(`/workouts/${id}/add-exercise`)}
            >
                <Plus size={14}/>
                Добавить упражнение
            </Button>
        )}

        {/* Упражнения */}
        <h3>Упражнения ({CompletedExercises}/{TotalExercises})</h3>
        <div style={{listStyle: "none", padding: 0}}>
            {workout.exercises?.map((ex: FormattedExercise) => (
                <div className="card"
                     key={ex.id}
                     style={{
                         border: "1px solid #ddd",
                         borderRadius: "8px",
                         padding: "1rem",
                         marginBottom: "0.5rem",
                     }}
                     onClick={(e) => {
                         if (workout.completed) {
                             e.stopPropagation();
                             return;
                         }
                         moveToCertainExerciseSession(workout.id, ex.index).then(() => {
                             navigate(`/sessions/${data?.progress.workout.id}`);
                         });
                     }}
                >
                    <div className={"card-header"}>{ex.name}</div>
                    <div className={"card-body"}>
                        {ex.sets?.map((set: FormattedSet) => {
                                return <div key={set.id} style={{listStyle: "none", padding: 0, margin: "1rem 0rem"}}>
                                    <SafeTextRenderer html={set.formatted_string}/>
                                </div>
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
                                background: "var(--color-active)",
                                transition: "width 0.3s",
                            }}/>
                        </div>
                    </div>
                </div>
            ))}
        </div>
    </div>;
};

export default WorkoutPage;
