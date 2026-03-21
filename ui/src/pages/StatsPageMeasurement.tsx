import React, {useCallback, useEffect, useRef, useState} from "react";
import {CalendarRange, Loader} from "lucide-react";
import Toast from "../components/Toast.tsx";
import Button from "../components/Button.tsx";
import {CartesianGrid, Line, LineChart, ResponsiveContainer, Tooltip, XAxis, YAxis} from "recharts";
import {getMeasurements, getMeasurementTypes} from "../api/measurements.ts";

const LIMIT = 10;

// Добавляем массив цветов (можно расширить по количеству замеров)
const LINE_COLORS = [
    "#1f77b4", // синий
    "#ff7f0e", // оранжевый
    "#2ca02c", // зелёный
    "#d62728", // красный
    "#9467bd", // фиолетовый
    "#8c564b", // коричневый
    "#e377c2", // розовый
    "#7f7f7f", // серый
    "#bcbd22", // лаймовый
    "#17becf", // голубой
];

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

    // состояние видимых линий
    const [visibleLines, setVisibleLines] = useState<Record<string, boolean>>({});

    // инициализация после загрузки типов замеров
    useEffect(() => {
        const initialVisibility: Record<string, boolean> = {};
        measurementTypes.forEach(m => {
            initialVisibility[m.code] = true; // по умолчанию все линии видны
        });
        setVisibleLines(initialVisibility);
    }, [measurementTypes]);

    // функция переключения линии
    const toggleLine = (code: string) => {
        setVisibleLines(prev => ({
            ...prev,
            [code]: !prev[code],
        }));
    };

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

        window.addEventListener("scroll", handleScroll, {passive: true});
        return () => window.removeEventListener("scroll", handleScroll);
    }, [loadData]);

    // =========================
    // 📈 chart data preparation
    // =========================
    const chartData = data
        .map(item => {
            const entry: any = {date: item.created_at || "—"};

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

    if (loading) return <Loader/>;

    return (
        <div className="page stack">
            <h1>Динамика замеров</h1>

            {/* ========================= */}
            {/* 📈 chart */}
            {/* ========================= */}
            {chartData.length > 0 && (
                <div style={{width: "100%", height: 300, marginTop: 16}}>
                    <ResponsiveContainer>
                        <LineChart data={chartData}>
                            <CartesianGrid strokeDasharray="3 3"/>
                            <XAxis dataKey="date"/>
                            <YAxis/>
                            <Tooltip/>
                            {selectedCode === "all"
                                ? measurementTypes.map((m, idx) => visibleLines[m.code] && (
                                    <Line
                                        key={m.code}
                                        type="monotone"
                                        dataKey={m.code}
                                        stroke={LINE_COLORS[idx % LINE_COLORS.length]}
                                        strokeWidth={2}
                                    />
                                ))
                                : <Line type="monotone" dataKey={selectedCode} stroke={LINE_COLORS[0]} strokeWidth={2} />}
                        </LineChart>
                    </ResponsiveContainer>
                </div>
            )}

            {/* ========================= */}
            {/* 📊 interactive legend */}
            {/* ========================= */}
            <div style={{
                display: "flex",
                flexWrap: "wrap",
                gap: 8,
                marginTop: 8,
                justifyContent: "flex-start"
            }}>
                {selectedCode === "all"
                    ? measurementTypes.map((m, idx) => (
                        <Button
                            key={m.code}
                            variant={visibleLines[m.code] ? "ghost" : "ghost"}
                            style={{ display: "flex", alignItems: "center", gap: 4, minWidth: 60 }}
                            onClick={() => toggleLine(m.code)}
                        >
                            <div style={{
                                width: 12,
                                height: 12,
                                backgroundColor: LINE_COLORS[idx % LINE_COLORS.length],
                                borderRadius: 3
                            }} />
                            <span style={{ fontSize: 12 }}>{m.name}</span>
                        </Button>
                    ))
                    : (
                        <div style={{ display: "flex", alignItems: "center", gap: 4 }}>
                            <div
                                style={{
                                    width: 16,
                                    height: 16,
                                    backgroundColor: LINE_COLORS[0],
                                    borderRadius: 4
                                }}
                            />
                            <span>{measurementsMap[selectedCode]?.name || selectedCode}</span>
                        </div>
                    )
                }
            </div>

            {/* ========================= */}
            {/* 📋 list */}
            {/* ========================= */}
            <div className="stack" style={{marginTop: 16}}>
                {data.map(item => {
                    const codesToShow = selectedCode === "all"
                        ? measurementTypes.map(m => m.code)
                        : [selectedCode];

                    return (
                        <div key={item.id} className="card" style={{
                            padding: "12px 16px",
                            borderRadius: 8,
                            marginBottom: 8,
                            boxShadow: "0 1px 3px rgba(0,0,0,0.1)",
                            backgroundColor: "var(--color-bg)",
                            display: "flex",
                            flexDirection: "column",
                            gap: 8
                        }}>
                            {/* дата */}
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
                                    {item.created_at}
                                </b>
                            </div>

                            {/* значения замеров */}
                            <div style={{ display: "flex", flexWrap: "wrap", gap: 8 }}>
                                {codesToShow.map(code => {
                                    const value = getValue(item, code);
                                    if (value == null) return null;

                                    return (
                                        <div key={code} style={{
                                            backgroundColor: "var(--color-card)",
                                            padding: "4px 8px",
                                            borderRadius: 6,
                                            fontSize: 14,
                                            minWidth: 50,
                                            textAlign: "center"
                                        }}>
                                            <b style={{ fontWeight: 500 }}>{measurementsMap[code]?.name || code}</b>: {value}
                                        </div>
                                    );
                                })}
                            </div>
                        </div>
                    );
                })}
            </div>

            {dataLoading && <Loader/>}

            {!hasMore && (
                <div style={{textAlign: "center", marginTop: 16}}>
                    <b>Всего: {total}</b>
                </div>
            )}

            {toast && <Toast message={toast} onClose={() => setToast(null)}/>}
        </div>
    );
};

export default StatsPageMeasurement;
