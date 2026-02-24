import API from "../services/API.js";
import MovieItemComponent from "./MovieItem.js";

export default class MoviesPage extends HTMLElement {
  _ulMovies = null;
  _selectOrder = null;
  _selectFilter = null;

  async render(query) {
    const urlParams = new URLSearchParams(window.location.search);
    const order = urlParams.get("order") ?? "";
    const genre = urlParams.get("genre") ?? "";

    const movies = await API.searchMovies(query, order, genre);

    while (this._ulMovies.firstChild) {
      this._ulMovies.removeChild(this._ulMovies.firstChild);
    }

    if (movies && movies.length > 0) {
      const fragment = document.createDocumentFragment();
      for (let i = 0; i < movies.length; i++) {
        const li = document.createElement("li");
        li.appendChild(new MovieItemComponent(movies[i]));
        fragment.appendChild(li);
      }
      this._ulMovies.appendChild(fragment);
    } else {
      const emptyMessage = document.createElement("h3");
      emptyMessage.textContent = "There are no movies with your search";
      this._ulMovies.appendChild(emptyMessage);
    }

    if (order) this._selectOrder.value = order;
    if (genre) this._selectFilter.value = genre;
  }

  async loadGenres() {
    const genres = await API.getGenres();
    if (!genres) return;

    while (this._selectFilter.firstChild) {
      this._selectFilter.removeChild(this._selectFilter.firstChild);
    }

    const fragment = document.createDocumentFragment();

    const defaultOption = document.createElement("option");
    defaultOption.value = "";
    defaultOption.textContent = "Filter by Genre";
    fragment.appendChild(defaultOption);

    for (let i = 0; i < genres.length; i++) {
      const option = document.createElement("option");
      option.value = genres[i].id;
      option.textContent = genres[i].name;
      fragment.appendChild(option);
    }

    this._selectFilter.appendChild(fragment);
  }

  connectedCallback() {
    const template = document.getElementById("template-movies");
    const content = template.content.cloneNode(true);
    this.appendChild(content);

    this._ulMovies = this.querySelector("ul");
    this._selectOrder = this.querySelector("#order");
    this._selectFilter = this.querySelector("#filter");

    const urlParams = new URLSearchParams(window.location.search);
    const query = urlParams.get("q");
    if (query) {
      this.querySelector("h2").textContent = "\u201C" + query + "\u201D movies";
      this.render(query);
    } else {
      app.showError();
    }

    this.loadGenres();
  }
}

customElements.define("movies-page", MoviesPage);
