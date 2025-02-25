document.addEventListener("DOMContentLoaded", () => {
  const subscribeButton = document.getElementById("subscribe-button");
  const statusDiv = document.getElementById("status");

  function showStatus(message, type) {
    statusDiv.textContent = message;
    statusDiv.className = type;
    statusDiv.style.display = "block";
  }

  async function subscribeToPush() {
    try {
      // Проверяем поддержку
      if (!("serviceWorker" in navigator) || !("PushManager" in window)) {
        throw new Error("Push-уведомления не поддерживаются");
      }

      showStatus("Подключение к сервису уведомлений...", "info");

      // Регистрируем Service Worker
      const registration = await navigator.serviceWorker.register(
        "/service-worker.js"
      );
      console.log("Service Worker зарегистрирован");

      // Проверяем существующую подписку
      let subscription = await registration.pushManager.getSubscription();
      if (subscription) {
        // Отписываемся от старой подписки
        await subscription.unsubscribe();
        console.log("Старая подписка удалена");
      }

      // Запрашиваем разрешение
      const permission = await Notification.requestPermission();
      if (permission !== "granted") {
        throw new Error("Необходимо разрешить уведомления");
      }

      // Получаем новый ключ с сервера
      const response = await fetch("/api/push-key");
      const { publicKey } = await response.json();

      // Создаем новую подписку
      subscription = await registration.pushManager.subscribe({
        userVisibleOnly: true,
        applicationServerKey: urlBase64ToUint8Array(publicKey),
      });

      // Отправляем на сервер
      const subscribeResponse = await fetch("/api/subscribe", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(subscription),
      });

      if (!subscribeResponse.ok) {
        throw new Error("Ошибка сохранения подписки на сервере");
      }

      showStatus("Подписка успешно оформлена!", "success");
      subscribeButton.disabled = true;
    } catch (error) {
      console.error("Ошибка:", error);
      showStatus(`Ошибка: ${error.message}`, "error");
      subscribeButton.disabled = false;
    }
  }

  subscribeButton.addEventListener("click", subscribeToPush);
});

function urlBase64ToUint8Array(base64String) {
  const padding = "=".repeat((4 - (base64String.length % 4)) % 4);
  const base64 = (base64String + padding).replace(/-/g, "+").replace(/_/g, "/");

  const rawData = window.atob(base64);
  const outputArray = new Uint8Array(rawData.length);

  for (let i = 0; i < rawData.length; ++i) {
    outputArray[i] = rawData.charCodeAt(i);
  }
  return outputArray;
}
