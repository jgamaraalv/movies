const CACHE_NAME = "movies-cache-v3";

// App shell assets to pre-cache (stable paths only — hashed assets
// are cached at runtime via stale-while-revalidate on first load)
const PRECACHE_URLS = [
  "/offline.html",
  "/index.html",
  "/images/logo.svg",
  "/images/icon.png",
];

self.addEventListener("install", (event) => {
  event.waitUntil(
    caches
      .open(CACHE_NAME)
      .then((cache) => cache.addAll(PRECACHE_URLS))
      .then(() => self.skipWaiting())
  );
});

// Activate event - clean up old caches
self.addEventListener("activate", (event) => {
  event.waitUntil(
    caches
      .keys()
      .then((cacheNames) => {
        return Promise.all(
          cacheNames
            .filter((name) => name !== CACHE_NAME)
            .map((name) => caches.delete(name))
        );
      })
      .then(() => {
        // Take control of clients immediately
        return self.clients.claim();
      })
  );
});

// Fetch event - handle caching strategies
self.addEventListener("fetch", (event) => {
  const requestUrl = new URL(event.request.url);

  // Only handle http/https requests (ignore chrome-extension://, etc.)
  if (requestUrl.protocol !== "http:" && requestUrl.protocol !== "https:") {
    return;
  }

  // Handle /api/ GET requests (network first, cache fallback)
  // POST/PUT/DELETE are not cacheable — let them pass through
  if (requestUrl.pathname.startsWith("/api/")) {
    if (event.request.method !== "GET") return;

    event.respondWith(
      fetch(event.request)
        .then((networkResponse) => {
          // Cache successful network response
          return caches.open(CACHE_NAME).then((cache) => {
            cache.put(event.request, networkResponse.clone());
            return networkResponse;
          });
        })
        .catch(() => {
          // If network fails, try cache
          return caches.match(event.request).then((cachedResponse) => {
            return (
              cachedResponse ||
              new Response(JSON.stringify({ error: "offline" }), {
                status: 503,
                headers: { "Content-Type": "application/json" },
              })
            );
          });
        })
    );
  }
  // Handle navigation requests (HTML pages) - network first, offline fallback
  else if (event.request.mode === "navigate") {
    event.respondWith(
      fetch(event.request)
        .then((networkResponse) => {
          return caches.open(CACHE_NAME).then((cache) => {
            cache.put(event.request, networkResponse.clone());
            return networkResponse;
          });
        })
        .catch(() => {
          return caches.match(event.request).then((cachedResponse) => {
            return cachedResponse || caches.match("/offline.html");
          });
        })
    );
  }
  // Handle static assets (JS, CSS, images) - stale-while-revalidate
  else {
    event.respondWith(
      caches.open(CACHE_NAME).then((cache) => {
        return cache.match(event.request).then((cachedResponse) => {
          const fetchPromise = fetch(event.request)
            .then((networkResponse) => {
              cache.put(event.request, networkResponse.clone());
              return networkResponse;
            })
            .catch(() => {
              // If no cached version and network fails, return a proper error response
              if (!cachedResponse) {
                return new Response("Network error", { status: 503, statusText: "Service Unavailable" });
              }
            });

          return cachedResponse || fetchPromise;
        });
      })
    );
  }
});
