import { routes } from "./Routes.js";

let _mainElement = null;

const Router = {
  init: () => {
    _mainElement = document.querySelector("main");

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

    Router.go(location.pathname + location.search);
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
