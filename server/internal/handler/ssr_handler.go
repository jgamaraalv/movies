package handler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jgamaraalv/movies.git/internal/domain/repository"
	"github.com/jgamaraalv/movies.git/internal/usecase/movie"
	"github.com/jgamaraalv/movies.git/models"
	"github.com/jgamaraalv/movies.git/pkg/logger"
)

type SSRHandler struct {
	movieHandler *MovieHandler
	logger       *logger.Logger
	publicDir    string
}

func NewSSRHandler(movieHandler *MovieHandler, log *logger.Logger) (*SSRHandler, error) {
	publicDir := os.Getenv("PUBLIC_DIR")
	if publicDir == "" {
		publicDir = "public"
		if _, err := os.Stat(publicDir); os.IsNotExist(err) {
			publicDir = "../../public"
		}
	}
	publicDir, _ = filepath.Abs(publicDir)

	return &SSRHandler{
		movieHandler: movieHandler,
		logger:       log,
		publicDir:    publicDir,
	}, nil
}

func (h *SSRHandler) isCrawler(r *http.Request) bool {
	userAgent := strings.ToLower(r.Header.Get("User-Agent"))
	crawlers := []string{
		"googlebot",
		"bingbot",
		"slurp",
		"duckduckbot",
		"baiduspider",
		"yandexbot",
		"sogou",
		"exabot",
		"facebot",
		"ia_archiver",
		"facebookexternalhit",
		"twitterbot",
		"rogerbot",
		"linkedinbot",
		"embedly",
		"quora link preview",
		"showyoubot",
		"outbrain",
		"pinterest",
		"developers.google.com/+/web/snippet",
		"slackbot",
		"vkShare",
		"W3C_Validator",
		"whatsapp",
		"flipboard",
		"tumblr",
		"bitlybot",
		"skypeuripreview",
		"nuzzel",
		"redditbot",
		"applebot",
		"flipboard",
		"tumblr",
		"bitlybot",
		"skypeuripreview",
		"nuzzel",
		"redditbot",
		"applebot",
		"crawler",
		"spider",
		"bot",
	}

	for _, crawler := range crawlers {
		if strings.Contains(userAgent, crawler) {
			return true
		}
	}

	if r.URL.Query().Get("_escaped_fragment_") != "" {
		return true
	}

	return false
}

// shouldUseSSR determines if we should use SSR for this request
func (h *SSRHandler) shouldUseSSR(r *http.Request) bool {
	// Use SSR for crawlers
	if h.isCrawler(r) {
		return true
	}

	// Use SSR for direct navigation (not AJAX requests)
	// Check if it's a direct page load (no X-Requested-With header)
	if r.Header.Get("X-Requested-With") == "" {
		return true
	}

	return false
}

// PageData holds data for SSR pages
type PageData struct {
	Title        string         `json:"title,omitempty"`
	Description  string         `json:"description,omitempty"`
	TopMovies    []models.Movie `json:"topMovies,omitempty"`
	RandomMovies []models.Movie `json:"randomMovies,omitempty"`
	Movies       []models.Movie `json:"movies,omitempty"`
	Movie        *models.Movie  `json:"movie,omitempty"`
	Genres       []models.Genre `json:"genres,omitempty"`
	Query        string         `json:"query,omitempty"`
	Order        string         `json:"order,omitempty"`
	Genre        string         `json:"genre,omitempty"`
}

// HomePage renders the home page with SSR
func (h *SSRHandler) HomePage(w http.ResponseWriter, r *http.Request) {
	if !h.shouldUseSSR(r) {
		// Serve SPA fallback
		http.ServeFile(w, r, filepath.Join(h.publicDir, "index.html"))
		return
	}

	// Fetch data for SSR
	topMoviesOutput, err := h.movieHandler.getTopMoviesUC.Execute()
	if err != nil {
		h.logger.Error("Failed to get top movies for SSR", err)
		http.ServeFile(w, r, filepath.Join(h.publicDir, "index.html"))
		return
	}

	randomMoviesOutput, err := h.movieHandler.getRandomMoviesUC.Execute()
	if err != nil {
		h.logger.Error("Failed to get random movies for SSR", err)
		http.ServeFile(w, r, filepath.Join(h.publicDir, "index.html"))
		return
	}

	// Render HTML with data
	h.renderPage(w, "home", PageData{
		Title:        "Movies - Discover Top Films",
		Description:  "Discover the top movies and find something great to watch today",
		TopMovies:    topMoviesOutput.Movies,
		RandomMovies: randomMoviesOutput.Movies,
	})
}

