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
    img.alt = this._movie.title + " Poster";

    const p = document.createElement("p");
    p.textContent = this._movie.title + " (" + this._movie.release_year + ")";

    article.appendChild(img);
    article.appendChild(p);
    a.appendChild(article);
    this.appendChild(a);
  }
}

customElements.define("movie-item", MovieItemComponent);
