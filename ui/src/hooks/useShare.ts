import {useState} from "react";

export const useShare = () => {
    const [isOpen, setIsOpen] = useState(false);
    const [url, setUrl] = useState<string | null>(null);

    const openShare = async (shareUrl: string) => {
        setUrl(shareUrl);

        if (navigator.share) {
            try {
                await navigator.share({
                    title: 'Моя тренировка',
                    url: shareUrl,
                });
                return;
            } catch (e) {
                // пользователь мог отменить — не считаем ошибкой
                return;
            }
        }

        // fallback UI
        setIsOpen(true);
    };

    const close = () => setIsOpen(false);

    return {isOpen, url, openShare, close};
};
