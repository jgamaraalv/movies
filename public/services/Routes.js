import HomePage from "../components/HomePage.js";
import MovieDetailsPage from "../components/MovieDetailsPage.js";

export const routes = [
  {
    path: "/",
    component: HomePage,
  },
  {
    path: /\/movies\/(\d+)/,
    component: MovieDetailsPage,
  },
];
