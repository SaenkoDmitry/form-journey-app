import React from "react";
import {useAuth} from "../context/AuthContext";
import TelegramLoginWidget from "./TelegramLoginWidget";
import Button from "../components/Button";

const ProfilePage: React.FC = () => {
    const {user, logout, loading} = useAuth();

    return (
        <div
            style={{
                maxWidth: 420,
                margin: '0 auto',
                padding: '1rem',
                display: 'flex',
                flexDirection: 'column',
                gap: 20,
            }}
        >
            {/* ---------------- NOT LOGGED IN ---------------- */}
            {!loading && !user && (
                <div
                    style={{
                        background: '#fff',
                        borderRadius: 20,
                        padding: '2rem 1.5rem',
                        boxShadow: '0 6px 20px rgba(0,0,0,0.06)',
                        textAlign: 'center',
                    }}
                >
                    <div style={{fontSize: 42, marginBottom: 12}}>🔐</div>

                    <div
                        style={{
                            fontSize: 16,
                            fontWeight: 600,
                            marginBottom: 16,
                        }}
                    >
                        Войдите в аккаунт
                    </div>

                    <TelegramLoginWidget/>
                </div>
            )}

            {/* ---------------- LOGGED IN ---------------- */}
            {user && (
                <>
                    <div
                        style={{
                            background: '#fff',
                            borderRadius: 20,
                            padding: '1.5rem',
                            boxShadow: '0 6px 20px rgba(0,0,0,0.06)',
                            textAlign: 'center',
                        }}
                    >
                        <div style={{fontSize: 42, marginBottom: 8}}>👤</div>

                        <div style={{fontSize: 18, fontWeight: 600}}>
                            {user.first_name}
                        </div>

                        {user.username && (
                            <div style={{opacity: 0.6, fontSize: 14}}>
                                @{user.username}
                            </div>
                        )}
                    </div>

                    <Button
                        variant="danger"
                        onClick={logout}
                        style={{
                            width: '100%',
                            height: 48,
                            fontSize: 16,
                            borderRadius: 14,
                        }}
                    >
                        Выйти из аккаунта
                    </Button>
                </>
            )}
        </div>
    );
};

export default ProfilePage;
