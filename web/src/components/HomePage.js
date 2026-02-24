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
  }

  async render() {
    const promises = [API.getTopMovies(), API.getRandomMovies()];

    if (Store.loggedIn) {
      promises.push(API.getRecommendations());
    }

    const results = await Promise.all(promises);
    const [topMovies, randomMovies] = results;
    const recommendations = results[2] || null;

    if (!topMovies || !randomMovies) return;

    this._renderMoviesInList(topMovies, this._ulTop10);
    this._renderMoviesInList(randomMovies, this._ulRandom);

    if (recommendations && recommendations.length > 0) {
      this._renderMoviesInList(recommendations, this._ulRecommended);
      this._recommendedSection.style.display = "";
    }
  }

  _renderMoviesInList(movies, ul) {
    while (ul.firstChild) {
      ul.removeChild(ul.firstChild);
    }

    const fragment = document.createDocumentFragment();
    for (let i = 0; i < movies.length; i++) {
      const li = document.createElement("li");
      li.appendChild(new MovieItemComponent(movies[i]));
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

    this.render();
  }
}

customElements.define("home-page", HomePage);
