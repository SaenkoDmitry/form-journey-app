import {defineConfig} from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
    plugins: [react()],
    build: {
        outDir: "../internal/web/dist", // если embed
        emptyOutDir: true,
    },
    server: {
        // разрешаем локальный сервер для туннельных хостов
        host: true,
        strictPort: false,
        port: 5173,
        allowedHosts: [
            'localhost',
            '127.0.0.1',
            'e63d8f032d6afd.lhr.life', // ваш туннельный host, чтобы работал telegram widget
        ],
        proxy: {
            '/api': 'http://localhost:8080',
        },
    },
    base: "/", // важно!
})
