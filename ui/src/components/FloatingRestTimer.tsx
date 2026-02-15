import { useEffect, useState } from "react";
import { useRestTimer } from "../context/RestTimerContext";
import { useLocation } from "react-router-dom";

export default function FloatingRestTimer() {
    const { remaining, seconds, running } = useRestTimer();
    const location = useLocation();

    // ✅ хуки всегда в начале
    const [position, setPosition] = useState({ x: 20, y: 100 });
    const [blink, setBlink] = useState(false);

    // 🔹 скрываем визуально, если на странице тренировки или таймер неактивен
    const shouldRender = running && !location.pathname.startsWith("/sessions/");

    // мигание последние 5 секунд
    useEffect(() => {
        if (!shouldRender) return;

        if (remaining <= 5 && remaining > 0) {
            const interval = setInterval(() => setBlink(prev => !prev), 500);
            return () => clearInterval(interval);
        } else {
            setBlink(false);
        }
    }, [remaining, shouldRender]);

    // вибрация по завершению
    useEffect(() => {
        if (!shouldRender) return;
        if (remaining === 0 && running) {
            navigator.vibrate?.([300, 150, 300]);
        }
    }, [remaining, running, shouldRender]);

    if (!shouldRender) return null; // ✅ рендер только после вызова всех хуков

    const progress = seconds > 0 ? 1 - remaining / seconds : 0;
    const radius = 26;
    const circumference = 2 * Math.PI * radius;

    const format = (t: number) => {
        const m = Math.floor(t / 60);
        const s = t % 60;
        return `${m}:${s.toString().padStart(2, "0")}`;
    };

    return (
        <div
            style={{
                position: "fixed",
                top: position.y,
                left: position.x,
                zIndex: 9999,
                touchAction: "none",
                cursor: "grab",
                width: "64px",
                height: "64px",
                background: "#fff",
                borderRadius: "50%",
                boxShadow: "0 8px 24px rgba(0,0,0,0.18)",
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
                userSelect: "none",
                opacity: blink ? 0.4 : 1,
                transition: "opacity 0.3s"
            }}
            onPointerDown={(e) => {
                const startX = e.clientX - position.x;
                const startY = e.clientY - position.y;

                const move = (ev: PointerEvent) => {
                    setPosition({
                        x: ev.clientX - startX,
                        y: ev.clientY - startY
                    });
                };

                const up = () => {
                    window.removeEventListener("pointermove", move);
                    window.removeEventListener("pointerup", up);
                };

                window.addEventListener("pointermove", move);
                window.addEventListener("pointerup", up);
            }}
        >
            <svg width="64" height="64">
                {/* фон круга убран, чтобы только прогресс-кольцо */}
                <circle
                    r={radius}
                    cx="32"
                    cy="32"
                    fill="none"
                />
                <circle
                    r={radius}
                    cx="32"
                    cy="32"
                    fill="none"
                    stroke="#4f46e5"
                    strokeWidth={6}
                    strokeDasharray={circumference}
                    strokeDashoffset={circumference * (1 - progress)}
                    strokeLinecap="round"
                    style={{ transition: "stroke-dashoffset 0.3s linear" }}
                />
                <text
                    x="32"
                    y="36"
                    textAnchor="middle"
                    fontSize="14"
                    fontWeight="600"
                    fill="#111"
                >
                    {remaining > 0 ? format(remaining) : ""}
                </text>
            </svg>
        </div>
    );
}
