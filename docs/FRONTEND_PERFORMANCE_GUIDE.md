# Guia de Performance Frontend

Este documento descreve os requisitos não funcionais de performance, boas práticas e padrões a serem seguidos no desenvolvimento de componentes frontend da aplicação.

---

## Requisitos Não Funcionais (RNF)

### RNF-01: Minimização de Reflows/Repaints

A UI deve realizar o menor número possível de operações de reflow/repaint. Múltiplas manipulações do DOM devem ser agrupadas em uma única operação usando `DocumentFragment`.

### RNF-02: Eliminação de Parsing de HTML

A aplicação não deve utilizar `innerHTML` ou `insertAdjacentHTML` para criação de elementos. Todos os elementos devem ser criados via `document.createElement()` e suas propriedades definidas diretamente.

### RNF-03: Cache de Referências do DOM

Elementos do DOM que são acessados múltiplas vezes devem ter suas referências armazenadas em cache (variáveis de instância) para evitar buscas repetidas.

### RNF-04: Otimização de Event Listeners

A aplicação deve utilizar event delegation quando apropriado, minimizando o número de event listeners registrados e evitando a criação de múltiplas closures.

### RNF-05: Paralelização de Operações Assíncronas

Chamadas de API independentes devem ser executadas em paralelo usando `Promise.all()` ao invés de sequencialmente com múltiplos `await`.

### RNF-06: Reutilização de Elementos

Quando possível, elementos existentes devem ser reutilizados ao invés de destruídos e recriados, especialmente para elementos custosos como `<iframe>`.

### RNF-07: Manipulação em DocumentFragment

Modificações em elementos devem ser feitas enquanto estão em um `DocumentFragment` (fora do DOM), antes de serem inseridos na árvore do documento.

---

## Boas Práticas para Desenvolvimento de Componentes

### 1. Cache de Elementos no `connectedCallback`

Sempre armazene referências de elementos que serão acessados posteriormente em propriedades da classe.

```javascript
connectedCallback() {
  const template = document.getElementById("template-example");
  const content = template.content.cloneNode(true);
  this.appendChild(content);

  // ✅ Cache das referências
  this._ulList = this.querySelector("ul");
  this._btnSubmit = this.querySelector("#submit");
}
```

### 2. Use DocumentFragment para Inserções em Lote

Ao inserir múltiplos elementos, sempre use `DocumentFragment` para agrupar as operações.

```javascript
const fragment = document.createDocumentFragment();
for (let i = 0; i < items.length; i++) {
  const li = document.createElement("li");
  li.textContent = items[i].name;
  fragment.appendChild(li);
}
this._ulList.appendChild(fragment); // Única operação no DOM
```

### 3. Prefira `removeChild` ao invés de `innerHTML = ""`

Para limpar o conteúdo de um elemento, use um loop com `removeChild`.

```javascript
while (element.firstChild) {
  element.removeChild(element.firstChild);
}
```

### 4. Use `textContent` ao invés de `innerHTML` para Texto

Quando o conteúdo é apenas texto, use `textContent` para evitar parsing de HTML.

```javascript
element.textContent = "Texto simples";
```

### 5. Prefira `for` Clássico ao invés de `forEach`

Para iterações de alta performance, use o loop `for` tradicional.

```javascript
for (let i = 0; i < array.length; i++) {
  // operações
}
```

### 6. Modifique Elementos ANTES de Inserir no DOM

Faça todas as modificações em elementos enquanto estão fora do DOM (no template clonado ou fragment).

```javascript
const content = template.content.cloneNode(true);
content.querySelector("h2").textContent = title; // Modificação fora do DOM
content.querySelector("img").src = imageUrl; // Modificação fora do DOM
this.appendChild(content); // Única inserção no DOM
```

### 7. Use Event Delegation para Múltiplos Handlers

Ao invés de adicionar listeners em cada elemento filho, adicione um único listener no container.

```javascript
container.addEventListener("click", (e) => {
  const btn = e.target.closest("button");
  if (!btn) return;

  if (btn.id === "btnSave") handleSave();
  else if (btn.id === "btnDelete") handleDelete();
});
```

### 8. Paralelize Chamadas de API Independentes

Use `Promise.all` para executar múltiplas requisições simultaneamente.

```javascript
const [users, products] = await Promise.all([
  API.getUsers(),
  API.getProducts(),
]);
```

---

## Patterns vs Anti-Patterns

### Inserção de Elementos em Lista

#### ❌ Anti-Pattern: Inserção Individual no Loop

```javascript
movies.forEach((movie) => {
  const li = document.createElement("li");
  li.textContent = movie.title;
  ulMovies.appendChild(li); // ⚠️ Reflow a cada iteração!
});
```

#### ✅ Pattern: DocumentFragment

```javascript
const fragment = document.createDocumentFragment();
for (let i = 0; i < movies.length; i++) {
  const li = document.createElement("li");
  li.textContent = movies[i].title;
  fragment.appendChild(li);
}
ulMovies.appendChild(fragment); // ✅ Único reflow
```

---

### Limpeza de Conteúdo

#### ❌ Anti-Pattern: innerHTML para Limpar

```javascript
ulMovies.innerHTML = ""; // ⚠️ Parsing de HTML desnecessário
```

#### ✅ Pattern: removeChild Loop

```javascript
while (ulMovies.firstChild) {
  ulMovies.removeChild(ulMovies.firstChild);
}
```

---

### Criação de Elementos com HTML

#### ❌ Anti-Pattern: innerHTML com Template String

```javascript
li.innerHTML = `
  <img src="${actor.image}" alt="${actor.name}">
  <p>${actor.name}</p>
