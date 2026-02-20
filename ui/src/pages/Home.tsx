import React, {useEffect, useRef, useState, useCallback} from 'react';
import {useNavigate} from 'react-router-dom';
import WorkoutCard from '../components/WorkoutCard';
import {useAuth} from '../context/AuthContext';
import '../styles/App.css';
import Button from "../components/Button.tsx";
import {deleteWorkout, getWorkouts} from "../api/workouts.ts";
import {Frown, Loader, Play, RotateCcw, Trash2} from "lucide-react";

const LIMIT = 10;

type Status = 'idle' | 'loading' | 'error';

const Home: React.FC = () => {
    const {user} = useAuth();

    const [workouts, setWorkouts] = useState<Workout[]>([]);
    const [pagination, setPagination] = useState<Pagination | null>(null);
    const [status, setStatus] = useState<Status>('idle');
    const [hasMore, setHasMore] = useState(true);

    const offsetRef = useRef(0);
    const loaderRef = useRef<HTMLDivElement>(null);
    const isFetchingRef = useRef(false); // 🔥 защита от дублей

    const navigate = useNavigate();

    // ---------------- DELETE ----------------
    const handleDelete = async (id: number) => {
        if (!confirm("Вы уверены, что хотите удалить тренировку?")) return;

        await deleteWorkout(id);
        setWorkouts(prev => prev.filter(w => w.id !== id));
    };

    // ---------------- FETCH ----------------
    const fetchWorkouts = useCallback(async () => {
        if (isFetchingRef.current || !hasMore) return;

        isFetchingRef.current = true;
        setStatus('loading');

        try {
            const nextOffset = offsetRef.current;

            const data: ShowMyWorkoutsResult =
                await getWorkouts(nextOffset, LIMIT);

            // 🔥 если сервер вернул пусто — прекращаем
            if (data.items.length === 0) {
                setHasMore(false);
                return;
            }

            setWorkouts(prev => [...prev, ...data.items]);
            setPagination(data.pagination);

            offsetRef.current += data.items.length;
            setHasMore(offsetRef.current < data.pagination.total);

            setStatus('idle');
        } catch (e) {
            setStatus('error');
        } finally {
            isFetchingRef.current = false;
        }
    }, [hasMore]);

    // ---------------- INITIAL LOAD ----------------
    useEffect(() => {
        if (!user) return;

        // сброс при смене пользователя
        setWorkouts([]);
        setPagination(null);
        setHasMore(true);
        offsetRef.current = 0;

        fetchWorkouts();
    }, [user]);

    // ---------------- OBSERVER ----------------
    useEffect(() => {
        if (!loaderRef.current) return;

        const observer = new IntersectionObserver((entries) => {
            if (entries[0].isIntersecting) {
                fetchWorkouts();
            }
        });

        observer.observe(loaderRef.current);

        return () => observer.disconnect();
    }, [fetchWorkouts]);

    const isEmpty = pagination && pagination.total === 0;

    return (
        <div className="page stack">
            <h1>Мои тренировки</h1>

            <Button
                variant="active"
                onClick={() => navigate('/start')}
            >
                <Play/> Начать новую тренировку
            </Button>

            {!isEmpty && (
                <div style={{display: 'flex', flexDirection: 'column', gap: 12}}>
                    {workouts.map((w, idx) => (
                        <div
                            key={w.id}
                            onClick={() => navigate(`/workouts/${w.id}`)}
                            className="workout-item"
                        >
                            <WorkoutCard w={w} idx={idx + 1}/>

                            <div className="workout-actions">
                                {!w.completed && (
                                    <Button
                                        variant="active"
                                        onClick={(e) => {
                                            navigate(`/sessions/${w.id}`);
                                            e.stopPropagation();
                                        }}
                                    >
                                        <Play size={14}/>
                                    </Button>
                                )}

                                <Button
                                    variant="danger"
                                    onClick={(e) => {
                                        e.stopPropagation();
                                        handleDelete(w.id);
                                    }}
                                >
                                    <Trash2 size={14}/>
                                </Button>
                            </div>
                        </div>
                    ))}
                </div>
            )}

            {status === 'loading' && <Loader/>}

            {status === 'error' && (
                <div className="error-block">
                    <div className="error-message">
                        <Frown size={28} />
                        <span>Ошибка загрузки</span>
                    </div>

                    <Button
                        variant="danger"
                        onClick={() => {
                            setStatus('idle');
                            fetchWorkouts();
                        }}
                    >
                        <RotateCcw size={16}/>
                        Повторить
                    </Button>
                </div>
            )}

            <div ref={loaderRef} style={{height: 20}}/>

            {pagination && pagination.total > 0 && (
                <p>
                    {Math.min(offsetRef.current, pagination.total)} из {pagination.total}
                </p>
            )}
        </div>
    );
};

export default Home;
