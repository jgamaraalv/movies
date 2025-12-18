# Estratégia de Server-Side Rendering (SSR)

## Visão Geral

Este documento descreve a implementação de **Server-Side Rendering (SSR) híbrido** no projeto Movies, que combina renderização no servidor para SEO com navegação client-side rápida para uma melhor experiência do usuário.

## Problema

O projeto atual é uma **SPA (Single Page Application)** totalmente renderizada pelo JavaScript no cliente. Isso oferece:

- ✅ Navegação rápida e fluida
- ✅ Experiência de usuário excelente
- ❌ **Problema**: Conteúdo vazio no HTML inicial
- ❌ **Problema**: Crawlers de busca não conseguem indexar o conteúdo
- ❌ **Problema**: SEO ruim (sem meta tags dinâmicas, sem conteúdo visível)

## Solução: SSR Híbrido

A solução implementada utiliza **SSR seletivo**:

- **Rotas públicas** (`/`, `/movies`, `/movies/:id`) → Renderizadas no servidor com dados
- **Rotas privadas** (`/account/*`) → Mantêm renderização client-side (SPA)
- **Detecção inteligente**: Crawlers recebem HTML renderizado, navegadores recebem SPA com hidratação

## Arquitetura

### Fluxo de Requisição

```
┌─────────────┐
│   Cliente   │
└──────┬──────┘
       │
       ▼
┌─────────────────────────────────────┐
│  Detecção: Crawler ou Navegador?   │
└──────┬──────────────────┬──────────┘
       │                   │
       ▼                   ▼
┌──────────────┐    ┌──────────────┐
│   Crawler    │    │  Navegador   │
└──────┬───────┘    └──────┬───────┘
       │                   │
       ▼                   ▼
┌──────────────┐    ┌──────────────┐
│  SSR (HTML)  │    │ SPA + Hidratação│
└──────────────┘    └──────────────┘
```

### Componentes da Solução

1. **SSR Handler** (`server/internal/handler/ssr_handler.go`)

   - Detecta crawlers pelo User-Agent
   - Renderiza HTML no servidor com dados do banco
   - Injeta dados JSON para hidratação no cliente

2. **Router Client-Side** (`web/src/services/Router.js`)

   - Detecta conteúdo SSR pré-renderizado
   - Hidrata componentes existentes em vez de re-renderizar
   - Mantém navegação SPA para rotas subsequentes

3. **Rotas SSR** (`server/cmd/api/main.go`)
   - `/` → Home page com top 10 e filmes aleatórios
   - `/movies/:id` → Detalhes do filme
   - `/movies?q=...` → Resultados de busca

## Detalhes de Implementação

### 1. Detecção de Crawlers

O SSR handler detecta crawlers através de:

- **User-Agent**: Lista de crawlers conhecidos (Googlebot, Bingbot, etc.)
- **Query Parameter**: `_escaped_fragment_` (usado por alguns crawlers)
- **Headers**: Ausência de `X-Requested-With` (indica navegação direta)

```go
func (h *SSRHandler) isCrawler(r *http.Request) bool {
    userAgent := strings.ToLower(r.Header.Get("User-Agent"))
    crawlers := []string{"googlebot", "bingbot", ...}
    // Verifica User-Agent
    // Verifica _escaped_fragment_
    // Verifica headers
}
```

### 2. Renderização Server-Side

Para cada rota pública, o handler:

1. Busca dados do banco via Use Cases existentes
2. Renderiza HTML com os dados
3. Injeta meta tags dinâmicas (title, description)
4. Injeta dados JSON para hidratação

```go
func (h *SSRHandler) HomePage(w http.ResponseWriter, r *http.Request) {
    // 1. Buscar dados
    topMoviesOutput, _ := h.movieHandler.getTopMoviesUC.Execute()
    randomMoviesOutput, _ := h.movieHandler.getRandomMoviesUC.Execute()

    // 2. Renderizar HTML
    h.renderPage(w, "home", PageData{
        TopMovies: topMoviesOutput.Movies,
        RandomMovies: randomMoviesOutput.Movies,
    })
}
```

### 3. Hidratação no Cliente

O Router.js detecta conteúdo SSR e hidrata em vez de re-renderizar:

```javascript
const Router = {
  init: () => {
    const ssrDataScript = document.getElementById("ssr-data");
    if (ssrDataScript && _mainElement.children.length > 0) {
      // Hidratar conteúdo existente
      Router.hydrate(JSON.parse(ssrDataScript.textContent));
    } else {
      // Renderização normal (SPA)
      Router.go(location.pathname + location.search);
    }
  },
  hydrate: (ssrData) => {
    // Anexar event listeners aos elementos já renderizados
    // Não precisa re-renderizar o HTML
  },
};
```

## Benefícios

### SEO

