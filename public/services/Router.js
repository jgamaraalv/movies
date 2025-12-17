import { routes } from "./Routes.js";

const Router = {
  init: () => {
    document.querySelectorAll("a.navlink").forEach((a) => {
      a.addEventListener("click", (event) => {
        event.preventDefault();
        const href = a.getAttribute("href");
        Router.go(href);
      });
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
        break;
      } else if (r.path instanceof RegExp) {
        const match = r.path.exec(route);

        if (match) {
          const params = match.slice(1);
          pageElement = new r.component();
          pageElement.params = params;
          break;
        }
      }
    }

    if (pageElement == null) {
      pageElement = document.createElement("h1");
      pageElement.textContent = "Page not found";
    }

    document.querySelector("main").innerHTML = "";
    document.querySelector("main").appendChild(pageElement);
  },
};

export default Router;
