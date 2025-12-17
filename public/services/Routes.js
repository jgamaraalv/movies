import HomePage from "../components/HomePage.js";
import MovieDetailsPage from "../components/MovieDetailsPage.js";
import MoviesPage from "../components/MoviesPage.js";
import RegisterPage from "../components/RegisterPage.js";
import LoginPage from "../components/LoginPage.js";

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
  {
    path: "/account/register",
    component: RegisterPage,
  },
  {
    path: "/account/login",
    component: LoginPage,
  },
];
