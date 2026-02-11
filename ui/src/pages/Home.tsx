import React, {useCallback, useEffect, useRef, useState} from 'react';
import {useNavigate} from 'react-router-dom';
import WorkoutCard from '../components/WorkoutCard';
import {useAuth} from '../context/AuthContext';
import '../styles/App.css';
import Button from "../components/Button.tsx";

const LIMIT = 10;

const Home: React.FC = () => {
    const {user, loading: authLoading} = useAuth();

    const [workouts, setWorkouts] = useState<Workout[]>([]);
    const [pagination, setPagination] = useState<Pagination | null>(null);
    const [offset, setOffset] = useState(0);
    const [loading, setLoading] = useState(false);

    const loaderRef = useRef<HTMLDivElement>(null);
    const navigate = useNavigate();

    // ---------------- WORKOUTS ----------------

    const fetchWorkouts = useCallback(
        async (nextOffset: number, append = true) => {
            if (loading) return;

            setLoading(true);

            const res = await fetch(
                `/api/workouts?offset=${nextOffset}&limit=${LIMIT}`,
                {credentials: 'include'}
            );

            const data = await res.json();

            setWorkouts((prev) =>
                append ? [...prev, ...data.items] : data.items
            );

            setPagination(data.pagination);
            setOffset(nextOffset + LIMIT);

            setLoading(false);
        },
        [loading]
    );

    // ---------------- INIT ----------------

    useEffect(() => {
        if (!user) return;

        fetchWorkouts(0, false);
    }, [user]);

    // ---------------- INFINITE SCROLL ----------------

    useEffect(() => {
        if (!loaderRef.current || !pagination) return;

        const observer = new IntersectionObserver((entries) => {
            if (!entries[0].isIntersecting) return;

            const hasMore =
                pagination.offset + pagination.limit < pagination.total;

            if (hasMore) fetchWorkouts(offset, true);
        });

        observer.observe(loaderRef.current);

        return () => observer.disconnect();
    }, [offset, pagination, fetchWorkouts]);

    // ---------------- UI ----------------

    if (authLoading) return <p>Загрузка...</p>;

    if (!user)
        return (
            <div style={{textAlign: 'center', marginTop: 40}}>
                <h3>Войдите через Telegram 👆</h3>
            </div>
        );

    return (
        <div className="page stack">
            <h1>Мои тренировки</h1>

            {user && (
                <Button
                    variant="primary"
                    onClick={() => alert('Начало новой тренировки!')}
                >
                    Начать тренировку
                </Button>
            )}

            <div style={{display: 'flex', flexDirection: 'column', gap: 12}}>
                {workouts.map((w, idx) => (
                    <div
                        key={w.id}
                        onClick={() => navigate(`/workout/${w.id}`)}
                        className="workout-item"
                    >
                        <WorkoutCard w={w} idx={idx+1}/>
                        {/*<hr/>*/}
                    </div>
                ))}
            </div>

            <div ref={loaderRef} style={{padding: 20}}>
                {loading && 'Загрузка...'}
            </div>

            {pagination && (
                <p>
                    {Math.min(
                        pagination.limit + pagination.offset,
                        pagination.total
                    )}{' '}
                    из {pagination.total}
                </p>
            )}
        </div>
    );
};

export default Home;
