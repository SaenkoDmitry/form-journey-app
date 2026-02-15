import React, { useEffect } from "react";
import "../styles/Toast.css";

interface ToastProps {
    message: string;
    onClose?: () => void;
}

export default function Toast({ message, onClose }: ToastProps) {
    useEffect(() => {
        const timer = setTimeout(() => {
            onClose?.();
        }, 3000); // показываем 3 секунды
        return () => clearTimeout(timer);
    }, [onClose]);

    return (
        <div className="toast">
            {message}
        </div>
    );
}
