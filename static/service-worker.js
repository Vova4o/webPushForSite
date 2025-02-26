self.addEventListener("install", (event) => {
  console.log("[Service Worker] Установка");
  self.skipWaiting();
});

self.addEventListener("activate", (event) => {
  console.log("[Service Worker] Активация");
  event.waitUntil(self.clients.claim());
});

self.addEventListener("push", (event) => {
  console.log("[Service Worker] Push событие получено", event);

  let data;
  try {
    const text = event.data.text();
    console.log("[Service Worker] Получены данные:", text);
    data = JSON.parse(text);
  } catch (e) {
    console.error("[Service Worker] Ошибка разбора данных:", e);
    data = {
      title: "Новое уведомление",
      body: "Получено уведомление без данных",
    };
  }

  const title = data.title || "Уведомление";
  const options = {
    body: data.body || "Пришло новое уведомление",
    icon: "/icon.png",
    badge: "/badge.png",
    data: {
      url: data.url || "/",
    },
    requireInteraction: true,
    vibrate: [100, 50, 100],
  };

  console.log("[Service Worker] Показываем уведомление:", { title, options });

  event.waitUntil(
    self.registration
      .showNotification(title, options)
      .then(() => console.log("[Service Worker] Уведомление показано"))
      .catch((error) =>
        console.error("[Service Worker] Ошибка показа уведомления:", error)
      )
  );
});

self.addEventListener("notificationclick", (event) => {
  console.log("[Service Worker] Клик по уведомлению", event);
  event.notification.close();

  const url = (event.notification.data && event.notification.data.url) || "/";

  event.waitUntil(
    clients.openWindow(url).then((windowClient) => {
      if (windowClient) {
        return windowClient.focus();
      }
    })
  );
});

// Добавляем обработчик для отладки
self.addEventListener("message", (event) => {
  console.log("[SW] Получено сообщение:", event.data);
});
