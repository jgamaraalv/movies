import API from "../services/API.js";
import Store from "../services/Store.js";
import MovieItemComponent from "./MovieItem.js";

export default class HomePage extends HTMLElement {
  constructor() {
    super();
    this._ulTop10 = null;
    this._ulRandom = null;
    this._recommendedSection = null;
    this._ulRecommended = null;
    this._savedIds = { favorites: new Set(), watchlist: new Set() };
    this._recsRefreshTimer = null;
  }

  async render() {
    const promises = [API.getTopMovies(), API.getRandomMovies()];

    if (Store.loggedIn) {
      promises.push(API.getRecommendations(), API.getFavorites(), API.getWatchlist());
    }

    const results = await Promise.all(promises);
    const [topMovies, randomMovies] = results;

    if (!topMovies || !randomMovies) return;

    this._savedIds = { favorites: new Set(), watchlist: new Set() };

    let recommendations = null;

    if (Store.loggedIn) {
      recommendations = results[2] || null;
      const favs = results[3];
      const watch = results[4];
      if (Array.isArray(favs)) favs.forEach((m) => this._savedIds.favorites.add(m.id));
      if (Array.isArray(watch)) watch.forEach((m) => this._savedIds.watchlist.add(m.id));
    }

    this._renderMoviesInList(topMovies, this._ulTop10, this._savedIds);
    this._renderMoviesInList(randomMovies, this._ulRandom, this._savedIds);

    if (recommendations && recommendations.length > 0) {
      this._renderMoviesInList(recommendations, this._ulRecommended, this._savedIds);
      this._recommendedSection.style.display = "";
    }
  }

  async _refreshRecommendations() {
    // Retry up to 3 times in case the goroutine is still computing (table briefly empty)
    for (let attempt = 0; attempt < 3; attempt++) {
      if (attempt > 0) await new Promise((r) => setTimeout(r, 2000));
      const recommendations = await API.getRecommendations();
      if (recommendations && recommendations.length > 0) {
        this._renderMoviesInList(recommendations, this._ulRecommended, this._savedIds);
        this._recommendedSection.style.display = "";
        return;
      }
    }
  }

  _renderMoviesInList(movies, ul, savedIds = null) {
    while (ul.firstChild) {
      ul.removeChild(ul.firstChild);
    }

    const fragment = document.createDocumentFragment();
    for (let i = 0; i < movies.length; i++) {
      const li = document.createElement("li");
      li.appendChild(new MovieItemComponent(movies[i], savedIds));
      fragment.appendChild(li);
    }
    ul.appendChild(fragment);
  }

  connectedCallback() {
    const template = document.getElementById("template-home");
    const content = template.content.cloneNode(true);
    this.appendChild(content);

    this._ulTop10 = this.querySelector("#top-10 ul");
    this._ulRandom = this.querySelector("#random ul");
    this._recommendedSection = this.querySelector("#recommended");
    this._ulRecommended = this.querySelector("#recommended ul");

    this.addEventListener("movie:saved", (e) => {
      const { movieId, collection } = e.detail;
      if (collection === "favorite") this._savedIds.favorites.add(movieId);
      else this._savedIds.watchlist.add(movieId);
      clearTimeout(this._recsRefreshTimer);
      this._recsRefreshTimer = setTimeout(
        () => this._refreshRecommendations(),
        3000
      );
    });

    this.render();
  }
}

customElements.define("home-page", HomePage);
