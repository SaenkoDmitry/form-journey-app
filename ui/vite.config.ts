import {defineConfig} from 'vite'
import react from '@vitejs/plugin-react'
import { VitePWA } from 'vite-plugin-pwa'

// https://vite.dev/config/
export default defineConfig({
    plugins: [
        react(),
        VitePWA({
            registerType: 'autoUpdate',
            manifest: {
                name: 'Form Journey · Training 🏔️',
                short_name: 'Form Journey',
                start_url: '/',
                display: 'standalone',
                background_color: '#ffffff',
                theme_color: '#ffffff',
                icons: [
                    {
                        src: '/web-app-manifest-192x192.png',
                        sizes: '192x192',
                        type: 'image/png'
                    },
                    {
                        src: '/web-app-manifest-512x512.png',
                        sizes: '512x512',
                        type: 'image/png'
                    }
                ]
            }
        })
    ],
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
            '498e41742ef017.lhr.life', // ваш туннельный host, чтобы работал telegram widget
        ],
        proxy: {
            '/api': 'http://localhost:8080',
        },
    },
    base: "/", // важно!
})
