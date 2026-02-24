class AnimatedLoading extends HTMLElement {
  constructor() {
    super();
  }

  connectedCallback() {
    let qty = this.dataset.elements ?? 1;
    let width = this.dataset.width ?? "100px";
    let height = this.dataset.height ?? "10px";
    for (let i = 0; i < qty; i++) {
      const wrapper = document.createElement("div");
      wrapper.setAttribute("class", "loading-wave");
      wrapper.style.width = width;
      wrapper.style.height = height;
      wrapper.style.margin = "0 5px";
      wrapper.style.display = "inline-block";
      wrapper.style.borderRadius = "8px";
      wrapper.style.flexShrink = "0";
      this.appendChild(wrapper);
    }
  }
}

customElements.define("animated-loading", AnimatedLoading);
