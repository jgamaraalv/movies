import { API } from "../services/API.js";

const PATH_HEART =
  "M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z";
const PATH_BOOKMARK =
  "M17 3H7c-1.1 0-2 .9-2 2v16l7-3 7 3V5c0-1.1-.9-2-2-2z";
const PATH_CHECK =
  "M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41L9 16.17z";

function createSvgIcon(pathD) {
  const svg = document.createElementNS("http://www.w3.org/2000/svg", "svg");
  svg.setAttribute("viewBox", "0 0 24 24");
  svg.setAttribute("width", "13");
  svg.setAttribute("height", "13");
  svg.setAttribute("fill", "currentColor");
  svg.setAttribute("aria-hidden", "true");
  const path = document.createElementNS("http://www.w3.org/2000/svg", "path");
  path.setAttribute("d", pathD);
  svg.appendChild(path);
  return svg;
}

export default class MovieItemComponent extends HTMLElement {
  constructor(movie, savedIds = null) {
    super();
    this._movie = movie;
    this._savedIds = savedIds;
  }

  connectedCallback() {
    // If content already exists (SSR), skip rendering (hydration mode)
    if (this.children.length > 0) {
      return;
    }

    // If no movie data (shouldn't happen in normal SPA flow), skip
    if (!this._movie) {
      return;
    }

    const a = document.createElement("a");
    a.href = "/movies/" + this._movie.id;
    a.className = "navlink";

    const article = document.createElement("article");

    const img = document.createElement("img");
    img.src = this._movie.poster_url;
    img.alt = this._movie.title + " poster";
    img.loading = "lazy";
    img.width = 185;
    img.height = 278;

    // Score badge
    if (this._movie.score) {
      const scoreBadge = document.createElement("span");
      scoreBadge.className = "movie-score";
      const star = document.createElement("span");
      star.className = "star";
      star.textContent = "\u2605";
      scoreBadge.appendChild(star);
      scoreBadge.appendChild(
        document.createTextNode(" " + Number(this._movie.score).toFixed(1))
      );
      article.appendChild(scoreBadge);
    }

    // Info overlay
    const info = document.createElement("div");
    info.className = "movie-info";

    const title = document.createElement("span");
    title.className = "movie-title";
    title.textContent = this._movie.title;

    const year = document.createElement("span");
    year.className = "movie-year";
    year.textContent = this._movie.release_year;

    // Quick action buttons
    const actions = document.createElement("div");
    actions.className = "movie-card-actions";

    const btnFav = document.createElement("button");
    btnFav.className = "movie-card-btn movie-card-btn--fav";
    btnFav.title = "Add to Favorites";
    btnFav.setAttribute("aria-label", "Add to Favorites");
    btnFav.appendChild(createSvgIcon(PATH_HEART));

    const btnWatch = document.createElement("button");
    btnWatch.className = "movie-card-btn movie-card-btn--watch";
    btnWatch.title = "Add to Watchlist";
    btnWatch.setAttribute("aria-label", "Add to Watchlist");
    btnWatch.appendChild(createSvgIcon(PATH_BOOKMARK));

    // Apply initial active state if movie is already saved
    if (this._savedIds?.favorites?.has(this._movie.id)) {
      btnFav.classList.add("movie-card-btn--fav--active");
    }
    if (this._savedIds?.watchlist?.has(this._movie.id)) {
      btnWatch.classList.add("movie-card-btn--watch--active");
    }

    btnFav.addEventListener("click", (e) => {
      e.preventDefault();
      e.stopPropagation();
      this._handleSave(btnFav, "favorite");
    });

    btnWatch.addEventListener("click", (e) => {
      e.preventDefault();
      e.stopPropagation();
      this._handleSave(btnWatch, "watchlist");
    });

    actions.appendChild(btnFav);
    actions.appendChild(btnWatch);

    info.appendChild(title);
    info.appendChild(year);
    info.appendChild(actions);

    article.appendChild(img);
    article.appendChild(info);
    a.appendChild(article);
    this.appendChild(a);
  }

  async _handleSave(button, collection) {
    if (!app.Store.loggedIn) {
      app.Router.go("/account/");
      return;
    }

    button.disabled = true;
    button.style.opacity = "0.5";
    const originalPath = collection === "favorite" ? PATH_HEART : PATH_BOOKMARK;
    const activeClass =
      collection === "favorite"
        ? "movie-card-btn--fav--active"
        : "movie-card-btn--watch--active";

    try {
      const response = await API.saveToCollection(this._movie.id, collection);
      button.style.opacity = "";
      if (response && response.success) {
        // Dispatch immediately so HomePage starts the refresh timer in parallel with animation
        this.dispatchEvent(
          new CustomEvent("movie:saved", {
            bubbles: true,
            detail: { movieId: this._movie.id, collection },
          })
        );

        button.classList.add("movie-card-btn--done");
        while (button.firstChild) button.removeChild(button.firstChild);
        button.appendChild(createSvgIcon(PATH_CHECK));
        setTimeout(() => {
          button.classList.remove("movie-card-btn--done");
          button.classList.add(activeClass);
          while (button.firstChild) button.removeChild(button.firstChild);
          button.appendChild(createSvgIcon(originalPath));
          button.disabled = false;
        }, 1500);
      } else {
        button.disabled = false;
      }
    } catch {
      button.style.opacity = "";
      button.disabled = false;
    }
  }
}

customElements.define("movie-item", MovieItemComponent);
