import React, {useCallback, useEffect, useRef, useState} from 'react';
import '../styles/MeasurementsPage.css';
import Button from "../components/Button.tsx";
import {deleteMeasurement as apiDeleteMeasurement} from "../api/measurements.ts";
import {useAuth} from "../context/AuthContext.tsx";
import {useNavigate} from "react-router-dom";

const PAGE_SIZE = 15;

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

const MeasurementsPage: React.FC = () => {
    const {user, loading: authLoading} = useAuth();
    const navigate = useNavigate();
    const [measurements, setMeasurements] = useState<Measurement[]>([]);
    const [count, setCount] = useState(0);
    const [offset, setOffset] = useState(0);
    const [loading, setLoading] = useState(false);
    const [hasMore, setHasMore] = useState(true);

    const [toast, setToast] = useState<string | null>(null);

    const [adding, setAdding] = useState(false);
    const [newMeasurement, setNewMeasurement] = useState<Partial<ToCreateMeasurement>>({});

    const tableObserver = useRef<IntersectionObserver | null>(null);
    const cardObserver = useRef<IntersectionObserver | null>(null);

    /* ===================== Центральный fetch ===================== */
    const fetchMeasurements = async (offset: number) => {
        try {
            setLoading(true);
            const res = await fetch(`/api/measurements?offset=${offset}&limit=${PAGE_SIZE}`);
            if (!res.ok) throw new Error("Ошибка при загрузке измерений");
            const data = await res.json() as { items: Measurement[]; count: number };

            setMeasurements(prev => {
                const ids = new Set(prev.map(m => m.id));
                const unique = data.items.filter(m => !ids.has(m.id));
                return [...prev, ...unique];
            });

            setCount(data.count);

            if (offset + data.items.length >= data.count) {
                setHasMore(false);
            }
        } catch {
            showToast("❌ Ошибка при загрузке измерений");
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        if (hasMore) fetchMeasurements(offset);
    }, [offset, hasMore]);

    /* ===================== Infinite scroll ===================== */
    const lastRowRef = useCallback((node: HTMLTableRowElement | null) => {
        if (loading || !hasMore) return;
        if (tableObserver.current) tableObserver.current.disconnect();

        tableObserver.current = new IntersectionObserver(entries => {
            if (entries[0].isIntersecting) setOffset(prev => prev + PAGE_SIZE);
        });

        if (node) tableObserver.current.observe(node);
    }, [loading, hasMore]);

    const lastCardRef = useCallback((node: HTMLDivElement | null) => {
        if (loading || !hasMore) return;
        if (cardObserver.current) cardObserver.current.disconnect();

        cardObserver.current = new IntersectionObserver(entries => {
            if (entries[0].isIntersecting) setOffset(prev => prev + PAGE_SIZE);
        });

        if (node) cardObserver.current.observe(node);
    }, [loading, hasMore]);

    /* ===================== Toast ===================== */
    const showToast = (text: string) => {
        setToast(text);
        setTimeout(() => setToast(null), 3000);
    };

    /* ===================== Добавление измерения ===================== */
    const handleInputChange = (field: keyof ToCreateMeasurement, value: string) => {
        setNewMeasurement(prev => ({...prev, [field]: Number(value)}));
    };

    const handleSaveNewMeasurement = async () => {
        try {
            const res = await fetch('/api/measurements', {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify(newMeasurement),
            });
            if (!res.ok) throw new Error("Ошибка при сохранении измерения");
            const data: Measurement = await res.json();

            setMeasurements(prev => [data, ...prev]);
            setCount(prev => prev + 1);
            setAdding(false);
            setNewMeasurement({});
            showToast("✅ Измерение добавлено");
        } catch {
            showToast("❌ Ошибка при добавлении измерения");
        }
    };

    const handleCancelNewMeasurement = () => {
        setAdding(false);
        setNewMeasurement({});
    };

    /* ===================== Удаление измерения ===================== */
    const handleDeleteMeasurement = async (id: number) => {
        const confirmed = window.confirm("Вы уверены, что хотите удалить измерение?");
        if (!confirmed) return;

        try {
            await apiDeleteMeasurement(id);
            setMeasurements(prev => prev.filter(m => m.id !== id));
            setCount(prev => prev - 1);
            showToast("✅ Измерение удалено");
        } catch {
            showToast("❌ Ошибка при удалении измерения");
        }
    };

    return <div className="measurements-page">
        <h1>Замеры</h1>

        {!adding && (
            <Button variant="active" onClick={() => setAdding(true)}>
                ➕ Добавить новое измерение
            </Button>
        )}

        {/* ======== TABLE (DESKTOP) ======== */}
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
                        <td style={{minWidth: "140px"}}>
                            <Button variant="active" onClick={handleSaveNewMeasurement}>💾</Button>
                            <Button style={{marginLeft: "10px"}} onClick={handleCancelNewMeasurement}>❌</Button>
                        </td>
                    </tr>
                )}

                {measurements.map((m, idx) => (
                    <tr key={m.id} ref={idx === measurements.length - 1 ? lastRowRef : null}>
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
                        <td>
                            <Button
                                style={{width: "44px", padding: "0 8px", borderRadius: "16px"}}
                                variant="danger"
                                onClick={() => handleDeleteMeasurement(m.id)}
                                disabled={!m.id} // <-- блокируем, если id = 0 или undefined
                            >🗑</Button>
                        </td>
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

        {/* ======== CARDS (MOBILE) ======== */}
        <div className="mobile-only cards-wrapper">
            {adding && (
                <div className="card-form">
                    {fields.map(field => (
                        <div key={field.key} className="card-form-field">
                            <b>{field.label}</b>
                            <input
                                style={{maxHeight: "40px"}}
                                type="number"
                                value={newMeasurement[field.key] ?? ''}
                                onChange={e => handleInputChange(field.key, e.target.value)}
                            />
                        </div>
                    ))}
                    <div className="card-form-buttons">
                        <Button variant="active" onClick={handleSaveNewMeasurement}>Сохранить</Button>
                        <Button variant="ghost" onClick={handleCancelNewMeasurement}>Отмена</Button>
                    </div>
                </div>
            )}

            {measurements.map((m, idx) => (
                <div
                    key={m.id}
                    ref={idx === measurements.length - 1 ? lastCardRef : null}
                    className="measurement-card"
                >
                    <div className="card-header"
                         style={{display: 'flex', justifyContent: 'space-between', alignItems: 'center'}}>
                        <div>
                            <span>📅 {m.created_at}</span>
                            <span style={{marginLeft: '10px'}}>⚖ {m.weight} кг</span>
                        </div>
                        <Button
                            style={{width: '44px', height: '44px'}}
                            variant="danger"
                            onClick={() => handleDeleteMeasurement(m.id)}
                            disabled={!m.id}
                        >
                            🗑
                        </Button>
                    </div>

                    <div className="card-body two-columns">
                        <div className="card-column">
                            <div className="card-row"><span>Плечи:</span><span>{m.shoulders}</span></div>
                            <div className="card-row"><span>Грудь:</span><span>{m.chest}</span></div>
                            <div className="card-row"><span>Л. рука:</span><span>{m.hand_left}</span></div>
                            <div className="card-row"><span>П. рука:</span><span>{m.hand_right}</span></div>
                            <div className="card-row"><span>Талия:</span><span>{m.waist}</span></div>
                        </div>
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

        {toast && <div className="toast">{toast}</div>}
    </div>;
};

export default MeasurementsPage;
