import React, {StrictMode} from 'react'
import {createRoot} from 'react-dom/client'
import './styles/index.css'
import App from './App.tsx'

import {BrowserRouter} from 'react-router-dom';
import {TimerProvider} from "./context/TimerContext.tsx";

createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <BrowserRouter>
            <TimerProvider>
                <App/>
            </TimerProvider>
        </BrowserRouter>
    </StrictMode>
);