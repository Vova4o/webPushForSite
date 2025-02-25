self.addEventListener("install", (event) => {
  console.log("[SW] Установка Service Worker");
  self.skipWaiting();
});

self.addEventListener("activate", (event) => {
  console.log("[SW] Service Worker активирован");
  event.waitUntil(clients.claim());
});

self.addEventListener("push", (event) => {
  console.log("[SW] Push событие получено", event);

  let payload;
  try {
    payload = event.data.json();
    console.log("[SW] Данные уведомления:", payload);
  } catch (err) {
    payload = {
      title: "Новое уведомление",
      body: event.data ? event.data.text() : "Нет данных",
    };
    console.error("[SW] Ошибка разбора данных:", err);
  }

  const notificationPromise = self.registration.showNotification(
    payload.title,
    {
      body: payload.body,
      icon: "/icon.png",
      data: payload,
      requireInteraction: true,
    }
  );

  event.waitUntil(notificationPromise);
});

self.addEventListener("notificationclick", (event) => {
  console.log("[SW] Клик по уведомлению", event);
  event.notification.close();

  const url = event.notification.data.url || "/";
  event.waitUntil(clients.openWindow(url));
});

// Добавляем обработчик для отладки
self.addEventListener("message", (event) => {
  console.log("[SW] Получено сообщение:", event.data);
});
