import {useLocation, useNavigate} from "react-router-dom";
import React, {useEffect, useState} from "react";
import {api} from "../api/client.ts";
import Button from "../components/Button.tsx";
import {ArrowLeft, Loader} from "lucide-react";
import Toast from "../components/Toast.tsx";

export default function ExerciseVideoPage() {
    const location = useLocation();
    const navigate = useNavigate();
    const [toast, setToast] = useState<string | null>(null);

    const originalUrl = location.state?.videoUrl;
    const [videoUrl, setVideoUrl] = useState<string | null>(null);
    const [error, setError] = useState<string | null>(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        if (!originalUrl) {
            setError("Видео не найдено");
            setLoading(false);
            return;
        }

        const fetchVideoLink = async () => {
            try {
                const data = await api<{ url: string }>(
                    `/api/video/link?url=${encodeURIComponent(originalUrl)}`
                );
                setVideoUrl(data.url);
            } catch (e: any) {
                setError(e.message || "Ошибка загрузки видео");
                setToast("Ошибка загрузки видео 😢");
            } finally {
                setLoading(false);
            }
        };

        fetchVideoLink();
    }, [originalUrl]);

    if (loading) return <Loader/>;
    if (error) return <div>{error}</div>;

    return (
        <div style={{padding: 16}}>
            <Button variant={"ghost"}
                    onClick={() => navigate(-1)}
                    style={{marginBottom: 10}}
            ><ArrowLeft/> Назад</Button>

            {videoUrl && (
                <video
                    src={videoUrl}
                    controls
                    playsInline
                    style={{width: "100%", borderRadius: 12}}
                />
            )}

            {toast && <Toast message={toast} onClose={() => setToast(null)}/>}
        </div>
    );
}
