import {defineConfig} from 'vite'
import react from '@vitejs/plugin-react'
import {VitePWA} from 'vite-plugin-pwa'

export default defineConfig({
    plugins: [
        react(),
        VitePWA({
            strategies: 'injectManifest',
            srcDir: 'src',
            filename: 'sw.ts',
            registerType: 'autoUpdate',
            devOptions: {
                enabled: true,
                type: 'module',
            },
            workbox: {
                navigateFallback: null
            },
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
    resolve: {
        dedupe: ['react', 'react-dom']
    },
    build: {
        outDir: "../internal/web/dist",
        emptyOutDir: true,
    },

    server: {
        host: true,
        strictPort: false,
        port: 5173,
        allowedHosts: [
            'localhost',
            '127.0.0.1',
            '629a62bc2e1bff.lhr.life',
        ],
        proxy: {
            '/api': 'http://localhost:8080',
        },
    },

    base: "/",
})
