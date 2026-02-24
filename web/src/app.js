import "./components/AnimatedLoading.js";
import "./components/YouTubeEmbed.js";
import Router from "./services/Router.js";
import API from "./services/API.js";
import Store from "./services/Store.js";

window.app = {
  API,
  Store,
  Router,
  showError: (
    message = "There was an error loading the page",
    goToHome = true
  ) => {
    document.querySelector("#alert-modal").showModal();
    document.querySelector("#alert-modal p").textContent = message;
    if (goToHome) app.Router.go("/");
    return;
  },
  closeError: () => {
    document.getElementById("alert-modal").close();
  },
  showOffline: () => {
    const main = document.querySelector("main");
    while (main.firstChild) main.removeChild(main.firstChild);

    const section = document.createElement("section");
    section.style.cssText =
      "display:flex;flex-direction:column;align-items:center;justify-content:center;min-height:60vh;text-align:center;padding:2rem";

    const icon = document.createElement("div");
    icon.style.cssText = "font-size:4rem;margin-bottom:1.5rem;opacity:0.6";
    icon.textContent = "\u{1F39E}";

    const h2 = document.createElement("h2");
    h2.style.cssText =
      "font-family:var(--font-display),sans-serif;font-size:2rem;letter-spacing:0.04em;margin-bottom:0.75rem";
    h2.textContent = "You're Offline";

    const p = document.createElement("p");
    p.style.cssText =
      "color:var(--text-secondary);font-size:1.1rem;max-width:360px;line-height:1.5;margin-bottom:2rem";
    p.textContent =
      "It looks like you've lost your internet connection. Check your network and try again.";

    const btn = document.createElement("button");
    btn.className = "action-btn";
    btn.textContent = "Retry";
    btn.addEventListener("click", () => location.reload());

    section.append(icon, h2, p, btn);
    main.appendChild(section);
  },
  search: (event) => {
    event.preventDefault();
    const keywords = document.querySelector("input[type=search]").value;
    if (keywords.length > 1) {
      app.Router.go(`/movies?q=${keywords}`);
    }
  },
  searchOrderChange: (order) => {
    const urlParams = new URLSearchParams(window.location.search);
    const q = urlParams.get("q");
    const genre = urlParams.get("genre") ?? "";
    app.Router.go(`/movies?q=${q}&order=${order}&genre=${genre}`);
  },
  searchFilterChange: (genre) => {
    const urlParams = new URLSearchParams(window.location.search);
    const q = urlParams.get("q");
    const order = urlParams.get("order") ?? "";
    app.Router.go(`/movies?q=${q}&order=${order}&genre=${genre}`);
  },
  register: async (event) => {
    event.preventDefault();
    const form = event.target;
    const btn = form.querySelector("button[type=submit]");
    const errorEl = document.getElementById("register-error");
    let errors = [];
    const name = document.getElementById("register-name").value;
    const email = document.getElementById("register-email").value;
    const password = document.getElementById("register-password").value;
    const passwordConfirm = document.getElementById(
      "register-password-confirm"
    ).value;

    errorEl.textContent = "";

    if (name.length < 4) errors.push("Enter your complete name");
    if (email.length < 8) errors.push("Enter a valid email address");
    if (password.length < 6) errors.push("Password must be at least 6 characters");
    if (password != passwordConfirm) errors.push("Passwords don't match");
    if (errors.length == 0) {
      btn.disabled = true;
      btn.classList.add("btn-loading");
      try {
        const response = await API.register(name, email, password);
        if (response.success) {
          app.Store.jwt = response.jwt;
          app.Router.go("/account/");
        } else {
          errorEl.textContent = response.message;
        }
      } finally {
        btn.disabled = false;
        btn.classList.remove("btn-loading");
      }
    } else {
      errorEl.textContent = errors.join(". ");
    }
  },
  login: async (event) => {
    event.preventDefault();
    const form = event.target;
    const btn = form.querySelector("button[type=submit]");
    const errorEl = document.getElementById("login-error");
    let errors = [];
    const email = document.getElementById("login-email").value;
    const password = document.getElementById("login-password").value;

    errorEl.textContent = "";

    if (email.length < 8) errors.push("Enter a valid email address");
    if (password.length < 6) errors.push("Password must be at least 6 characters");
    if (errors.length == 0) {
      btn.disabled = true;
      btn.classList.add("btn-loading");
      try {
        const response = await API.authenticate(email, password);
        if (response.success) {
          app.Store.jwt = response.jwt;
          app.Router.go("/account/");
        } else {
          errorEl.textContent = response.message;
        }
      } finally {
        btn.disabled = false;
        btn.classList.remove("btn-loading");
      }
    } else {
      errorEl.textContent = errors.join(". ");
    }
  },
  logout: () => {
    Store.jwt = null;
    app.Router.go("/");
  },
  saveToCollection: async (movie_id, collection) => {
    if (app.Store.loggedIn) {
      try {
        const response = await API.saveToCollection(movie_id, collection);
        if (response.success) {
          switch (collection) {
            case "favorite":
              app.Router.go("/account/favorites");
              break;
            case "watchlist":
              app.Router.go("/account/watchlist");
          }
        } else {
          app.showError("We couldn't save the movie.");
        }
      } catch (e) {
        console.log(e);
      }
    } else {
      app.Router.go("/account/");
    }
  },
};

window.addEventListener("DOMContentLoaded", () => {
  app.Router.init();
  if ("serviceWorker" in navigator) {
    navigator.serviceWorker.register("/sw.js").catch((err) => {
      console.error("SW registration failed:", err);
    });
  }
});
