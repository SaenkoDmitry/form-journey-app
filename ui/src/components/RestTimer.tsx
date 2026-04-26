import { useEffect, useRef, useState } from "react";
import { createPortal } from "react-dom";
import { useRestTimer } from "../context/RestTimerContext";
import Button from "./Button";
import "../styles/RestTimer.css";
import { Pause, Play, RotateCcw, ChevronDown } from "lucide-react";
import { toast } from "react-hot-toast";
import { startTimer } from "../api/timers.ts";

type Props = {
    seconds: number;
    autoStartTrigger?: number;
    workoutID?: number;
};

const PRESETS = [30, 60, 90, 120, 180];

export default function RestTimer({ seconds, autoStartTrigger, workoutID }: Props) {
    const {
        remaining,
        running,
        start: localStart,
        pause,
        reset,
        seconds: totalSeconds
    } = useRestTimer();

    const [customSeconds, setCustomSeconds] = useState(seconds);
    const [showEditor, setShowEditor] = useState(false);
    const inputRef = useRef<HTMLInputElement>(null);
    const sheetRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        setCustomSeconds(seconds);
    }, [seconds]);

    // 🔥 Закрытие по клику вне sheet + блокировка click-through
    useEffect(() => {
        if (!showEditor) return;

        const handlePointerDown = (e: PointerEvent) => {
            // Если клик вне sheet — закрываем и блокируем дальнейшую обработку
            if (sheetRef.current && !sheetRef.current.contains(e.target as Node)) {
                e.preventDefault();
                e.stopPropagation();
                setShowEditor(false);
            }
        };

        // 🔥 capture: true — ловим ДО того, как событие дойдёт до кнопок под overlay
        document.addEventListener("pointerdown", handlePointerDown, true);

        return () => {
            document.removeEventListener("pointerdown", handlePointerDown, true);
        };
    }, [showEditor]);

    // Закрытие по Escape
    useEffect(() => {
        if (!showEditor) return;

        const handleKey = (e: KeyboardEvent) => {
            if (e.key === "Escape") setShowEditor(false);
        };
        document.addEventListener("keydown", handleKey);
        return () => document.removeEventListener("keydown", handleKey);
    }, [showEditor]);

    const start = async (secs: number) => {
        if (secs <= 0) return;
        if (!workoutID) return;

        try {
            const resp = await startTimer(workoutID, secs);
            if (resp?.id != null) {
                localStorage.setItem("currentTimerID", resp.id.toString());
            }
            localStart(secs, resp.id);
        } catch (err) {
            console.error("Failed to start server timer", err);
            toast.error("Не удалось зарегистрировать таймер на сервере");
        }
    };

    useEffect(() => {
        if (!autoStartTrigger || !workoutID) return;
        localStorage.setItem("floatingTimerWorkoutID", workoutID.toString());
        start(customSeconds);
    }, [autoStartTrigger, workoutID]);

    useEffect(() => {
        if (remaining === 0 && running) {
            navigator.vibrate?.([300, 150, 300]);
            toast.success("Таймер завершён!");
            localStorage.removeItem("floatingTimerWorkoutID");
        }
    }, [remaining, running]);

    const format = (t: number) => {
        const m = Math.floor(t / 60);
        const s = t % 60;
        return `${m}:${s.toString().padStart(2, "0")}`;
    };

    const progress = totalSeconds > 0 ? 1 - remaining / totalSeconds : 0;
    const radius = 28;
    const circumference = 2 * Math.PI * radius;

    const applyPreset = (s: number) => {
        setCustomSeconds(s);
        setShowEditor(false);
        reset();
    };

    const adjustBy = (delta: number) => {
        setCustomSeconds(prev => Math.max(5, prev + delta));
    };

    const sheetContent = showEditor && (
        <div className="rest-sheet-overlay">
            <div className="rest-sheet" ref={sheetRef}>
                <div className="sheet-handle" />

                <div className="sheet-header">
                    <h3>Время отдыха</h3>
                    <button className="sheet-done" onClick={() => setShowEditor(false)}>
                        Готово
                    </button>
                </div>

                <div className="presets">
                    {PRESETS.map(s => (
                        <button
                            key={s}
                            className={`preset-chip ${customSeconds === s ? "active" : ""}`}
                            onClick={() => applyPreset(s)}
                        >
                            {s < 60 ? `${s}с` : `${s / 60}м`}
                        </button>
                    ))}
                </div>

                <div className="fine-tune">
                    <button className="tune-btn" onClick={() => adjustBy(-5)}>−5</button>

                    <div className="tune-input-wrap">
                        <input
                            ref={inputRef}
                            type="number"
                            inputMode="numeric"
                            pattern="[0-9]*"
                            value={customSeconds}
                            onChange={(e) => {
                                const val = parseInt(e.target.value) || 0;
                                setCustomSeconds(Math.max(5, val));
                            }}
                        />
                        <span>сек</span>
                    </div>

                    <button className="tune-btn" onClick={() => adjustBy(5)}>+5</button>
                </div>

                <button
                    className="sheet-start-btn"
                    onClick={() => {
                        setShowEditor(false);
                        reset();
                        start(customSeconds);
                    }}
                >
                    Старт {customSeconds}с
                </button>
            </div>
        </div>
    );

    return (
        <div className={`rest-timer ${running ? "active" : ""}`}>
            <div className="timer-inner">
                <div className="timer-top-row">
                    <div className="circle">
                        <svg width="60" height="60" viewBox="0 0 60 60">
                            <circle className="bg" strokeWidth="8" r={radius} cx="30" cy="30" />
                            <circle
                                className="progress"
                                strokeWidth="8"
                                r={radius}
                                cx="30"
                                cy="30"
                                strokeDasharray={circumference}
                                strokeDashoffset={circumference * (1 - progress)}
                            />
                        </svg>
                        <div className="time">{format(remaining)}</div>
                    </div>

                    <div className="actions">
                        <Button
                            variant={running ? "primary" : "active"}
                            onClick={() => {
                                if (running) {
                                    pause();
                                } else if (remaining > 0) {
                                    start(remaining);
                                } else {
                                    localStorage.setItem("floatingTimerWorkoutID", workoutID?.toString() ?? "");
                                    start(customSeconds);
                                }
                            }}
                        >
                            {running ? <Pause size={14} /> : <Play size={14} />}
                            {running ? "Пауза" : "Старт"}
                        </Button>

                        <Button variant="ghost" onClick={reset}>
                            <RotateCcw size={14} /> Сброс
                        </Button>
                    </div>
                </div>

                <button
                    className="rest-value-btn"
                    onClick={() => {
                        if (!running) {
                            setShowEditor(true);
                            setTimeout(() => inputRef.current?.focus(), 350);
                        }
                    }}
                    disabled={running}
                >
                    <span className="rest-label">Отдых</span>
                    <span className="rest-time">{customSeconds} сек</span>
                    {!running && <ChevronDown size={14} />}
                </button>
            </div>

            {sheetContent && createPortal(sheetContent, document.body)}
        </div>
    );
}
