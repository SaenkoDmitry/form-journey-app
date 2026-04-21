import {useNavigate, useParams} from 'react-router-dom';
import React, {useEffect, useState} from 'react';
import SafeTextRenderer from "../components/SafeTextRenderer.tsx";
import {Check, Loader, Play, Plus, Share2} from "lucide-react";
import Button from "../components/Button.tsx";
import {createShare, getPublicWorkout, getWorkout} from "../api/workouts.ts";
import {moveToCertainExerciseSession} from "../api/sessions.ts";
import {getExerciseGroups} from "../api/exercises.ts";
import ShareSheet from "../components/ShareSheet.tsx";
import {useShare} from "../hooks/useShare.ts";
import Toast from "../components/Toast.tsx";

const WorkoutPage = () => {

    const {id, token} = useParams<{ id?: string; token?: string }>();
    const isPublicMode = !!token;

    const [data, setData] = useState<ReadWorkoutDTO | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const navigate = useNavigate();
    const [groupsMap, setGroupsMap] = useState<Record<string, Group>>({});
    const [shareUrl, setShareUrl] = useState<string | null>(null);
    const [copied, setCopied] = useState(false);
    const [toast, setToast] = useState<string | null>(null);

    const {openShare, isOpen, url, close} = useShare();

    const handleShare = async () => {
        if (!data?.progress?.workout?.id) return;

        try {
            const result = await createShare(data.progress.workout.id);
            await openShare(result.share_url);
        } catch (e: any) {
            console.log('e.status:', e.status);
            console.log(e.message, e.message == "rate limit exceeded")
            if (e.message?.trim() === 'rate limit exceeded') {
                setToast('Слишком много запросов');
            } else {
                setToast('Не удалось создать ссылку');
            }
        }
    };

    useEffect(() => {
        const fetchWorkout = async () => {
            try {
                setLoading(true);

                let workoutData: ReadWorkoutDTO;

                if (isPublicMode) {
                    workoutData = await getPublicWorkout(token!);
                } else {
                    const [wd, groups]: [ReadWorkoutDTO, Group[]] = await Promise.all([
                        getWorkout(Number(id)),
                        getExerciseGroups()
                    ]);
                    workoutData = wd;

                    const map = groups.reduce<Record<string, Group>>((acc, group) => {
                        acc[group.code] = group;
                        return acc;
                    }, {});
                    setGroupsMap(map);
                }

                setData(workoutData);

            } catch (err: any) {
                setError(err.message || 'Не удалось загрузить данные');
            } finally {
                setLoading(false);
            }
        };

        if (isPublicMode && token) {
            fetchWorkout();
        } else if (id) {
            fetchWorkout();
        }
    }, [id, token, isPublicMode]);

    if (loading) return <Loader/>;
    if (error) return <p style={{color: 'red'}}>{error}</p>;
    if (!data) return <p>Данные тренировки не найдены</p>;

    const {progress, stats} = data;
    const {workout, ProgressPercent, RemainingMin, SessionStarted, CompletedExercises, TotalExercises} = progress;

    return <div className={"page stack"}>
        {isPublicMode && (
            <div style={{
                background: 'var(--color-card-alt)',
                color: 'var(--color-text-muted)',
                fontWeight: 'bold',
                padding: '8px 16px',
                borderRadius: '16px',
                textAlign: 'center',
            }}>
                🔗 Тренировка пользователя '{data?.user_first_name || ''}'
            </div>
        )}

        <h2>
            {workout.day_type_name || `Тренировка ${workout.id}`}
            {!isPublicMode && workout.completed && (
                <Button
                    style={{marginLeft: 10}}
                    variant="ghost"
                    onClick={handleShare}
                    disabled={copied}
                >
                    {copied ? <Check size={14} /> : <Share2 size={14} />}
                </Button>
            )}
        </h2>

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
            <div className={"card"}>
                {stats.cardio_time > 0 && <p><strong>🫀 Время кардио:</strong> {stats.cardio_time} мин</p>}
                {stats.total_weight > 0 && <p><strong>🏋 Общий вес:</strong> {stats.total_weight} кг</p>}
                <br/>
                {stats.exercise_map && <strong>Группы мышц:</strong>}
                {stats.exercise_map && [...new Set(
                    Object.values(stats.exercise_map).map(ex => ex.group_name)
                )].map(groupName => {
                    const exercisesInGroup = Object.values(stats.exercise_map).filter(
                        ex => ex.group_name === groupName
                    );
                    const exerciseCount = exercisesInGroup.length;
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

                    return (
                        <p key={groupName}>
                            • {groupName} — {exerciseCount} упражн.,
                            {totalWeight > 0 ? ` общий вес: ${totalWeight} кг` : ''}
                            {totalTime > 0 ? ` общее время: ${totalTime} мин` : ''}
                        </p>
                    );
                })}
            </div>
        </div>}

        {!isPublicMode && data.progress.SessionStarted && (
            <>
                <Button variant={"active"} onClick={() => navigate(`/sessions/${data?.progress.workout.id}`)}>
                    <Play size={14}/>К тренировке</Button>
                <Button
                    variant="primary"
                    onClick={() => navigate(`/workouts/${id}/add-exercise`)}
                >
                    <Plus size={14}/>
                    Добавить упражнение
                </Button>
            </>
        )}

        {/* Упражнения */}
        <h3>Упражнения ({CompletedExercises}/{TotalExercises})</h3>
        <div style={{listStyle: "none", padding: 0}}>
            {workout.exercises?.map((ex: FormattedExercise) => (
                <div className="card"
                     key={ex.id}
                     style={{
                         borderRadius: "20px",
                         padding: "1rem",
                         marginBottom: "0.5rem",
                         cursor: !isPublicMode && !workout.completed ? 'pointer' : 'default'
                     }}
                     onClick={(e) => {
                         if (isPublicMode || workout.completed) {
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

        {isPublicMode && (
            <div style={{
                textAlign: 'center',
                marginTop: '2rem',
                padding: '1.5rem',
                borderTop: '1px solid #eee'
            }}>
                <Button variant="primary" onClick={() => navigate('/')}>
                    Перейти в приложение
                </Button>
            </div>
        )}

        <ShareSheet isOpen={isOpen} url={url} onClose={close} />

        {toast && (
            <Toast message={toast} onClose={() => setToast(null)}/>
        )}
    </div>;
};

export default WorkoutPage;
