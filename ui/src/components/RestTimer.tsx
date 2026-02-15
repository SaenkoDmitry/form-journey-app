import { useEffect, useRef, useState } from "react";
import Button from "./Button";
import "../styles/RestTimer.css";

const STORAGE_KEY = "rest_timer_end";

export default function RestTimer({
                                      seconds,
                                      onFinish,
                                      autoStartTrigger,
                                  }: {
    seconds: number;
    onFinish?: () => void;
    autoStartTrigger?: number;
}) {
    const [endTime, setEndTime] = useState<number | null>(null);
    const [remaining, setRemaining] = useState(seconds);
    const [running, setRunning] = useState(false);

    const intervalRef = useRef<number | null>(null);

    // восстановление
    useEffect(() => {
        const saved = localStorage.getItem(STORAGE_KEY);
        if (saved) {
            const parsed = Number(saved);
            if (parsed > Date.now()) {
                setEndTime(parsed);
                setRunning(true);
            } else {
                localStorage.removeItem(STORAGE_KEY);
            }
        }
    }, []);

    // автостарт
    useEffect(() => {
        if (!autoStartTrigger) return;
        start();
    }, [autoStartTrigger]);

    useEffect(() => {
        if (!running || !endTime) return;

        intervalRef.current = window.setInterval(() => {
            const diff = Math.max(0, Math.floor((endTime - Date.now()) / 1000));
            setRemaining(diff);

            if (diff <= 0) finish();
        }, 500);

        return () => {
            if (intervalRef.current) clearInterval(intervalRef.current);
        };
    }, [running, endTime]);

    const start = () => {
        const newEnd = Date.now() + seconds * 1000;
        setEndTime(newEnd);
        localStorage.setItem(STORAGE_KEY, String(newEnd));
        setRunning(true);
    };

    const pause = () => {
        setRunning(false);
        localStorage.removeItem(STORAGE_KEY);
    };

    const reset = () => {
        pause();
        setRemaining(seconds);
        setEndTime(null);
    };

    const finish = () => {
        pause();
        setRemaining(0);
        navigator.vibrate?.([300, 150, 300]);
        onFinish?.();
    };

    const format = (t: number) => {
        const m = Math.floor(t / 60);
        const s = t % 60;
        return `${m}:${s.toString().padStart(2, "0")}`;
    };

    const progress = seconds > 0
        ? 1 - remaining / seconds
        : 0;

    const radius = 28;
    const circumference = 2 * Math.PI * radius;

    return (
        <div className={`rest-timer ${running ? "active" : ""}`}>
            <div className="timer-inner">

                <div className="circle">
                    <svg width="70" height="70">
                        <circle
                            className="bg"
                            strokeWidth="6"
                            r={radius}
                            cx="35"
                            cy="35"
                        />
                        <circle
                            className="progress"
                            strokeWidth="6"
                            r={radius}
                            cx="35"
                            cy="35"
                            strokeDasharray={circumference}
                            strokeDashoffset={circumference * (1 - progress)}
                        />
                    </svg>
                    <div className="time">{format(remaining)}</div>
                </div>

                <div className="actions">
                    {!running ? (
                        <Button onClick={start}>Старт</Button>
                    ) : (
                        <Button onClick={pause}>Пауза</Button>
                    )}
                    <Button variant="ghost" onClick={reset}>
                        Сброс
                    </Button>
                </div>

            </div>
        </div>
    );
}
