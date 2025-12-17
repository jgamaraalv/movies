import HomePage from "./components/HomePage.js";
import "./components/AnimatedLoading.js";
import "./components/YouTubeEmbed.js";
import MovieDetailsPage from "./components/MovieDetailsPage.js";

window.app = { 
    search: (event) => {
        event.preventDefault();
        const keywords = document.querySelector("input[type=search]").value;
        
    },    
}

window.addEventListener("DOMContentLoaded", () => {
  document.querySelector("main").appendChild(new HomePage());
  document.querySelector("main").appendChild(new MovieDetailsPage());
})