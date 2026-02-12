import React, { useEffect } from 'react';
import { useAuth } from '../context/AuthContext';
import { useNavigate } from 'react-router-dom';

const RequireAuth: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const { user, loading } = useAuth();
    const navigate = useNavigate();

    useEffect(() => {
        if (!loading && !user) {
            // Пользователь не залогинен — кидаем на страницу логина
            navigate('/profile');
        }
    }, [user, loading]);

    if (loading || !user) return <p>Загрузка...</p>; // пока проверяем авторизацию
    return <>{children}</>;
};

export default RequireAuth;
