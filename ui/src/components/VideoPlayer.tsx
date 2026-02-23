import {useEffect, useState} from "react";
import {Loader} from "lucide-react";
import {useVideo} from "../hooks/useVideo.ts";

interface VideoPlayerProps {
    url: string;
}

export default function VideoPlayer({url}: VideoPlayerProps) {
    const {videoUrl, loading, error, fetchVideo} = useVideo(url);
    const [open, setOpen] = useState(false);

    useEffect(() => {
        fetchVideo();
    }, []);

    return (
        <div style={{marginTop: 8, borderRadius: 8}}>
            {loading && <Loader/>}
            {error && <div>{error}</div>}
            {videoUrl && !loading && (
                <video
                    src={videoUrl}
                    controls
                    playsInline
                    style={{width: "100%", borderRadius: 12}}
                />
            )}
        </div>
    );
}
