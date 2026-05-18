package web

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/zrgst/arx-dash/internal/arx"
)

// Server holder dependencies for HTTP-laget.
type Server struct {
	arxCache  *arx.Cache
	mux       *http.ServeMux
	templates *template.Template
}

// PageData er data som sendes til HTML templates.
type PageData struct {
	Title        string
	ActivePage   string
	CacheLoaded  bool
	CacheUpdated string
	Persons      []arx.Person
	Cards        []arx.Card
}

// NewServer setter opp routes.
func NewServer(arxCache *arx.Cache) *Server {
	funcs := template.FuncMap{
		"mask":  maskString,
		"yesNo": yesNo,
	}

	templates := template.Must(
		template.New("").Funcs(funcs).ParseGlob(filepath.Join("web", "templates", "*.html")),
	)

	s := &Server{
		arxCache:  arxCache,
		mux:       http.NewServeMux(),
		templates: templates,
	}

	s.routes()

	return s
}

// Handler returnerer http.Handler slik at main kan starte serveren.
func (s *Server) Handler() http.Handler {
	return s.mux
}

func (s *Server) routes() {
	// Static CSS.
	s.mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	// HTML pages.
	s.mux.HandleFunc("GET /", s.handleHomePage)
	s.mux.HandleFunc("GET /api/personer", s.handlePersonsPage)
	s.mux.HandleFunc("GET /api/kort", s.handleCardsPage)

	// HTML form refresh.
	s.mux.HandleFunc("POST /refresh", s.handleRefreshPage)

	// JSON API.
	s.mux.HandleFunc("GET /api/persons-export", s.handlePersonsExportAPI)
	s.mux.HandleFunc("GET /api/persons", s.handlePersonsAPI)
	s.mux.HandleFunc("GET /api/cards", s.handleCardsAPI)
	s.mux.HandleFunc("POST /api/refresh", s.handleRefreshAPI)
	s.mux.HandleFunc("GET /api/cache-status", s.handleCacheStatusAPI)
}

func (s *Server) handleHomePage(w http.ResponseWriter, r *http.Request) {
	data := s.basePageData("Start", "home")
	s.render(w, "home.html", data)
}

func (s *Server) handlePersonsExport(w http.ResponseWriter, r *http.Request) {
	export, err := s.arxCache.GetPersonsExport(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, export)
}

func (s *Server) handlePersonsPage(w http.ResponseWriter, r *http.Request) {
	export, err := s.arxCache.GetPersonsExport(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}

	data := s.basePageData("Personer", "persons")
	data.Persons = export.Persons

	s.render(w, "persons.html", data)
}

func (s *Server) handleCardsPage(w http.ResponseWriter, r *http.Request) {
	export, err := s.arxCache.GetPersonsExport(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}

	data := s.basePageData("Kort", "cards")
	data.Cards = export.Cards

	s.render(w, "cards.html", data)
}

func (s *Server) handleRefreshPage(w http.ResponseWriter, r *http.Request) {
	_, err := s.arxCache.RefreshPersonsExport(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}

	// Send brukeren tilbake til siden de kom fra.
	redirectTo := r.Header.Get("Referer")
	if redirectTo == "" {
		redirectTo = "/"
	}

	http.Redirect(w, r, redirectTo, http.StatusSeeOther)
}

func (s *Server) handlePersonsExportAPI(w http.ResponseWriter, r *http.Request) {
	export, err := s.arxCache.GetPersonsExport(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, export)
}

func (s *Server) handlePersonsAPI(w http.ResponseWriter, r *http.Request) {
	export, err := s.arxCache.GetPersonsExport(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, export.Persons)
}

func (s *Server) handleCardsAPI(w http.ResponseWriter, r *http.Request) {
	export, err := s.arxCache.GetPersonsExport(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, export.Cards)
}

func (s *Server) handleRefreshAPI(w http.ResponseWriter, r *http.Request) {
	export, err := s.arxCache.RefreshPersonsExport(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, map[string]any{
		"ok":        true,
		"timestamp": export.Timestamp,
		"persons":   len(export.Persons),
		"cards":     len(export.Cards),
		"updatedAt": s.arxCache.UpdatedAt(),
	})
}

func (s *Server) handleCacheStatusAPI(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, map[string]any{
		"loaded":    s.arxCache.Loaded(),
		"updatesAt": s.arxCache.UpdatedAt(),
	})
}

func (s *Server) basePageData(title string, activePage string) PageData {
	updatedAt := ""

	if s.arxCache.Loaded() {
		updatedAt = s.arxCache.UpdatedAt().Format("2006-01-02 15:04:05")
	}

	return PageData{
		Title:        title,
		ActivePage:   activePage,
		CacheLoaded:  s.arxCache.Loaded(),
		CacheUpdated: updatedAt,
	}
}

func (s *Server) render(w http.ResponseWriter, templateName string, data PageData) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err := s.templates.ExecuteTemplate(w, templateName, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func writeJSON(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")

	err := encoder.Encode(value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func writeError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func maskString(value string) string {
	if value == "" {
		return ""
	}

	return "****"
}

func yesNo(value bool) string {
	if value {
		return "Ja"
	}

	return "Nei"
}
