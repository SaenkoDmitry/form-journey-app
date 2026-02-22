import React, { useEffect, useState } from "react";
import { api } from "../api/client.ts";
import Button from "../components/Button.tsx";
import { ArrowDown, ArrowUp, Loader } from "lucide-react";
import Toast from "../components/Toast.tsx";

interface ExerciseVideoInlineProps {
    originalUrl: string;
}

export default function ExerciseVideoInline({ originalUrl }: ExerciseVideoInlineProps) {
    const [videoUrl, setVideoUrl] = useState<string | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [toast, setToast] = useState<string | null>(null);
    const [open, setOpen] = useState(false);

    useEffect(() => {
        if (!originalUrl) {
            setError("Видео не найдено");
            setLoading(false);
            return;
        }

        const CACHE_KEY = `video-${originalUrl}`;
        const cached = localStorage.getItem(CACHE_KEY);
        if (cached) {
            const { url, expires } = JSON.parse(cached);
            if (Date.now() < expires) {
                setVideoUrl(url);
                setLoading(false);
                return;
            }
        }

        const fetchVideoLink = async () => {
            try {
                const data = await api<{ url: string }>(
                    `/api/video/link?url=${encodeURIComponent(originalUrl)}`
                );
                setVideoUrl(data.url);

                // локальный кэш на 4 минуты
                localStorage.setItem(CACHE_KEY, JSON.stringify({
                    url: data.url,
                    expires: Date.now() + 4 * 60 * 1000
                }));
            } catch (e: any) {
                setError(e.message || "Ошибка загрузки видео");
                setToast("Ошибка загрузки видео 😢");
            } finally {
                setLoading(false);
            }
        };

        fetchVideoLink();
    }, [originalUrl]);

    return (
        <div style={{ marginTop: 16 }}>
            <Button
                variant="ghost"
                onClick={() => setOpen(!open)}
                style={{ marginBottom: 8 }}
            >
                {open ? <ArrowUp /> : <ArrowDown />}
                {open ? "Свернуть видео" : "Техника упражнения"}
            </Button>

            {open && (
                <div style={{ padding: 8, borderRadius: 8, border: "1px solid #eee" }}>
                    {loading && <Loader />}
                    {error && <div>{error}</div>}
                    {videoUrl && !loading && (
                        <video
                            src={videoUrl}
                            controls
                            playsInline
                            style={{ width: "100%", borderRadius: 12 }}
                        />
                    )}
                </div>
            )}

            {toast && <Toast message={toast} onClose={() => setToast(null)} />}
        </div>
    );
}