- ✅ HTML completo no servidor para crawlers
- ✅ Meta tags dinâmicas (title, description)
- ✅ Conteúdo indexável pelos motores de busca
- ✅ Open Graph tags (futuro)

### Performance

- ✅ First Contentful Paint (FCP) mais rápido
- ✅ Time to Interactive (TTI) mantido rápido
- ✅ Navegação client-side continua rápida após primeira carga

### Experiência do Usuário

- ✅ Conteúdo visível imediatamente (sem loading)
- ✅ Navegação SPA mantida para rotas subsequentes
- ✅ Funcionalidade interativa preservada

## Rotas com SSR

| Rota            | SSR    | Motivo                                      |
| --------------- | ------ | ------------------------------------------- |
| `/`             | ✅ Sim | Página inicial pública, importante para SEO |
| `/movies/:id`   | ✅ Sim | Detalhes de filme, compartilhamento social  |
| `/movies?q=...` | ✅ Sim | Resultados de busca, indexação              |
| `/account/*`    | ❌ Não | Páginas privadas, não precisam de SEO       |

## Rotas sem SSR (SPA)

- `/account/login` - Página privada
- `/account/register` - Página privada
- `/account/` - Dashboard do usuário
- `/account/favorites` - Coleção privada
- `/account/watchlist` - Coleção privada

## Estratégia de Cache

### Para Crawlers

- Sem cache (sempre renderiza no servidor)
- Dados sempre atualizados do banco

### Para Navegadores

- Service Worker continua funcionando
- Cache de assets estáticos
- Navegação SPA usa cache quando possível

## Meta Tags Dinâmicas

Cada página SSR inclui meta tags otimizadas:

```html
<!-- Home Page -->
<title>Movies - Discover Top Films</title>
<meta
  name="description"
  content="Discover the top movies and find something great to watch today"
/>

<!-- Movie Details -->
<title>The Matrix - Welcome to the Real World</title>
<meta
  name="description"
  content="A computer hacker learns about the true nature of reality..."
/>

<!-- Search Results -->
<title>'action' movies</title>
```

## Estrutura de Dados SSR

Os dados SSR são injetados como JSON no HTML:

```html
<script id="ssr-data" type="application/json">
  {
    "pageType": "home",
    "data": {
      "topMovies": [...],
      "randomMovies": [...]
    }
  }
</script>
```

O cliente usa esses dados para:

- Hidratação de componentes
- Evitar requisições desnecessárias na primeira carga
- Manter sincronização entre servidor e cliente

## Compatibilidade

### Crawlers Suportados

- ✅ Googlebot
- ✅ Bingbot
- ✅ Facebook Crawler
- ✅ Twitter Bot
- ✅ LinkedIn Bot
- ✅ E outros principais crawlers

### Navegadores

- ✅ Todos os navegadores modernos
- ✅ Fallback para SPA se SSR falhar
- ✅ Funciona sem JavaScript (crawlers)

## Desenvolvimento

### Testando SSR Localmente

1. **Simular crawler**:

```bash
curl -A "Googlebot" http://localhost:8080/
```

2. **Verificar HTML renderizado**:

```bash
curl -A "Googlebot" http://localhost:8080/movies/1 | grep -A 10 "<main>"
```

3. **Testar hidratação**:
   - Abrir página no navegador
   - Verificar console do navegador
   - Verificar que componentes estão hidratados

### Debugging

- **Logs SSR**: Verificar logs do servidor Go para erros de renderização
- **Console do Cliente**: Verificar se hidratação está funcionando
- **Network Tab**: Verificar que dados SSR estão sendo injetados

## Limitações e Considerações

### Performance do Servidor

- SSR adiciona carga no servidor (renderização por requisição)
- Considerar cache de páginas SSR no futuro se necessário
- Monitorar tempo de resposta

### Complexidade

- Código de renderização duplicado (servidor + cliente)
- Manter sincronização entre SSR e componentes client-side
- Testes mais complexos (SSR + SPA)

### Futuras Melhorias

1. **Cache de SSR**

   - Cache de páginas renderizadas (Redis/Memcached)
   - Invalidação quando dados mudam

2. **Incremental Static Regeneration (ISR)**

   - Pré-renderizar páginas estáticas
   - Atualizar em background

3. **Streaming SSR**

   - Enviar HTML em chunks
   - Melhorar Time to First Byte (TTFB)

4. **Open Graph Tags**
   - Meta tags para compartilhamento social
   - Imagens de preview dinâmicas

## Conclusão

A implementação de SSR híbrido oferece:

- ✅ Melhor SEO sem sacrificar performance
- ✅ Conteúdo indexável pelos motores de busca
- ✅ Experiência de usuário mantida
- ✅ Arquitetura flexível e extensível

A solução é **progressiva**: funciona mesmo se SSR falhar (fallback para SPA), garantindo que a aplicação sempre funcione.
