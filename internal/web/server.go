package web

import (
	"encoding/json"
	"net/http"

	"github.com/zrgst/arx-dash/internal/arx"
)

// Server holder dependencies for HTTP-laget.
type Server struct {
	arx *arx.Client
	mux *http.ServeMux
}

// NewServer setter opp routes.
func NewServer(arxClient *arx.Client) *Server {
	s := &Server{
		arx: arxClient,
		mux: http.NewServeMux(),
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
}

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	_, _ = w.Write([]byte(`
<!doctype html>
<html lang="no">
  <head>
    <meta charset="UTF-8">
    <title>ARX Dashboard Go</title>
  </head>
  <body>
    <h1>ARX Dashboard Go</h1>
    <p>Første Go-versjon kjører.</p>

    <ul>
      <li><a href="/api/persons-export">/api/persons-export</a></li>
      <li><a href="/api/persons">/api/persons</a></li>
      <li><a href="/api/cards">/api/cards</a></li>
    </ul>
  </body>
</html>
`))
}

func (s *Server) handlePersonsExport(w http.ResponseWriter, r *http.Request) {
	export, err := s.arx.ExportPersons(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, export)
}

func (s *Server) handlePersons(w http.ResponseWriter, r *http.Request) {
	export, err := s.arx.ExportPersons(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, export.Persons)
}

func (s *Server) handleCards(w http.ResponseWriter, r *http.Request) {
	export, err := s.arx.ExportPersons(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, export.Cards)
}

func writeJSON(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(value); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func writeError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
