import React, {useEffect, useRef} from 'react';
import {useAuth} from '../context/AuthContext.tsx';

const TelegramLoginWidget: React.FC = () => {
    const {user} = useAuth();
    const widgetRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        if (user || !widgetRef.current) return;

        widgetRef.current.innerHTML = '';

        const botUsername = process.env.NODE_ENV === 'development'
            ? 'fitness_gym_buddy_dev_bot'
            : 'form_journey_bot';

        const callbackUrl = `${window.location.origin}/api/telegram/callback`;

        // Проверяем, iOS PWA или нет
        const isIOS = /iPad|iPhone|iPod/.test(navigator.userAgent);
        const isPWA = window.matchMedia('(display-mode: standalone)').matches;

        const isIOSPWA = isIOS && isPWA;

        if (isIOSPWA) {
            // iOS PWA → используем redirect flow через window.location
            const button = document.createElement('button');
            button.textContent = 'Login via Telegram';
            button.style.fontSize = '16px';
            button.style.padding = '12px 24px';
            button.onclick = () => {
                // Просто редиректим на backend callback
                window.location.href = callbackUrl;
            };
            widgetRef.current.appendChild(button);
        } else {
            // Остальные → fetch flow через data-auth-url
            const script = document.createElement('script');
            script.src = 'https://telegram.org/js/telegram-widget.js?15';
            script.async = true;
            script.setAttribute('data-telegram-login', botUsername);
            script.setAttribute('data-size', 'large');
            script.setAttribute('data-userpic', 'true');
            script.setAttribute('data-auth-url', callbackUrl); // fetch flow
            widgetRef.current.appendChild(script);
        }

        // Redirect flow: убираем data-auth-url
        // Telegram будет редиректить на /api/telegram/callback, который ставит cookie и делает редирект обратно

    }, [user]);


    if (user) return null;
    return <div ref={widgetRef}/>;
};


export default TelegramLoginWidget;