// MovieDetailsPage renders movie details page with SSR
func (h *SSRHandler) MovieDetailsPage(w http.ResponseWriter, r *http.Request) {
	if !h.shouldUseSSR(r) {
		http.ServeFile(w, r, filepath.Join(h.publicDir, "index.html"))
		return
	}

	// Extract movie ID from path
	path := strings.TrimPrefix(r.URL.Path, "/movies/")
	idStr := strings.TrimSuffix(path, "/")

	id, ok := h.movieHandler.parseID(w, idStr)
	if !ok {
		http.NotFound(w, r)
		return
	}

	// Fetch movie data
	input := movie.GetMovieByIDInput{ID: id}
	output, err := h.movieHandler.getMovieByIDUC.Execute(input)
	if err != nil {
		if err == repository.ErrMovieNotFound {
			http.NotFound(w, r)
			return
		}
		h.logger.Error("Failed to get movie for SSR", err)
		http.ServeFile(w, r, filepath.Join(h.publicDir, "index.html"))
		return
	}

	movieData := output.Movie
	title := movieData.Title
	if movieData.Tagline != nil {
		title += " - " + *movieData.Tagline
	}

	description := title
	if movieData.Overview != nil {
		description = *movieData.Overview
	}

	h.renderPage(w, "movie-details", PageData{
		Title:       title,
		Description: description,
		Movie:       &movieData,
	})
}

// MoviesPage renders search results page with SSR
func (h *SSRHandler) MoviesPage(w http.ResponseWriter, r *http.Request) {
	if !h.shouldUseSSR(r) {
		http.ServeFile(w, r, filepath.Join(h.publicDir, "index.html"))
		return
	}

	query := r.URL.Query().Get("q")
	order := r.URL.Query().Get("order")
	genreStr := r.URL.Query().Get("genre")

	var genre *int
	if genreStr != "" {
		genreInt, ok := h.movieHandler.parseID(w, genreStr)
		if !ok {
			http.ServeFile(w, r, filepath.Join(h.publicDir, "index.html"))
			return
		}
		genre = &genreInt
	}

	// Fetch genres for filter
	genresOutput, err := h.movieHandler.getGenresUC.Execute()
	if err != nil {
		h.logger.Error("Failed to get genres for SSR", err)
	}

	// Search movies
	searchInput := movie.SearchMoviesInput{
		Query: query,
		Order: order,
		Genre: genre,
	}

	searchOutput, err := h.movieHandler.searchMoviesUC.Execute(searchInput)
	if err != nil {
		h.logger.Error("Failed to search movies for SSR", err)
		http.ServeFile(w, r, filepath.Join(h.publicDir, "index.html"))
		return
	}

	title := "Movies"
	if query != "" {
		title = "'" + query + "' movies"
	}

	h.renderPage(w, "movies", PageData{
		Title:  title,
		Movies: searchOutput.Movies,
		Genres: genresOutput.Genres,
		Query:  query,
		Order:  order,
		Genre:  genreStr,
	})
}

