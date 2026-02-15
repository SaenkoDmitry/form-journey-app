import { useEffect, useState, useRef } from "react";
import { useRestTimer } from "../context/RestTimerContext";
import { useLocation, useNavigate } from "react-router-dom";
import "../styles/FloatingRestTimer.css";

export default function FloatingRestTimer() {
    const { remaining, seconds, running } = useRestTimer();
    const location = useLocation();
    const navigate = useNavigate();

    const [position, setPosition] = useState({ x: 20, y: 100 });
    const [blinking, setBlinking] = useState(false);
    const touchRef = useRef<{ startX: number; startY: number } | null>(null);

    useEffect(() => {
        const saved = localStorage.getItem("floatingTimerPosition");
        if (saved) setPosition(JSON.parse(saved));
    }, []);

    useEffect(() => {
        localStorage.setItem("floatingTimerPosition", JSON.stringify(position));
    }, [position]);

    const shouldRender = running && !location.pathname.startsWith("/sessions/");

    // Пульс на последние 5 секунд
    useEffect(() => {
        if (!shouldRender) return;
        setBlinking(remaining > 0 && remaining <= 5);
    }, [remaining, shouldRender]);

    // Вибрация по завершению
    useEffect(() => {
        if (!shouldRender) return;
        if (remaining === 0 && running) {
            navigator.vibrate?.([300, 150, 300]);
        }
    }, [remaining, running, shouldRender]);

    if (!shouldRender) return null;

    const progress = seconds > 0 ? 1 - remaining / seconds : 0;
    const radius = 34; // увеличенный радиус
    const circumference = 2 * Math.PI * radius;

    // touch для перемещения
    const onTouchStart = (e: React.TouchEvent) => {
        const touch = e.touches[0];
        touchRef.current = { startX: touch.clientX - position.x, startY: touch.clientY - position.y };
    };

    const onTouchMove = (e: React.TouchEvent) => {
        if (!touchRef.current) return;
        const touch = e.touches[0];
        setPosition({ x: touch.clientX - touchRef.current.startX, y: touch.clientY - touchRef.current.startY });
    };

    const onTouchEnd = () => { touchRef.current = null; };

    const handleClick = () => {
        // Переход на текущую тренировку
        let link = localStorage.getItem("floatingTimerLink");
        if (link != "") {
            navigate(link);
        }
    };

    const minutes = Math.floor(remaining / 60);
    const secs = (remaining % 60).toString().padStart(2, "0");

    return (
        <div
            className={`floating-rest-timer ${blinking ? "blinking" : ""}`}
            style={{ top: position.y, left: position.x }}
            onTouchStart={onTouchStart}
            onTouchMove={onTouchMove}
            onTouchEnd={onTouchEnd}
            onClick={handleClick}
        >
            <svg>
                <circle r={radius} cx="40" cy="40" />
                <circle
                    className="progress"
                    r={radius}
                    cx="40"
                    cy="40"
                    strokeDasharray={circumference}
                    strokeDashoffset={circumference * (1 - progress)}
                />
                <text x="40" y="46" textAnchor="middle" className="timer-text">
                    {remaining > 0 ? `${minutes}:${secs}` : ""}
                </text>
            </svg>
        </div>
    );
}
