import React from 'react';
import {Route, Routes} from 'react-router-dom';
import {AuthProvider} from './context/AuthContext';
import MainLayout from './components/MainLayout';

import Home from './pages/Home';
import WorkoutPage from './pages/WorkoutPage';
import StatsPage from './pages/StatsPage';
import ProgramsPage from './pages/ProgramsPage';
import ProgramDetailsPage from './pages/ProgramDetailsPage';
import DayDetailsPage from './pages/DayDetailsPage';
import MeasurementsPage from './pages/MeasurementsPage';
import LibraryPage from './pages/LibraryPage';
import ProfilePage from './pages/ProfilePage';

import RequireAuth from './components/RequireAuth';

const App = () => (
    <AuthProvider>
        <Routes>
            {/* Публичная страница профиля */}
            <Route path="/profile" element={<MainLayout><ProfilePage/></MainLayout>}/>

            {/* Защищённые страницы */}
            <Route
                path="/"
                element={
                    <RequireAuth>
                        <MainLayout><Home/></MainLayout>
                    </RequireAuth>
                }
            />
            <Route
                path="/workout/:id"
                element={
                    <RequireAuth>
                        <MainLayout><WorkoutPage/></MainLayout>
                    </RequireAuth>
                }
            />
            <Route
                path="/stats"
                element={
                    <RequireAuth>
                        <MainLayout><StatsPage/></MainLayout>
                    </RequireAuth>
                }
            />
            <Route
                path="/programs"
                element={
                    <RequireAuth>
                        <MainLayout><ProgramsPage/></MainLayout>
                    </RequireAuth>
                }
            />
            <Route
                path="/programs/:id"
                element={
                    <RequireAuth>
                        <MainLayout><ProgramDetailsPage/></MainLayout>
                    </RequireAuth>
                }
            />
            <Route
                path="/programs/:programId/days/:dayId"
                element={
                    <RequireAuth>
                        <MainLayout><DayDetailsPage/></MainLayout>
                    </RequireAuth>
                }
            />
            <Route
                path="/measurements"
                element={
                    <RequireAuth>
                        <MainLayout><MeasurementsPage/></MainLayout>
                    </RequireAuth>
                }
            />
            <Route
                path="/library"
                element={
                    <RequireAuth>
                        <MainLayout><LibraryPage/></MainLayout>
                    </RequireAuth>
                }
            />
        </Routes>
    </AuthProvider>
);

export default App;
