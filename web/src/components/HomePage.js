import API from "../services/API.js";
import MovieItemComponent from "./MovieItem.js";

export default class HomePage extends HTMLElement {
  constructor() {
    super();
    this._ulTop10 = null;
    this._ulRandom = null;
  }

  async render() {
    const [topMovies, randomMovies] = await Promise.all([
      API.getTopMovies(),
      API.getRandomMovies(),
    ]);

    this._renderMoviesInList(topMovies, this._ulTop10);
    this._renderMoviesInList(randomMovies, this._ulRandom);
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

    this.render();
  }
}

customElements.define("home-page", HomePage);
