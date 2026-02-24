export default class MovieItemComponent extends HTMLElement {
  constructor(movie) {
    super();
    this._movie = movie;
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

    info.appendChild(title);
    info.appendChild(year);

    article.appendChild(img);
    article.appendChild(info);
    a.appendChild(article);
    this.appendChild(a);
  }
}

customElements.define("movie-item", MovieItemComponent);
