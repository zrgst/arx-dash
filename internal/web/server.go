package web

import (
	"encoding/json"
	"net/http"

	"github.com/zrgst/arx-dash/internal/arx"
)

// Server holder dependencies for HTTP-laget.
type Server struct {
	arxCache *arx.Cache
	mux      *http.ServeMux
}

// NewServer setter opp routes.
func NewServer(arxCache *arx.Cache) *Server {
	s := &Server{
		arxCache: arxCache,
		mux:      http.NewServeMux(),
	}

	s.routes()

	return s
}

// Handler returnerer http.Handler slik at main kan starte serveren.
func (s *Server) Handler() http.Handler {
	return s.mux
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /", s.handleHome)
	s.mux.HandleFunc("GET /api/persons-export", s.handlePersonsExport)
	s.mux.HandleFunc("GET /api/persons", s.handlePersons)
	s.mux.HandleFunc("GET /api/cards", s.handleCards)

	// Tvinger ny henting fra ARX
	s.mux.HandleFunc("POST /api/refresh", s.handleRefresh)

	// Viser cache-status
	s.mux.HandleFunc("GET /api/cache-status", s.handleCacheStatus)
}

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	_, _ = w.Write([]byte(`
		<!doctype html>
		<html lang="no">
			<head>
				<meta charset="UTF-8">
				<title>ARX-Dash</title>
			</head>
			<body>
				<h1>ARX|DASH</h1>
				<p>Første versjon i Go kjører.</p>

				<ul>
					<li><a href="/api/persons-export">/api/persons-export</a></li>
					<li><a href="/api/persons">/api/persons</a></li>
					<li><a href="/api/cards">/api/cards</a></li>
					<li><a href="/apie/cache-status">/api/cache-status</a></li>
				</ul>

				<form method="post" action="/api/refresh">
					<button type="submit">Refresh ARX Cache</button>
				</form>
			</body>
		</html>
		`))
}

func (s *Server) handlePersonsExport(w http.ResponseWriter, r *http.Request) {
	export, err := s.arxCache.GetPersonsExport(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, export)
}

func (s *Server) handlePersons(w http.ResponseWriter, r *http.Request) {
	export, err := s.arxCache.GetPersonsExport(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, export.Persons)
}

func (s *Server) handleCards(w http.ResponseWriter, r *http.Request) {
	export, err := s.arxCache.RefreshPersonsExport(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, export.Cards)
}

func (s *Server) handleRefresh(w http.ResponseWriter, r *http.Request) {
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

func (s *Server) handleCacheStatus(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, map[string]any{
		"loaded":    s.arxCache.Loaded(),
		"updatesAt": s.arxCache.UpdatedAt(),
	})
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
