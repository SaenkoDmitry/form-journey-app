import {CalendarRange, Loader} from "lucide-react";
import React, {useCallback, useEffect, useRef, useState} from "react";
import Toast from "../components/Toast.tsx";
import {getExerciseGroups, getExerciseStats, getExerciseTypesByGroup,} from "../api/exercises.ts";
import {useParams} from "react-router-dom";

import {CartesianGrid, Line, LineChart, ResponsiveContainer, Tooltip, XAxis, YAxis} from "recharts";

import SafeTextRenderer from "../components/SafeTextRenderer.tsx";
import Button from "../components/Button.tsx";

const LIMIT = 10;

type MetricType = "max" | "avg" | "volume";

const StatsPageGroupExercise: React.FC = () => {
    const {groupCode, exerciseID} = useParams();

    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [toast, setToast] = useState<string | null>(null);

    const [groupsMap, setGroupsMap] = useState<Record<string, Group>>({});
    const [exercisesMap, setExercisesMap] = useState<Record<number, ExerciseType>>({});

    const [stats, setStats] = useState<ExerciseStat[]>([]);
    const [total, setTotal] = useState(0);

    const [hasMore, setHasMore] = useState(true);
    const [statsLoading, setStatsLoading] = useState(false);

    const [metric, setMetric] = useState<MetricType>("max");

    const offsetRef = useRef(0);
    const isFetchingRef = useRef(false);
    const hasMoreRef = useRef(true);

    // синхронизация ref
    useEffect(() => {
        hasMoreRef.current = hasMore;
    }, [hasMore]);

    // =========================
    // 📦 metadata
    // =========================
    useEffect(() => {
        if (!groupCode) return;

        (async () => {
            try {
                setLoading(true);

                const [exerciseTypes, groups] = await Promise.all([
                    getExerciseTypesByGroup(groupCode),
                    getExerciseGroups()
                ]);

                setGroupsMap(
                    Object.fromEntries(groups.map(g => [g.code, g]))
                );

                setExercisesMap(
                    Object.fromEntries(exerciseTypes.map(e => [e.id, e]))
                );

            } catch (err: any) {
                setError(err.message || "Ошибка загрузки");
            } finally {
                setLoading(false);
            }
        })();
    }, [groupCode]);

    // =========================
    // 📊 load stats
    // =========================
    const loadStats = useCallback(async () => {
        if (!exerciseID || isFetchingRef.current || !hasMoreRef.current) return;

        isFetchingRef.current = true;
        setStatsLoading(true);

        try {
            const currentOffset = offsetRef.current;

            const res = await getExerciseStats(
                Number(exerciseID),
                currentOffset,
                LIMIT
            );

            const items = res.items || [];
            const newTotal = res.total || 0;

            // стоп если пусто
            if (items.length === 0) {
                setHasMore(false);
                hasMoreRef.current = false;
                return;
            }

            setStats(prev => {
                const ids = new Set(prev.map(s => s.id));
                const filtered = items.filter(s => !ids.has(s.id));
                return [...prev, ...filtered];
            });

            const newOffset = currentOffset + items.length;

            offsetRef.current = newOffset;
            setTotal(newTotal);

            const more = newOffset < newTotal;
            setHasMore(more);
            hasMoreRef.current = more;

        } catch (err: any) {
            setToast(err.message || "Ошибка загрузки ❌");
        } finally {
            isFetchingRef.current = false;
            setStatsLoading(false);
        }
    }, [exerciseID]);

    // =========================
    // 🔄 reset
    // =========================
    useEffect(() => {
        if (!exerciseID) return;

        setStats([]);
        setHasMore(true);

        offsetRef.current = 0;
        isFetchingRef.current = false;
        hasMoreRef.current = true;

        loadStats();
    }, [exerciseID, loadStats]);

    // =========================
    // 📜 SCROLL LISTENER (🔥 главное)
    // =========================
    useEffect(() => {
        let ticking = false;

        const handleScroll = () => {
            if (ticking) return;

            ticking = true;

            requestAnimationFrame(() => {
                const scrollTop = window.scrollY;
                const windowHeight = window.innerHeight;
                const fullHeight = document.documentElement.scrollHeight;

                // триггер за 300px до низа
                if (scrollTop + windowHeight >= fullHeight - 300) {
                    loadStats();
                }

                ticking = false;
            });
        };

        window.addEventListener("scroll", handleScroll, {passive: true});

        return () => {
            window.removeEventListener("scroll", handleScroll);
        };
    }, [loadStats]);

    function calcMaxWeight(sets: FormattedSet[]): number {
        const weights = sets.map(s => s.fact_weight || s.weight || 0).filter(Boolean);
        return weights.length ? Math.max(...weights) : 0;
    }

    function calcAvgWeight(sets: FormattedSet[]): number {
        const weights = sets.map(s => s.fact_weight || s.weight || 0).filter(Boolean);
        return weights.length ? weights.reduce((a, b) => a + b, 0) / weights.length : 0;
    }

    function calcVolume(sets: FormattedSet[]): number {
        return sets.reduce((sum, s) => {
            const w = s.fact_weight || s.weight || 0;
            const r = s.fact_reps || s.reps || 0;
            return sum + w * r;
        }, 0)
    }

    // =========================
    // 📈 chart data
    // =========================
    const chartData = stats.map(stat => {
        const sets: FormattedSet[] = stat.sets || [];

        const maxWeight = calcMaxWeight(sets);
        const avgWeight = calcAvgWeight(sets);
        const volume = calcVolume(sets);

        return {
            date: stat.date,
            maxWeight,
            avgWeight,
            volume
        };
    }).reverse();

    const metricMap = {
        max: {key: "maxWeight", label: "Макс вес"},
        avg: {key: "avgWeight", label: "Средний вес"},
        volume: {key: "volume", label: "Объём"}
    };

    const currentMetric = metricMap[metric];

    if (loading) return <Loader/>;
    if (error) return <p style={{color: "red"}}>{error}</p>;

    const exerciseName = exercisesMap[Number(exerciseID)]?.name;

    return (
        <div className="page stack">

            <h1>Динамика: {groupsMap[groupCode!]?.name}</h1>

            {exerciseName && (
                <div style={{color: "var(--color-text-muted)"}}>
                    <b>{exerciseName}</b>
                </div>
            )}

            <div style={{display: "flex", gap: 6}}>
                {(["max", "avg", "volume"] as MetricType[]).map(m => (
                    <Button key={m} variant={metric === m ? "active" : "ghost"} onClick={() => setMetric(m)}>
                        {metricMap[m].label}
                    </Button>
                ))}
            </div>

            {chartData.length > 0 && (
                <div style={{width: "100%", height: 300}}>
                    <ResponsiveContainer>
                        <LineChart data={chartData}>
                            <CartesianGrid strokeDasharray="3 3"/>
                            <XAxis dataKey="date"/>
                            <YAxis/>
                            <Tooltip/>
                            <Line dataKey={currentMetric.key}/>
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
                                borderRadius: 12
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

                            <hr/>
                            <div style={{paddingTop: 8}}>
                                Макс: {calcMaxWeight(stat.sets || [])} |
                                Средний: {calcAvgWeight(stat.sets || [])} |
                                Объем: {calcVolume(stat.sets || [])}
                            </div>
                        </div>
                    </div>
                ))}
            </div>

            {statsLoading && <Loader/>}

            {!hasMore && (
                <div style={{textAlign: "center"}}>
                    <b>Всего: {total}</b>
                </div>
            )}

            {toast && (
                <Toast message={toast} onClose={() => setToast(null)}/>
            )}
        </div>
    );
};

export default StatsPageGroupExercise;
