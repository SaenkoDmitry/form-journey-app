import { createContext, useContext, useEffect, useState } from "react";

const STORAGE_KEY = "global_rest_end";

interface TimerContextType {
    endTime: number | null;
    start: (seconds: number) => void;
    stop: () => void;
}

const TimerContext = createContext<TimerContextType | null>(null);

export const TimerProvider = ({ children }: { children: React.ReactNode }) => {
    const [endTime, setEndTime] = useState<number | null>(null);

    // восстановление после перезагрузки
    useEffect(() => {
        const saved = localStorage.getItem(STORAGE_KEY);
        if (saved) {
            const parsed = Number(saved);
            if (parsed > Date.now()) {
                setEndTime(parsed);
            } else {
                localStorage.removeItem(STORAGE_KEY);
            }
        }
    }, []);

    const start = (seconds: number) => {
        const newEnd = Date.now() + seconds * 1000;
        localStorage.setItem(STORAGE_KEY, String(newEnd));
        setEndTime(newEnd);
    };

    const stop = () => {
        localStorage.removeItem(STORAGE_KEY);
        setEndTime(null);
    };

    return (
        <TimerContext.Provider value={{ endTime, start, stop }}>
            {children}
        </TimerContext.Provider>
    );
};

export const useGlobalTimer = () => {
    const ctx = useContext(TimerContext);
    if (!ctx) throw new Error("TimerProvider not found");
    return ctx;
};
