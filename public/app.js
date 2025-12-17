import HomePage from "./components/HomePage.js";
import "./components/AnimatedLoading.js";
import "./components/YouTubeEmbed.js";
import MovieDetailsPage from "./components/MovieDetailsPage.js";
import Router from "./services/Router.js";

window.app = {
  Router,
  search: (event) => {
    event.preventDefault();
    const keywords = document.querySelector("input[type=search]").value;
  },
};

window.addEventListener("DOMContentLoaded", () => {
  app.Router.init();
});
