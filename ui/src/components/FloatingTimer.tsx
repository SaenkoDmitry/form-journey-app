import { useEffect, useState, useRef } from "react";
import { useGlobalTimer } from "../context/TimerContext";
import "../styles/FloatingTimer.css";

const POS_STORAGE = "floating_timer_pos";

export default function FloatingTimer() {
    const { endTime, stop } = useGlobalTimer();

    const [remaining, setRemaining] = useState(0);
    const [duration, setDuration] = useState(0);
    const [position, setPosition] = useState({ x: 20, y: 120 });
    const [finished, setFinished] = useState(false);
    const [showNotification, setShowNotification] = useState(false);

    const dragging = useRef(false);
    const offset = useRef({ x: 0, y: 0 });

    // восстановление позиции
    useEffect(() => {
        const saved = localStorage.getItem(POS_STORAGE);
        if (saved) setPosition(JSON.parse(saved));
    }, []);

    // таймер
    useEffect(() => {
        if (!endTime) return;

        setFinished(false);
        setShowNotification(false);

        const total = Math.floor((endTime - Date.now()) / 1000);
        setDuration(total > 0 ? total : 0);

        const interval = setInterval(() => {
            const diff = Math.max(0, Math.floor((endTime - Date.now()) / 1000));
            setRemaining(diff);

            if (diff <= 0) {
                clearInterval(interval);
                setRemaining(0);
                setFinished(true);
                stop();

                // вибрация
                navigator.vibrate?.([400, 200, 400]);

                // визуальное уведомление внутри сайта
                setTimeout(() => setShowNotification(true), 50);

                // Web Notification
                if ("Notification" in window) {
                    if (Notification.permission === "granted") {
                        new Notification("Отдых закончен 💪");
                    } else if (Notification.permission !== "denied") {
                        Notification.requestPermission().then(permission => {
                            if (permission === "granted") {
                                new Notification("Отдых закончен 💪");
                            }
                        });
                    }
                }
            }
        }, 500);

        return () => clearInterval(interval);
    }, [endTime, stop]);

    // Drag handlers
    const handlePointerDown = (clientX: number, clientY: number) => {
        dragging.current = true;
        offset.current = { x: clientX - position.x, y: clientY - position.y };
    };
    const handleMouseDown = (e: React.MouseEvent) => handlePointerDown(e.clientX, e.clientY);
    const handleTouchStart = (e: React.TouchEvent) => handlePointerDown(e.touches[0].clientX, e.touches[0].clientY);

    const handlePointerMove = (clientX: number, clientY: number) => {
        if (!dragging.current) return;

        let newX = clientX - offset.current.x;
        let newY = clientY - offset.current.y;

        const padding = 10;
        const width = window.innerWidth;
        const height = window.innerHeight;
        const size = 70;

        newX = Math.min(Math.max(newX, padding), width - size - padding);
        newY = Math.min(Math.max(newY, padding), height - size - padding);

        setPosition({ x: newX, y: newY });
    };

    const handleMouseMove = (e: MouseEvent) => handlePointerMove(e.clientX, e.clientY);
    const handleTouchMove = (e: TouchEvent) => handlePointerMove(e.touches[0].clientX, e.touches[0].clientY);
    const handlePointerUp = () => {
        if (!dragging.current) return;
        dragging.current = false;
        localStorage.setItem(POS_STORAGE, JSON.stringify(position));
    };

    useEffect(() => {
        window.addEventListener("mousemove", handleMouseMove);
        window.addEventListener("mouseup", handlePointerUp);
        window.addEventListener("touchmove", handleTouchMove);
        window.addEventListener("touchend", handlePointerUp);

        return () => {
            window.removeEventListener("mousemove", handleMouseMove);
            window.removeEventListener("mouseup", handlePointerUp);
            window.removeEventListener("touchmove", handleTouchMove);
            window.removeEventListener("touchend", handlePointerUp);
        };
    }, [position]);

    if (!endTime) return null;

    const radius = 32;
    const circumference = 2 * Math.PI * radius;
    const progress = duration > 0 ? 1 - remaining / duration : 0;

    return (
        <div
            className={`floating-timer ${finished ? "done" : ""}`}
            style={{ left: position.x, top: position.y, width: 70, height: 70 }}
            onMouseDown={handleMouseDown}
            onTouchStart={handleTouchStart}
        >
            <svg width="70" height="70">
                <circle className="bg" r={radius} cx="35" cy="35" />
                <circle
                    className="progress"
                    r={radius}
                    cx="35"
                    cy="35"
                    strokeDasharray={circumference}
                    strokeDashoffset={circumference * (1 - progress)}
                />
            </svg>

            <div className="floating-time">{remaining}</div>

            <button className="floating-close" onClick={stop}>✕</button>

            {finished && showNotification && (
                <div className="floating-notification">Отдых закончен 💪</div>
            )}
        </div>
    );
}