// renderPage renders an HTML page with SSR data
func (h *SSRHandler) renderPage(w http.ResponseWriter, pageType string, data PageData) {
	// Read the base HTML template
	htmlPath := filepath.Join(h.publicDir, "index.html")
	htmlBytes, err := os.ReadFile(htmlPath)
	if err != nil {
		h.logger.Error("Failed to read HTML template", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	html := string(htmlBytes)

	// Inject SSR data as JSON in a script tag for hydration
	ssrData := map[string]interface{}{
		"pageType": pageType,
		"data":     data,
	}

	// Convert data to JSON
	jsonBytes, err := json.Marshal(ssrData)
	if err != nil {
		h.logger.Error("Failed to marshal SSR data", err)
	} else {
		// For JSON in <script type="application/json">, we only need to escape </script>
		// Do NOT use HTMLEscapeString as it corrupts the JSON
		jsonData := strings.ReplaceAll(string(jsonBytes), "</script>", "<\\/script>")
		// Inject before closing </body>
		ssrScript := `<script id="ssr-data" type="application/json">` + jsonData + `</script>`
		html = strings.Replace(html, "</body>", ssrScript+"</body>", 1)
	}

	// Update title and meta description
	if data.Title != "" {
		html = strings.Replace(html, "<title>Moovies</title>", "<title>"+template.HTMLEscapeString(data.Title)+"</title>", 1)
	}

	if data.Description != "" {
		metaDesc := `<meta name="description" content="` + template.HTMLEscapeString(data.Description) + `">`
		html = strings.Replace(html, "</head>", metaDesc+"</head>", 1)
	}

	// Inject pre-rendered content into <main>
	mainContent := h.renderMainContent(pageType, data)
	html = strings.Replace(html, "<main></main>", "<main>"+mainContent+"</main>", 1)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

// renderMainContent renders the main content for each page type
func (h *SSRHandler) renderMainContent(pageType string, data PageData) string {
	switch pageType {
	case "home":
		return h.renderHomeContent(data)
	case "movie-details":
		return h.renderMovieDetailsContent(data)
	case "movies":
		return h.renderMoviesContent(data)
	default:
		return ""
	}
}

// renderHomeContent renders the home page content
func (h *SSRHandler) renderHomeContent(data PageData) string {
	var html strings.Builder

	html.WriteString(`<section class="vertical-scroll" id="top-10"><h2>Top 10 This Week</h2><ul>`)
	for _, movie := range data.TopMovies {
		html.WriteString(h.renderMovieItem(movie))
	}
	html.WriteString(`</ul></section>`)

	html.WriteString(`<section class="vertical-scroll" id="random"><h2>Discover Something New</h2><ul>`)
	for _, movie := range data.RandomMovies {
		html.WriteString(h.renderMovieItem(movie))
	}
	html.WriteString(`</ul></section>`)

	return html.String()
}

// renderMovieDetailsContent renders movie details page content
func (h *SSRHandler) renderMovieDetailsContent(data PageData) string {
	if data.Movie == nil {
		return ""
	}

	m := data.Movie
	var html strings.Builder
	html.WriteString(`<article id="movie">`)

	html.WriteString(`<h2>` + template.HTMLEscapeString(m.Title) + `</h2>`)
	if m.Tagline != nil && *m.Tagline != "" {
		html.WriteString(`<h3>` + template.HTMLEscapeString(*m.Tagline) + `</h3>`)
	}

	html.WriteString(`<header>`)

	if m.PosterURL != nil && *m.PosterURL != "" {
		html.WriteString(`<img src="` + template.HTMLEscapeString(*m.PosterURL) + `" alt="` + template.HTMLEscapeString(m.Title) + ` Poster" />`)
	}

	if m.TrailerURL != nil && *m.TrailerURL != "" {
		html.WriteString(`<youtube-embed id="trailer" data-url="` + template.HTMLEscapeString(*m.TrailerURL) + `">YouTube loading...</youtube-embed>`)
	}

	html.WriteString(`<section id="actions" data-id="` + strconv.Itoa(m.ID) + `"><dl id="metadata">`)

	html.WriteString(`<dt>Release Year</dt><dd>` + strconv.Itoa(m.ReleaseYear) + `</dd>`)
	if m.Score != nil {
		html.WriteString(`<dt>Score</dt><dd>` + fmt.Sprintf("%.1f", *m.Score) + ` / 10</dd>`)
	}
	if m.Popularity != nil {
		html.WriteString(`<dt>Popularity</dt><dd>` + fmt.Sprintf("%.1f", *m.Popularity) + `</dd>`)
	}

	html.WriteString(`</dl><button id="btnFavorites">Add to Favorites</button><button id="btnWatchlist">Add to Watchlist</button></section></header>`)

	if len(m.Genres) > 0 {
		html.WriteString(`<ul id="genres">`)
		for _, genre := range m.Genres {
			html.WriteString(`<li>` + template.HTMLEscapeString(genre.Name) + `</li>`)
		}
		html.WriteString(`</ul>`)
	}

	if m.Overview != nil && *m.Overview != "" {
		html.WriteString(`<p id="overview">` + template.HTMLEscapeString(*m.Overview) + `</p>`)
	}

	if len(m.Casting) > 0 {
		html.WriteString(`<ul id="cast">`)
		for _, actor := range m.Casting {
			html.WriteString(`<li>`)
			if actor.ImageURL != nil && *actor.ImageURL != "" {
				html.WriteString(`<img src="` + template.HTMLEscapeString(*actor.ImageURL) + `" alt="` + template.HTMLEscapeString(actor.FirstName+" "+actor.LastName) + `" />`)
			} else {
				html.WriteString(`<img src="/images/generic_actor.jpg" alt="` + template.HTMLEscapeString(actor.FirstName+" "+actor.LastName) + `" />`)
			}
			html.WriteString(`<p>` + template.HTMLEscapeString(actor.FirstName+" "+actor.LastName) + `</p>`)
			html.WriteString(`</li>`)
		}
		html.WriteString(`</ul>`)
	}

	html.WriteString(`</article>`)
	return html.String()
}

// renderMoviesContent renders search results page content
func (h *SSRHandler) renderMoviesContent(data PageData) string {
	var html strings.Builder

	html.WriteString(`<section><div id="search-header">`)
	if data.Query != "" {
		html.WriteString(`<h2>'` + template.HTMLEscapeString(data.Query) + `' movies</h2>`)
	} else {
		html.WriteString(`<h2>Movies</h2>`)
	}

	html.WriteString(`<section id="filters"><select id="filter" onchange="app.searchFilterChange(this.value)"><option>Filter by Genre</option>`)
	for _, genre := range data.Genres {
		selected := ""
		if strconv.Itoa(genre.ID) == data.Genre {
			selected = " selected"
		}
		html.WriteString(`<option value="` + strconv.Itoa(genre.ID) + `"` + selected + `>` + template.HTMLEscapeString(genre.Name) + `</option>`)
	}
	html.WriteString(`</select><select id="order" onchange="app.searchOrderChange(this.value)">`)
	orders := []struct{ value, label string }{
		{"popularity", "Sort by Popularity"},
		{"score", "Sort by Score"},
		{"date", "Sort by Release Date"},
		{"name", "Sort by Name"},
	}
	for _, o := range orders {
		selected := ""
		if o.value == data.Order {
			selected = " selected"
		}
		html.WriteString(`<option value="` + o.value + `"` + selected + `>` + o.label + `</option>`)
	}
	html.WriteString(`</select></section></div><ul id="movies-result">`)

	for _, movie := range data.Movies {
		html.WriteString(h.renderMovieItem(movie))
	}

	html.WriteString(`</ul></section>`)
	return html.String()
}

// renderMovieItem renders a single movie item
func (h *SSRHandler) renderMovieItem(movie models.Movie) string {
	var html strings.Builder
	html.WriteString(`<li><movie-item>`)

	posterURL := ""
	if movie.PosterURL != nil {
		posterURL = *movie.PosterURL
	}

	html.WriteString(`<a href="/movies/` + strconv.Itoa(movie.ID) + `" class="navlink">`)
	html.WriteString(`<article>`)
	if posterURL != "" {
		html.WriteString(`<img src="` + template.HTMLEscapeString(posterURL) + `" alt="` + template.HTMLEscapeString(movie.Title) + ` Poster" />`)
	}
	html.WriteString(`<p>` + template.HTMLEscapeString(movie.Title) + ` (` + strconv.Itoa(movie.ReleaseYear) + `)</p>`)
	html.WriteString(`</article></a></movie-item></li>`)

	return html.String()
}
