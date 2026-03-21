import React, { useCallback, useEffect, useRef, useState } from "react";
import { Loader } from "lucide-react";
import Toast from "../components/Toast.tsx";
import Button from "../components/Button.tsx";
import {
    CartesianGrid,
    Line,
    LineChart,
    ResponsiveContainer,
    Tooltip,
    XAxis,
    YAxis
} from "recharts";
import { getMeasurements, getMeasurementTypes } from "../api/measurements.ts";

const LIMIT = 10;

const StatsPageMeasurement: React.FC = () => {
    const [toast, setToast] = useState<string | null>(null);
    const [loading, setLoading] = useState(true);
    const [measurementTypes, setMeasurementTypes] = useState<MeasurementTypeDTO[]>([]);
    const [measurementsMap, setMeasurementsMap] = useState<Record<string, MeasurementTypeDTO>>({});
    const [data, setData] = useState<Measurement[]>([]);
    const [total, setTotal] = useState(0);
    const [hasMore, setHasMore] = useState(true);
    const [dataLoading, setDataLoading] = useState(false);
    const [selectedCode, setSelectedCode] = useState<string | "all">("all");

    const offsetRef = useRef(0);
    const isFetchingRef = useRef(false);
    const hasMoreRef = useRef(true);

    useEffect(() => {
        hasMoreRef.current = hasMore;
    }, [hasMore]);

    // =========================
    // 📦 load measurement types
    // =========================
    useEffect(() => {
        (async () => {
            try {
                setLoading(true);
                const types = await getMeasurementTypes();
                setMeasurementTypes(types);

                const map = types.reduce<Record<string, MeasurementTypeDTO>>((acc, m) => {
                    acc[m.code] = m;
                    return acc;
                }, {});
                setMeasurementsMap(map);
            } catch (err: any) {
                setToast("Ошибка загрузки видов измерений ❌");
            } finally {
                setLoading(false);
            }
        })();
    }, []);

    // =========================
    // 🔍 get value by code
    // =========================
    const getValue = (item: any, code: string) => {
        if (item[code] != null) return item[code];
        if (item.measurements?.[code] != null) return item.measurements[code];
        return null;
    };

    // =========================
    // 📊 load data with pagination
    // =========================
    const loadData = useCallback(async () => {
        if (isFetchingRef.current || !hasMoreRef.current) return;

        isFetchingRef.current = true;
        setDataLoading(true);

        try {
            let collected: any[] = [];

            while (collected.length < LIMIT && hasMoreRef.current) {
                const res = await getMeasurements(offsetRef.current, LIMIT);
                const items = res.items || [];

                if (items.length === 0) {
                    setHasMore(false);
                    hasMoreRef.current = false;
                    break;
                }

                offsetRef.current += items.length;
                setTotal(res.count || 0);

                collected.push(...items);

                if (offsetRef.current >= res.total) {
                    setHasMore(false);
                    hasMoreRef.current = false;
                    break;
                }
            }

            if (collected.length === 0) return;

            setData(prev => {
                const ids = new Set(prev.map(i => i.id));
                const unique = collected.filter(i => !ids.has(i.id));
                return [...prev, ...unique];
            });

        } catch (err: any) {
            setToast(err.message || "Ошибка загрузки ❌");
        } finally {
            isFetchingRef.current = false;
            setDataLoading(false);
        }
    }, []);

    // =========================
    // 🔄 reset on filter change
    // =========================
    useEffect(() => {
        setData([]);
        setHasMore(true);
        offsetRef.current = 0;
        isFetchingRef.current = false;
        hasMoreRef.current = true;

        loadData();
    }, [loadData, selectedCode]);

    // =========================
    // 📜 infinite scroll
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

                if (scrollTop + windowHeight >= fullHeight - 300) {
                    loadData();
                }

                ticking = false;
            });
        };

        window.addEventListener("scroll", handleScroll, { passive: true });
        return () => window.removeEventListener("scroll", handleScroll);
    }, [loadData]);

    // =========================
    // 📈 chart data preparation
    // =========================
    const chartData = data
        .map(item => {
            const entry: any = { date: item.created_at || "—" };

            if (selectedCode === "all") {
                measurementTypes.forEach(m => {
                    entry[m.code] = getValue(item, m.code);
                });
            } else {
                entry[selectedCode] = getValue(item, selectedCode);
            }

            return entry;
        })
        .filter(item => {
            if (selectedCode === "all") {
                return measurementTypes.some(m => item[m.code] != null);
            } else {
                return item[selectedCode] != null;
            }
        })
        .reverse();

    if (loading) return <Loader />;

    return (
        <div className="page stack">
            <h1>Динамика замеров</h1>

            {/* ========================= */}
            {/* 🔘 filter / legend */}
            {/* ========================= */}
            <div style={{ display: "flex", gap: 8, flexWrap: "wrap" }}>
                <Button
                    variant={selectedCode === "all" ? "active" : "ghost"}
                    onClick={() => setSelectedCode("all")}
                >
                    Все
                </Button>
                {measurementTypes.map(m => (
                    <Button
                        key={m.code}
                        variant={selectedCode === m.code ? "active" : "ghost"}
                        onClick={() => setSelectedCode(m.code)}
                    >
                        {m.name}
                    </Button>
                ))}
            </div>

            {/* ========================= */}
            {/* 📈 chart */}
            {/* ========================= */}
            {chartData.length > 0 && (
                <div style={{ width: "100%", height: 300, marginTop: 16 }}>
                    <ResponsiveContainer>
                        <LineChart data={chartData}>
                            <CartesianGrid strokeDasharray="3 3" />
                            <XAxis dataKey="date" />
                            <YAxis />
                            <Tooltip />
                            {selectedCode === "all"
                                ? measurementTypes.map(m => (
                                    <Line
                                        key={m.code}
                                        type="monotone"
                                        dataKey={m.code}
                                        strokeWidth={2}
                                    />
                                ))
                                : <Line type="monotone" dataKey={selectedCode} strokeWidth={2} />}
                        </LineChart>
                    </ResponsiveContainer>
                </div>
            )}

            {/* ========================= */}
            {/* 📋 list */}
            {/* ========================= */}
            <div className="stack" style={{ marginTop: 16 }}>
                {data.map(item => {
                    const codesToShow = selectedCode === "all"
                        ? measurementTypes.map(m => m.code)
                        : [selectedCode];

                    return (
                        <div key={item.id} className="card">
                            <b>{item.created_at}</b>
                            {codesToShow.map(code => {
                                const value = getValue(item, code);
                                if (value == null) return null;
                                return (
                                    <div key={code}>
                                        {measurementsMap[code]?.name || code}: {value}
                                    </div>
                                );
                            })}
                        </div>
                    );
                })}
            </div>

            {dataLoading && <Loader />}

            {!hasMore && (
                <div style={{ textAlign: "center", marginTop: 16 }}>
                    <b>Всего: {total}</b>
                </div>
            )}

            {toast && <Toast message={toast} onClose={() => setToast(null)} />}
        </div>
    );
};

export default StatsPageMeasurement;