`; // ⚠️ Parsing de HTML + possível XSS
```

#### ✅ Pattern: DOM API

```javascript
const img = document.createElement("img");
img.src = actor.image;
img.alt = actor.name;

const p = document.createElement("p");
p.textContent = actor.name;

li.appendChild(img);
li.appendChild(p);
```

---

### Acesso a Elementos do DOM

#### ❌ Anti-Pattern: Busca Repetida

```javascript
async render() {
  this.querySelector("h2").textContent = title;     // ⚠️ Busca 1
  this.querySelector("ul").innerHTML = "";          // ⚠️ Busca 2
  // ... mais tarde no código
  this.querySelector("ul").appendChild(fragment);   // ⚠️ Busca 3 (mesmo elemento!)
}
```

#### ✅ Pattern: Cache de Referências

```javascript
connectedCallback() {
  // Cache no momento da conexão
  this._heading = this.querySelector("h2");
  this._list = this.querySelector("ul");
}

async render() {
  this._heading.textContent = title;  // ✅ Acesso direto
  // limpar lista...
  this._list.appendChild(fragment);   // ✅ Acesso direto
}
```

---

### Event Listeners

#### ❌ Anti-Pattern: Múltiplos Listeners com Closures

```javascript
this.querySelector("#btnFavorites").addEventListener("click", () => {
  saveToFavorites(movieId);
});
this.querySelector("#btnWatchlist").addEventListener("click", () => {
  saveToWatchlist(movieId);
});
// ⚠️ 2 listeners + 2 closures + 2 buscas no DOM
```

#### ✅ Pattern: Event Delegation

```javascript
const movieId = this._movie.id;
actionsContainer.addEventListener("click", (e) => {
  const btn = e.target.closest("button");
  if (!btn) return;

  if (btn.id === "btnFavorites") saveToFavorites(movieId);
  else if (btn.id === "btnWatchlist") saveToWatchlist(movieId);
});
// ✅ 1 listener + 1 closure
```

---

### Chamadas de API

#### ❌ Anti-Pattern: Await Sequencial

```javascript
const topMovies = await API.getTopMovies(); // ⏱️ 200ms
const randomMovies = await API.getRandomMovies(); // ⏱️ 200ms
// Total: ~400ms
```

#### ✅ Pattern: Promise.all Paralelo

```javascript
const [topMovies, randomMovies] = await Promise.all([
  API.getTopMovies(),
  API.getRandomMovies(),
]);
// Total: ~200ms (tempo da mais lenta)
```

---

### Reutilização de Elementos Custosos

#### ❌ Anti-Pattern: Recriar Elemento a Cada Update

```javascript
attributeChangedCallback(prop, oldValue, newValue) {
  if (prop === "data-url") {
    this.innerHTML = `<iframe src="${embedUrl}"></iframe>`;
    // ⚠️ Destrói e recria iframe a cada mudança
  }
}
```

#### ✅ Pattern: Reutilizar e Atualizar

```javascript
attributeChangedCallback(prop, oldValue, newValue) {
  if (prop === "data-url" && newValue) {
    if (!this._iframe) {
      this._iframe = document.createElement("iframe");
      // configurações do iframe...
      this.appendChild(this._iframe);
    }
    this._iframe.src = embedUrl; // ✅ Apenas atualiza o src
  }
}
```

---

### Iteração sobre Arrays

#### ❌ Anti-Pattern: forEach com Closure

```javascript
movies.forEach((movie) => {
  // ⚠️ Cria closure a cada iteração
  const li = document.createElement("li");
  li.textContent = movie.title;
  fragment.appendChild(li);
});
```

#### ✅ Pattern: for Loop Clássico

```javascript
for (let i = 0; i < movies.length; i++) {
  const li = document.createElement("li");
  li.textContent = movies[i].title;
  fragment.appendChild(li);
}
```

---

### Modificação de Template Clonado

#### ❌ Anti-Pattern: Modificar Após Inserção

```javascript
const content = template.content.cloneNode(true);
this.appendChild(content); // Insere no DOM

// ⚠️ Cada modificação abaixo causa reflow
this.querySelector("h2").textContent = title;
this.querySelector("img").src = imageUrl;
this.querySelector("p").textContent = description;
```

#### ✅ Pattern: Modificar Antes da Inserção

```javascript
const content = template.content.cloneNode(true);

// ✅ Modificações no DocumentFragment (fora do DOM)
content.querySelector("h2").textContent = title;
content.querySelector("img").src = imageUrl;
content.querySelector("p").textContent = description;

this.appendChild(content); // Único reflow
```

---

## Resumo de Impacto

| Técnica                    | Impacto                            |
| -------------------------- | ---------------------------------- |
| DocumentFragment           | Reduz N reflows para 1             |
| Cache de referências       | Elimina buscas repetidas no DOM    |
| removeChild vs innerHTML   | Elimina parsing de HTML            |
| DOM API vs innerHTML       | Elimina parsing + previne XSS      |
| Event Delegation           | Reduz listeners e closures         |
| Promise.all                | Reduz tempo de carregamento        |
| for vs forEach             | Elimina overhead de closures       |
| Reutilização de elementos  | Evita destruição/recriação custosa |
| Modificar antes de inserir | Agrupa reflows em uma operação     |

---

## Referências

- [MDN - DocumentFragment](https://developer.mozilla.org/en-US/docs/Web/API/DocumentFragment)
- [MDN - Reflow](https://developer.mozilla.org/en-US/docs/Glossary/Reflow)
- [Google - Avoid Large, Complex Layouts](https://web.dev/avoid-large-complex-layouts-and-layout-thrashing/)
- [Google - DOM Size](https://web.dev/dom-size/)
