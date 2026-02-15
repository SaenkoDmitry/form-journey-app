import {useEffect} from "react";
import {useRestTimer} from "../context/RestTimerContext";
import Button from "./Button";
import "../styles/RestTimer.css";
import {Pause, Play, RotateCcw} from "lucide-react";

type Props = {
    seconds: number;
    autoStartTrigger?: number;
};

export default function RestTimer({
                                      seconds,
                                      autoStartTrigger,
                                  }: Props) {

    const {
        remaining,
        running,
        start,
        pause,
        reset,
        seconds: totalSeconds
    } = useRestTimer();

    // 🔥 автостарт после завершения подхода
    useEffect(() => {
        if (!autoStartTrigger) return;
        start(seconds);
    }, [autoStartTrigger]);

    const format = (t: number) => {
        const m = Math.floor(t / 60);
        const s = t % 60;
        return `${m}:${s.toString().padStart(2, "0")}`;
    };

    const progress =
        totalSeconds > 0
            ? 1 - remaining / totalSeconds
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
                            strokeDashoffset={
                                circumference * (1 - progress)
                            }
                        />
                    </svg>

                    <div className="time">
                        {format(remaining)}
                    </div>
                </div>

                <div className="actions">
                    {!running ? (
                        <Button
                            variant="active"
                            onClick={() => start(seconds)}
                        >
                            <Play size={14}/>Старт
                        </Button>
                    ) : (
                        <Button
                            variant="primary"
                            onClick={pause}
                        >
                            <Pause size={14}/>Пауза
                        </Button>
                    )}

                    <Button
                        variant="ghost"
                        onClick={reset}
                    >
                        <RotateCcw size={14}/>Сброс
                    </Button>
                </div>

            </div>
        </div>
    );
}
