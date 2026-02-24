import { API } from "../services/API.js";

export default class MovieDetailsPage extends HTMLElement {
  _movieId = null;
  _movie = null;

  async render() {
    try {
      this._movie = await API.getMovieById(this._movieId);
    } catch {
      return;
    }
    if (!this._movie) return;

    const template = document.getElementById("template-movie-details");
    const content = template.content.cloneNode(true);

    content.querySelector("h2").textContent = this._movie.title;
    content.querySelector("h3").textContent = this._movie.tagline;
    const posterImg = content.querySelector("img");
    posterImg.src = this._movie.poster_url;
    posterImg.alt = this._movie.title + " poster";
    content.querySelector("#trailer").dataset.url = this._movie.trailer_url;
    content.querySelector("#overview").textContent = this._movie.overview;

    this._renderMetadata(content.querySelector("#metadata"));
    this._renderGenres(content.querySelector("#genres"));
    this._renderCast(content.querySelector("#cast"));
    this._bindActions(content);

    this.appendChild(content);
  }

  _renderMetadata(dl) {
    const fragment = document.createDocumentFragment();
    const metadata = [
      ["Release Year", this._movie.release_year],
      ["Score", `${this._movie.score} / 10`],
      ["Popularity", this._movie.popularity],
    ];

    for (let i = 0; i < metadata.length; i++) {
      const dt = document.createElement("dt");
      dt.textContent = metadata[i][0];
      const dd = document.createElement("dd");
      dd.textContent = metadata[i][1];
      fragment.appendChild(dt);
      fragment.appendChild(dd);
    }

    dl.appendChild(fragment);
  }

  _renderGenres(ul) {
    const genres = this._movie.genres;
    const fragment = document.createDocumentFragment();

    for (let i = 0; i < genres.length; i++) {
      const li = document.createElement("li");
      li.textContent = genres[i].name;
      fragment.appendChild(li);
    }

    ul.appendChild(fragment);
  }

  _renderCast(ul) {
    const casting = this._movie.casting;
    const fragment = document.createDocumentFragment();

    for (let i = 0; i < casting.length; i++) {
      const actor = casting[i];
      const li = document.createElement("li");

      const img = document.createElement("img");
      img.src = actor.image_url ?? "/images/generic_actor.jpg";
      img.alt = `${actor.first_name} ${actor.last_name}`;
      img.loading = "lazy";
      img.width = 56;
      img.height = 80;

      const p = document.createElement("p");
      p.textContent = `${actor.first_name} ${actor.last_name}`;

      li.appendChild(img);
      li.appendChild(p);
      fragment.appendChild(li);
    }

    ul.appendChild(fragment);
  }

  _bindActions(content) {
    const movieId = this._movie.id;
    const actionsContainer = content.querySelector("#actions");

    actionsContainer.addEventListener("click", (e) => {
      const btn = e.target.closest("button");
      if (!btn) return;

      if (btn.id === "btnFavorites") {
        app.saveToCollection(movieId, "favorite");
      } else if (btn.id === "btnWatchlist") {
        app.saveToCollection(movieId, "watchlist");
      }
    });
  }

  connectedCallback() {
    this._movieId = this.params[0];
    this.render();
  }
}

customElements.define("movie-details-page", MovieDetailsPage);
