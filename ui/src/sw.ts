/// <reference lib="webworker" />
import { precacheAndRoute } from 'workbox-precaching'
import { registerRoute } from 'workbox-routing'
import { NetworkOnly } from 'workbox-strategies'

declare let self: ServiceWorkerGlobalScope

precacheAndRoute(self.__WB_MANIFEST)

// SW не должен перехватывать Telegram callback
registerRoute(
    ({ url }) => url.pathname.startsWith('/api/') && !url.pathname.includes('/telegram/callback'),
    new NetworkOnly()
)

// PUSH
self.addEventListener('push', (event) => {
    console.log('Push received in SW:', event.data?.text());
    const data = event.data?.json() || {};
    console.log('Push data', data);

    const options: NotificationOptions = {
        body: data?.body || 'Таймер завершён!',
        icon: '/web-app-manifest-192x192.png',
        badge: '/web-app-manifest-192x192.png',
        requireInteraction: true,
        tag: data?.tag || '',
        data: {
            url: data?.url || '/'
        }
    }
    console.log('Push data options', options);
    setTimeout(() => { console.log("Waited 2 seconds!"); }, 2000);

    event.waitUntil(
        self.registration.showNotification(
            data.title || 'Form Journey',
            options
        )
    )
})

// CLICK
self.addEventListener('notificationclick', (event) => {
    event.notification.close()

    const url = event.notification.data?.url || '/'

    event.waitUntil(
        self.clients.matchAll({ type: 'window', includeUncontrolled: true })
            .then((clientsArr) => {
                for (const client of clientsArr) {
                    if (client.url.includes(url) && 'focus' in client)
                        return client.focus()
                }
                return self.clients.openWindow(url)
            })
    )
})
