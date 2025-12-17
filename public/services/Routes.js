import HomePage from "../components/HomePage.js";
import MovieDetailsPage from "../components/MovieDetailsPage.js";
import MoviesPage from "../components/MoviesPage.js";

export const routes = [
  {
    path: "/",
    component: HomePage,
  },

  {
    path: "/movies",
    component: MoviesPage,
  },
  {
    path: /\/movies\/(\d+)/,
    component: MovieDetailsPage,
  },
];
