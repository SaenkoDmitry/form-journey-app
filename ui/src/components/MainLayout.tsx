import React, {useEffect, useState} from 'react';
import {Link, useLocation, useNavigate} from 'react-router-dom';
import {useAuth} from '../context/AuthContext';
import Button from "./Button.tsx";

import {BookOpen, Dumbbell, FolderKanban, Moon, Ruler, Sun, User} from "lucide-react";
import FloatingRestTimer from "./FloatingRestTimer.tsx";
import Toast from "./Toast.tsx";
import {ICONS} from "./IconPicker.tsx";
import {useUserIcon} from "../hooks/useUserIcons.ts";

const tabs = [
    {name: 'Тренировки', path: '/', icon: Dumbbell},
    {name: 'Программы', path: '/programs', icon: FolderKanban},
    {name: 'Замеры', path: '/measurements', icon: Ruler},
    {name: 'Библиотека', path: '/library', icon: BookOpen},
    {name: 'Профиль', path: '/profile', icon: User},
];

const MainLayout: React.FC<{ children: React.ReactNode }> = ({children}) => {
    const location = useLocation();
    const {user, logout, loading} = useAuth();
    const navigate = useNavigate();
    const [toast, setToast] = useState<string | null>(null);

    const {icon} = useUserIcon();
    const CurrentIcon = ICONS[icon];

    const [darkMode, setDarkMode] = useState<boolean>(() => {
        // Читаем из localStorage
        const saved = localStorage.getItem("darkMode");
        if (saved !== null) return saved === "true";

        // Если нет сохранения — используем системную тему
        return window.matchMedia &&
            window.matchMedia("(prefers-color-scheme: dark)").matches;
    });

    // Применяем класс и сохраняем выбор
    useEffect(() => {
        const root = document.documentElement;
        if (darkMode) {
            root.classList.add("dark-theme");
        } else {
            root.classList.remove("dark-theme");
        }
        localStorage.setItem("darkMode", darkMode.toString());
    }, [darkMode]);

    useEffect(() => {
        const handler = () => setToast("Отдых закончен 💪");
        window.addEventListener("rest_timer_finished", handler);
        return () => window.removeEventListener("rest_timer_finished", handler);
    }, []);

    const isMobile = window.innerWidth <= 768;

    const tapFeedback = () => {
        if (navigator.vibrate) navigator.vibrate(10);
    };

    return (
        <div
            style={{
                minHeight: '100dvh',
            }}
        >
            {/* ---------------- DESKTOP NAVBAR ONLY ---------------- */}
            {!isMobile && (
                <nav
                    style={{
                        display: 'flex',
                        justifyContent: 'space-between',
                        alignItems: 'center',
                        padding: '0.6rem 1rem',
                        position: 'sticky',
                        top: 0,
                        zIndex: 10,
                        background: 'var(--color-primary)',
                        color: 'var(--color-text)',
                    }}
                >
                    <div style={{display: 'flex', gap: 12}}>
                        {tabs.slice(0, 4).map((tab) => (
                            <Link
                                key={tab.path}
                                to={tab.path}
                                style={{
                                    padding: '0.5rem 1rem',
                                    textDecoration: 'none',
                                    color: 'white',
                                    fontWeight:
                                        location.pathname === tab.path ? 700 : 400,
                                }}
                            >
                                {tab.name}
                            </Link>
                        ))}
                    </div>

                    <div style={{display: 'flex', alignItems: 'center', gap: 12}}>

                        {/* --- Кнопка переключения темы --- */}
                        <div style={{display: "flex", justifyContent: "flex-end"}}>
                            <Button
                                variant="attention"
                                onClick={() => setDarkMode(!darkMode)}
                                style={{display: "flex", alignItems: "center", gap: 8}}
                            >
                                {darkMode ? <Sun/> : <Moon/>}
                            </Button>
                        </div>

                        {user && (
                            <>
                                <Button variant={"ghost"}
                                        onClick={() => navigate('/profile')}>{user.first_name}
                                    <CurrentIcon/>
                                </Button>
                            </>
                        )}
                    </div>
                </nav>
            )}

            {/* ---------------- CONTENT ---------------- */}
            <div
                style={{
                    padding: '0.5rem',
                    paddingBottom: isMobile ? 110 : 16,
                }}
            >
                {children}
            </div>

            {/* ---------------- FLOATING PILL ---------------- */}
            {isMobile && (
                <div
                    style={{
                        position: 'fixed',
                        left: 12,
                        right: 12,
                        bottom: 'calc(env(safe-area-inset-bottom) + 10px)',
                        height: 70,
                        background: 'var(--color-card)',
                        borderRadius: 24,
                        display: 'flex',
                        padding: 6,
                        gap: 4,
                        boxShadow:
                            '0 8px 30px rgba(0,0,0,0.15), 0 2px 6px rgba(0,0,0,0.06)',
                        zIndex: 20,
                    }}
                >
                    {tabs.map((tab) => {
                        const active = location.pathname === tab.path;
                        const Icon = tab.name == 'Профиль' ? CurrentIcon : tab.icon;

                        return (
                            <Link
                                key={tab.path}
                                to={tab.path}
                                onClick={tapFeedback}
                                style={{
                                    flex: 1,
                                    display: 'flex',
                                    flexDirection: 'column',
                                    alignItems: 'center',
                                    justifyContent: 'center',
                                    borderRadius: 18,
                                    textDecoration: 'none',
                                    transition: 'all .15s ease',
                                    background: active
                                        ? 'var(--color-primary)'
                                        : 'transparent',
                                    color: active ? 'var(--color-btn-text)' : 'var(--color-text)',
                                }}
                            >
                                <Icon size={20}/>
                                <span
                                    style={{
                                        marginTop: 4,
                                        fontSize: 11,
                                        fontWeight: 600,
                                    }}
                                >
                                    {tab.name}
                                </span>
                            </Link>
                        );
                    })}
                </div>
            )}

            <FloatingRestTimer/>
            {toast && <Toast message={toast} onClose={() => setToast(null)}/>}
        </div>
    );
};

export default MainLayout;
