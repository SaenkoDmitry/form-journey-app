import {CalendarRange, Loader} from "lucide-react";
import React, {useEffect, useRef, useState} from "react";
import Toast from "../components/Toast.tsx";
import {
    getExerciseGroups,
    getExerciseStats,
    getExerciseTypesByGroup,
} from "../api/exercises.ts";
import {useParams} from "react-router-dom";

import {
    CartesianGrid,
    Line,
    LineChart,
    ResponsiveContainer,
    Tooltip,
    XAxis,
    YAxis
} from "recharts";
import SafeTextRenderer from "../components/SafeTextRenderer.tsx";
import Button from "../components/Button.tsx";

const LIMIT = 10;

type MetricType = "max" | "avg" | "volume";

const StatsPageGroupExercise: React.FC = () => {
    const [toast, setToast] = useState<string | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const {groupCode, exerciseID} = useParams();

    const [groupsMap, setGroupsMap] = useState<Record<string, Group>>({});
    const [exercisesMap, setExercisesMap] = useState<Record<number, ExerciseType>>({});

    const [stats, setStats] = useState<ExerciseStat[]>([]);
    const [total, setTotal] = useState<number>(0);
    const [offset, setOffset] = useState(0);
    const [hasMore, setHasMore] = useState(true);
    const [statsLoading, setStatsLoading] = useState(false);

    const [metric, setMetric] = useState<MetricType>("max");

    const loaderRef = useRef<HTMLDivElement | null>(null);

    // 🔥 фикс гонок
    const offsetRef = useRef(0);
    const isFetchingRef = useRef(false);

    // загрузка метаданных
    useEffect(() => {
        if (!groupCode) return;

        const fetchGroupData = async () => {
            try {
                setLoading(true);

                const [exerciseTypes, groups]: [ExerciseType[], Group[]] = await Promise.all([
                    getExerciseTypesByGroup(groupCode),
                    getExerciseGroups()
                ]);

                const groupsMap = groups.reduce<Record<string, Group>>((acc, group) => {
                    acc[group.code] = group;
                    return acc;
                }, {});

                const exercisesMap = exerciseTypes.reduce<Record<number, ExerciseType>>((acc, ex) => {
                    acc[ex.id] = ex;
                    return acc;
                }, {});

                setExercisesMap(exercisesMap);
                setGroupsMap(groupsMap);

            } catch (err: any) {
                setError(err.message || 'Не удалось загрузить данные');
            } finally {
                setLoading(false);
            }
        };

        fetchGroupData();
    }, [groupCode]);

    // загрузка статистики
    const loadStats = async () => {
        if (!exerciseID || isFetchingRef.current || !hasMore) return;

        isFetchingRef.current = true;
        setStatsLoading(true);

        try {
            const currentOffset = offsetRef.current;

            const res = await getExerciseStats(Number(exerciseID), currentOffset, LIMIT);
            const items = res.items || [];
            const total = res.total || 0;

            setStats(prev => {
                const existingIds = new Set(prev.map(s => s.id));
                const filtered = items.filter(s => !existingIds.has(s.id));
                return [...prev, ...filtered];
            });
            setTotal(total);

            offsetRef.current += LIMIT;
            setOffset(offsetRef.current);

            if (items.length < LIMIT) {
                setHasMore(false);
            }

        } catch (err: any) {
            setToast(err.message || "Ошибка загрузки статистики ❌");
        } finally {
            isFetchingRef.current = false;
            setStatsLoading(false);
        }
    };

    // при смене упражнения
    useEffect(() => {
        if (!exerciseID) return;

        setStats([]);
        setOffset(0);
        setHasMore(true);

        offsetRef.current = 0;
        isFetchingRef.current = false;

        loadStats();
    }, [exerciseID]);

    // infinite scroll
    useEffect(() => {
        if (!loaderRef.current) return;

        const observer = new IntersectionObserver(
            entries => {
                if (entries[0].isIntersecting) {
                    loadStats();
                }
            },
            {threshold: 1.0}
        );

        observer.observe(loaderRef.current);

        return () => observer.disconnect();
    }, [hasMore]);

    if (loading) return <Loader/>;
    if (error) return <p style={{color: 'red'}}>{error}</p>;

    const exerciseName = exercisesMap[Number(exerciseID)]?.name;

    // 📊 данные графика
    const chartData = stats
        .map(stat => {
            const sets = stat.sets || [];

            const weights = sets
                .map(s => s.fact_weight || s.weight || 0)
                .filter(w => w > 0);

            const maxWeight = weights.length ? Math.max(...weights) : 0;

            const avgWeight = weights.length
                ? weights.reduce((sum, w) => sum + w, 0) / weights.length
                : 0;

            const volume = sets.reduce((sum, s) => {
                const w = s.fact_weight || s.weight || 0;
                const r = s.fact_reps || s.reps || 0;
                return sum + w * r;
            }, 0);

            return {
                date: stat.date,
                maxWeight,
                avgWeight,
                volume
            };
        })
        .reverse();

    const metricConfig = {
        max: {key: "maxWeight", label: "Макс вес"},
        avg: {key: "avgWeight", label: "Средний вес"},
        volume: {key: "volume", label: "Объём"},
    };

    const currentMetric = metricConfig[metric];

    return (
        <div>
            <div className={"page stack"}>

                {groupCode && <h1>Динамика: {groupsMap[groupCode]?.name}</h1>}
                {exerciseName && (
                    <div style={{color: 'var(--color-text-muted)'}}>
                        <b>{exerciseName}</b>
                    </div>
                )}

                <div style={{display: "flex", gap: 4}}>
                    <Button variant="primary" onClick={() => setMetric("max")}>Макс</Button>
                    <Button variant="primary" onClick={() => setMetric("avg")}>Средний</Button>
                    <Button variant="primary" onClick={() => setMetric("volume")}>Объём</Button>
                </div>

                {chartData.length > 0 && (
                    <div style={{width: "100%", height: 300}}>
                        <ResponsiveContainer>
                            <LineChart data={chartData}>
                                <CartesianGrid strokeDasharray="3 3"/>
                                <XAxis dataKey="date"/>
                                <YAxis/>
                                <Tooltip/>
                                <Line
                                    type="monotone"
                                    dataKey={currentMetric.key}
                                    name={currentMetric.label}
                                />
                            </LineChart>
                        </ResponsiveContainer>
                    </div>
                )}

                <div className="stack">
                    {stats.map(stat => (
                        <div
                            key={stat.id}
                            className="card"
                            style={{
                                borderRadius: 12,
                                overflow: "hidden",
                                boxShadow: "0 2px 8px rgba(0,0,0,0.05)"
                            }}
                        >
                            <div
                                style={{
                                    display: "flex",
                                    alignItems: "center",
                                    gap: 8,
                                    padding: "8px 12px",
                                    background: "var(--color-primary-soft)",
                                    borderBottom: "1px solid rgba(0,0,0,0.05)"
                                }}
                            >
                                <CalendarRange size={18}/>
                                <b style={{fontSize: 14}}>
                                    {stat.date}
                                </b>
                            </div>

                            <div style={{padding: 10}}>
                                {stat.sets?.map((s, index) => (
                                    <div
                                        key={s.id}
                                        style={{
                                            padding: "6px 0",
                                            borderBottom: index !== stat.sets.length - 1
                                                ? "1px dashed rgba(0,0,0,0.06)"
                                                : "none"
                                        }}
                                    >
                                        <SafeTextRenderer html={s.formatted_string}/>
                                    </div>
                                ))}
                            </div>
                        </div>
                    ))}
                </div>

                <div ref={loaderRef} style={{height: 40}}>
                    {statsLoading && <Loader/>}
                </div>

                {!hasMore && (
                    <div style={{textAlign: "center"}}>
                        Всего: {total}
                    </div>
                )}
            </div>

            {toast && <Toast message={toast} onClose={() => setToast(null)}/>}
        </div>
    );
};

export default StatsPageGroupExercise;
