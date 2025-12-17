import HomePage from "./components/HomePage.js";

window.app = { 
    search: (event) => {
        event.preventDefault();
        const keywords = document.querySelector("input[type=search]").value;
        
    },    
}

window.addEventListener("DOMContentLoaded", () => {
  document.querySelector("main").appendChild(new HomePage());
})