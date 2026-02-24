export class YouTubeEmbed extends HTMLElement {
  _iframe = null;

  static get observedAttributes() {
    return ["data-url"];
  }

  attributeChangedCallback(prop, oldValue, newValue) {
    if (prop === "data-url" && newValue) {
      const videoId = newValue.substring(newValue.indexOf("?v") + 3);

      if (!this._iframe) {
        this._iframe = document.createElement("iframe");
        this._iframe.width = "100%";
        this._iframe.height = "100%";
        this._iframe.style.minHeight = "280px";
        this._iframe.style.display = "block";
        this._iframe.style.border = "none";
        this._iframe.title = "YouTube video player";
        this._iframe.frameBorder = "0";
        this._iframe.allow =
          "accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share";
        this._iframe.referrerPolicy = "strict-origin-when-cross-origin";
        this._iframe.allowFullscreen = true;
        this.appendChild(this._iframe);
      }

      this._iframe.src = "https://www.youtube.com/embed/" + videoId;
    }
  }
}

customElements.define("youtube-embed", YouTubeEmbed);
