import React, {useCallback, useEffect, useRef, useState} from 'react';
import './MeasurementsPage.css';
import Button from "../components/Button.tsx";

const PAGE_SIZE = 15;

const MeasurementsPage: React.FC = () => {
    const [measurements, setMeasurements] = useState<Measurement[]>([]);
    const [count, setCount] = useState<number>(0);
    const [offset, setOffset] = useState(0);
    const [loading, setLoading] = useState(false);
    const [hasMore, setHasMore] = useState(true);

    const [adding, setAdding] = useState(false);
    const [newMeasurement, setNewMeasurement] = useState<Partial<ToCreateMeasurement>>({});

    const tableObserver = useRef<IntersectionObserver | null>(null);
    const cardObserver = useRef<IntersectionObserver | null>(null);

    /* ================= infinite scroll для таблицы ================= */
    const lastRowRef = useCallback(
        (node: HTMLTableRowElement | null) => {
            if (loading || !hasMore) return;
            if (tableObserver.current) tableObserver.current.disconnect();

            tableObserver.current = new IntersectionObserver(entries => {
                if (entries[0].isIntersecting) {
                    setOffset(prev => prev + PAGE_SIZE);
                }
            });

            if (node) tableObserver.current.observe(node);
        },
        [loading, hasMore]
    );

    /* ================= infinite scroll для карточек ================= */
    const lastCardRef = useCallback(
        (node: HTMLDivElement | null) => {
            if (loading || !hasMore) return;
            if (cardObserver.current) cardObserver.current.disconnect();

            cardObserver.current = new IntersectionObserver(entries => {
                if (entries[0].isIntersecting) {
                    setOffset(prev => prev + PAGE_SIZE);
                }
            });

            if (node) cardObserver.current.observe(node);
        },
        [loading, hasMore]
    );

    /* ================= загрузка данных ================= */
    useEffect(() => {
        if (!hasMore) return;

        setLoading(true);

        fetch(`/api/measurements?offset=${offset}&limit=${PAGE_SIZE}`)
            .then(res => res.json())
            .then(data => {
                setMeasurements(prev => {
                    const ids = new Set(prev.map(m => m.id));
                    const unique = data.items.filter((m: Measurement) => !ids.has(m.id));
                    return [...prev, ...unique];
                });

                setCount(data.count);

                if (offset + data.items.length >= data.count) {
                    setHasMore(false);
                }
            })
            .finally(() => setLoading(false));
    }, [offset, hasMore]);

    /* ================= обработка формы ================= */
    const handleInputChange = (field: keyof ToCreateMeasurement, value: string) => {
        setNewMeasurement(prev => ({
            ...prev,
            [field]: Number(value)
        }));
    };

    const handleSaveNewMeasurement = () => {
        fetch('/api/measurements', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify(newMeasurement),
        })
            .then(res => res.json())
            .then(data => {
                setMeasurements(prev => [data, ...prev]);
                setCount(prev => prev + 1);
                setAdding(false);
                setNewMeasurement({});
            });
    };

    const handleCancelNewMeasurement = () => {
        setAdding(false);
        setNewMeasurement({});
    };

    return (
        <div className="measurements-page">
            <h1>Замеры</h1>

            {/* ===== DESKTOP BUTTON ===== */}
            {!adding && (
                <Button onClick={() => setAdding(true)}>
                    ➕ Добавить новое измерение
                </Button>
            )}

            {/* ================= TABLE (DESKTOP) ================= */}
            <div className="desktop-only table-wrapper">
                <table className="measurements-table">
                    <thead>
                    <tr>
                        <th>№</th>
                        <th>Дата</th>
                        <th>Плечи</th>
                        <th>Грудь</th>
                        <th>Л. рука</th>
                        <th>П. рука</th>
                        <th>Талия</th>
                        <th>Ягодицы</th>
                        <th>Л. бедро</th>
                        <th>П. бедро</th>
                        <th>Л. икра</th>
                        <th>П. икра</th>
                        <th>Вес</th>
                        <th>Действия</th>
                    </tr>
                    </thead>

                    <tbody>
                    {adding && (
                        <tr className="new-measurement-row">
                            <td>-</td>
                            <td>-</td>
                            {fields.map(field => (
                                <td key={field.key}>
                                    <input
                                        type="number"
                                        value={newMeasurement[field.key] ?? ''}
                                        onChange={e => handleInputChange(field.key, e.target.value)}
                                    />
                                </td>
                            ))}
                            <td>
                                <button onClick={handleSaveNewMeasurement}>Сохранить</button>
                                <button onClick={handleCancelNewMeasurement}>Отмена</button>
                            </td>
                        </tr>
                    )}

                    {measurements.map((m, idx) => (
                        <tr
                            key={m.id}
                            ref={idx === measurements.length - 1 ? lastRowRef : null}
                        >
                            <td>{idx + 1}</td>
                            <td>{m.created_at}</td>
                            <td>{m.shoulders}</td>
                            <td>{m.chest}</td>
                            <td>{m.hand_left}</td>
                            <td>{m.hand_right}</td>
                            <td>{m.waist}</td>
                            <td>{m.buttocks}</td>
                            <td>{m.hip_left}</td>
                            <td>{m.hip_right}</td>
                            <td>{m.calf_left}</td>
                            <td>{m.calf_right}</td>
                            <td className="weight">{m.weight}</td>
                            <td/>
                        </tr>
                    ))}

                    {loading && (
                        <tr>
                            <td colSpan={14}>Загрузка...</td>
                        </tr>
                    )}
                    </tbody>
                </table>
            </div>

            {/* ================= CARDS (MOBILE) ================= */}
            <div className="mobile-only cards-wrapper">
                {adding && (
                    <div className="card-form">
                        {fields.map(field => (
                            <div key={field.key} className="card-form-field">
                                <label>{field.label}</label>
                                <input
                                    type="number"
                                    value={newMeasurement[field.key] ?? ''}
                                    onChange={e => handleInputChange(field.key, e.target.value)}
                                />
                            </div>
                        ))}
                        <div className="card-form-buttons">
                            <button onClick={handleSaveNewMeasurement}>Сохранить</button>
                            <button onClick={handleCancelNewMeasurement}>Отмена</button>
                        </div>
                    </div>
                )}

                {measurements.map((m, idx) => (
                    <div
                        key={m.id}
                        ref={idx === measurements.length - 1 ? lastCardRef : null}
                        className="measurement-card"
                    >
                        <div className="card-header">
                            <span>📅 {m.created_at}</span>
                            <span>⚖ {m.weight} кг</span>
                        </div>

                        <div className="card-body two-columns">
                            {/* Левый столбец */}
                            <div className="card-column">
                                <div className="card-row"><span>Плечи:</span><span>{m.shoulders}</span></div>
                                <div className="card-row"><span>Грудь:</span><span>{m.chest}</span></div>
                                <div className="card-row"><span>Л. рука:</span><span>{m.hand_left}</span></div>
                                <div className="card-row"><span>П. рука:</span><span>{m.hand_right}</span></div>
                                <div className="card-row"><span>Талия:</span><span>{m.waist}</span></div>
                            </div>

                            {/* Правый столбец */}
                            <div className="card-column">
                                <div className="card-row"><span>Ягодицы:</span><span>{m.buttocks}</span></div>
                                <div className="card-row"><span>Л. бедро:</span><span>{m.hip_left}</span></div>
                                <div className="card-row"><span>П. бедро:</span><span>{m.hip_right}</span></div>
                                <div className="card-row"><span>Л. икра:</span><span>{m.calf_left}</span></div>
                                <div className="card-row"><span>П. икра:</span><span>{m.calf_right}</span></div>
                            </div>
                        </div>
                    </div>
                ))}
            </div>

        </div>
    );
};

const fields: { key: keyof ToCreateMeasurement; label: string }[] = [
    {key: 'shoulders', label: 'Плечи'},
    {key: 'chest', label: 'Грудь'},
    {key: 'hand_left', label: 'Л. рука'},
    {key: 'hand_right', label: 'П. рука'},
    {key: 'waist', label: 'Талия'},
    {key: 'buttocks', label: 'Ягодицы'},
    {key: 'hip_left', label: 'Л. бедро'},
    {key: 'hip_right', label: 'П. бедро'},
    {key: 'calf_left', label: 'Л. икра'},
    {key: 'calf_right', label: 'П. икра'},
    {key: 'weight', label: 'Вес'},
];

export default MeasurementsPage;
