import MovieItemComponent from "./MovieItem.js";

export default class CollectionPage extends HTMLElement {
  constructor(endpoint, title) {
    super();
    this.endpoint = endpoint;
    this.title = title;
    this._ulMovies = null;
  }

  async render() {
    const movies = await this.endpoint();

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
      emptyMessage.textContent = "There are no movies";
      this._ulMovies.appendChild(emptyMessage);
    }
  }

  connectedCallback() {
    const template = document.getElementById("template-collection");
    const content = template.content.cloneNode(true);
    this.appendChild(content);

    this._ulMovies = this.querySelector("ul");

    this.render();
  }
}
