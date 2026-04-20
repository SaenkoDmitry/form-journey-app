import React from "react";

type Props = {
    isOpen: boolean;
    url: string | null;
    onClose: () => void;
};

const ShareSheet: React.FC<Props> = ({isOpen, url, onClose}) => {
    if (!isOpen || !url) return null;

    const copy = async () => {
        try {
            await navigator.clipboard.writeText(url);
            alert("Ссылка скопирована");
        } catch {
            alert("Не удалось скопировать");
        }
        onClose();
    };

    return (
        <div style={{
            position: "fixed",
            inset: 0,
            background: "rgba(0,0,0,0.4)",
            display: "flex",
            alignItems: "flex-end",
            zIndex: 1000,
        }} onClick={onClose}>
            <div style={{
                background: "#fff",
                width: "100%",
                borderTopLeftRadius: 20,
                borderTopRightRadius: 20,
                padding: 16,
                animation: "slideUp 0.25s ease",
            }} onClick={e => e.stopPropagation()}>
                <div style={{
                    width: 40,
                    height: 5,
                    background: "#ccc",
                    borderRadius: 10,
                    margin: "0 auto 12px",
                }}/>

                <h3 style={{marginBottom: 16}}>Поделиться</h3>

                <div style={{
                    display: "grid",
                    gridTemplateColumns: "repeat(4, 1fr)",
                    gap: 12,
                    marginBottom: 16,
                }}>
                    <button onClick={copy}>📋 Копировать</button>

                    <a href={`https://t.me/share/url?url=${url}`} target="_blank">
                        📩 Telegram
                    </a>

                    <a href={`https://wa.me/?text=${url}`} target="_blank">
                        💬 WhatsApp
                    </a>

                    <a href={`mailto:?body=${url}`}>
                        ✉️ Email
                    </a>
                </div>

                <button style={{
                    width: "100%",
                    padding: 12,
                    borderRadius: 12,
                    border: "none",
                    background: "#eee",
                }} onClick={onClose}>
                    Отмена
                </button>
            </div>
        </div>
    );
};

export default ShareSheet;
