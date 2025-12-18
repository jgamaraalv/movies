import { routes } from "./Routes.js";

let _mainElement = null;

const Router = {
  init: () => {
    _mainElement = document.querySelector("main");

    // Check if page was server-side rendered
    const ssrDataScript = document.getElementById("ssr-data");
    if (ssrDataScript && _mainElement && _mainElement.children.length > 0) {
      // Page was SSR'd, hydrate instead of re-rendering
      try {
        const ssrData = JSON.parse(ssrDataScript.textContent);
        Router.hydrate(ssrData);
      } catch (e) {
        console.warn("Failed to parse SSR data, falling back to SPA:", e);
        Router.go(location.pathname + location.search, false);
      }
    } else {
      // Normal SPA initialization
      Router.go(location.pathname + location.search, false);
    }

    document.body.addEventListener("click", (event) => {
      const link = event.target.closest("a.navlink");
      if (link) {
        event.preventDefault();
        const url = new URL(link.href);
        Router.go(url.pathname + url.search);
      }
    });

    window.addEventListener("popstate", () => {
      Router.go(location.pathname, false);
    });
  },
  hydrate: (ssrData) => {
    // Hydrate the existing SSR content by attaching event listeners
    // The HTML is already rendered, we just need to make it interactive

    // Find and initialize Web Components that might need hydration
    const movieItems = _mainElement.querySelectorAll("movie-item");
    movieItems.forEach((item) => {
      // Components are already rendered, just ensure they're interactive
      if (item.connectedCallback) {
        // Component might need to re-run connectedCallback for event binding
        // But since HTML is already there, we just need to ensure events are bound
      }
    });

    // Bind YouTube embed if present
    const youtubeEmbed = _mainElement.querySelector("youtube-embed");
    if (youtubeEmbed && youtubeEmbed.dataset.url) {
      // YouTube embed component will handle its own initialization
    }

    // Bind action buttons for movie details page
    const actionsContainer = _mainElement.querySelector("#actions");
    if (actionsContainer) {
      const movieId = actionsContainer.dataset.id;
      actionsContainer.addEventListener("click", (e) => {
        const btn = e.target.closest("button");
        if (!btn) return;

        if (btn.id === "btnFavorites") {
          app.saveToCollection(movieId, "favorite");
        } else if (btn.id === "btnWatchlist") {
          app.saveToCollection(movieId, "watchlist");
        }
      });
    }

    // Remove SSR data script after hydration
    const ssrDataScript = document.getElementById("ssr-data");
    if (ssrDataScript) {
      ssrDataScript.remove();
    }
  },
  go: (route, addToHistory = true) => {
    if (addToHistory) {
      history.pushState(null, "", route);
    }

    const queryIndex = route.indexOf("?");
    const routePath =
      queryIndex !== -1 ? route.substring(0, queryIndex) : route;
    let pageElement = null;

    for (let i = 0; i < routes.length; i++) {
      const r = routes[i];
      if (typeof r.path === "string") {
        if (r.path === routePath) {
          pageElement = new r.component();
          pageElement.loggedIn = r.loggedIn;
          break;
        }
      } else {
        const match = r.path.exec(route);
        if (match) {
          pageElement = new r.component();
          pageElement.loggedIn = r.loggedIn;
          pageElement.params = match.slice(1);
          break;
        }
      }
    }

    if (pageElement) {
      if (pageElement.loggedIn && app.Store.loggedIn === false) {
        app.Router.go("/account/login");
        return;
      }
    } else {
      pageElement = document.createElement("h1");
      pageElement.textContent = "Page not found";
    }

    // Clear any existing SSR content
    if (!document.startViewTransition) {
      while (_mainElement.firstChild) {
        _mainElement.removeChild(_mainElement.firstChild);
      }
      _mainElement.appendChild(pageElement);
    } else {
      const oldPage = _mainElement.firstElementChild;
      if (oldPage) oldPage.style.viewTransitionName = "old";
      pageElement.style.viewTransitionName = "new";
      document.startViewTransition(() => {
        while (_mainElement.firstChild) {
          _mainElement.removeChild(_mainElement.firstChild);
        }
        _mainElement.appendChild(pageElement);
      });
    }
  },
};

export default Router;
