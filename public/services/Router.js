import { routes } from "./Routes.js";

const Router = {
  init: () => {
    window.addEventListener("click", (event) => {
      const link = event.target.closest("a.navlink");
      if (link) {
        event.preventDefault();
        const href = link.getAttribute("href");
        Router.go(href);
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

    const routePath = route.includes("?") ? route.split("?")[0] : route;
    let pageElement = null;

    for (const r of routes) {
      if (typeof r.path === "string" && r.path === routePath) {
        pageElement = new r.component();
        pageElement.loggedIn = r.loggedIn;
      } else if (r.path instanceof RegExp) {
        const match = r.path.exec(route);
        if (match) {
          const params = match.slice(1);
          pageElement = new r.component();
          pageElement.loggedIn = r.loggedIn;

          pageElement.params = params;
        }
      }
      if (pageElement) {
        if (pageElement.loggedIn && app.Store.loggedIn == false) {
          app.Router.go("/account/login");
          return;
        }
        break;
      }
    }

    if (pageElement == null) {
      pageElement = document.createElement("h1");
      pageElement.textContent = "Page not found";
    }

    function updatePage() {
      document.querySelector("main").innerHTML = "";
      document.querySelector("main").appendChild(pageElement);
    }

    if (!document.startViewTransition) {
      updatePage();
    } else {
      const oldPage = document.querySelector("main").firstElementChild;
      if (oldPage) oldPage.style.viewTransitionName = "old";
      pageElement.style.viewTransitionName = "new";
      document.startViewTransition(() => updatePage());
    }
  },
};

export default Router;
