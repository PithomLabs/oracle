package ui

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/PithomLabs/oracle/internal/application"
)

//go:embed templates/*.html
var templateFS embed.FS

// Server is the embedded Trust UI.
type Server struct {
	app           *application.App
	tmpl          *template.Template
	operatorToken string
}

// NewServer creates a new Trust UI server.
func NewServer(app *application.App, operatorToken string) *Server {
	tmpl := template.Must(template.ParseFS(templateFS, "templates/*.html"))
	return &Server{app: app, tmpl: tmpl, operatorToken: operatorToken}
}

// RegisterRoutes registers the Trust UI routes on the given mux.
func (s *Server) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/ui", s.HandleIndex)
	mux.HandleFunc("/ui/insights", s.HandleInsights)
	mux.HandleFunc("/ui/debts", s.HandleDebts)
	mux.HandleFunc("/ui/api/login", s.HandleLogin)
	mux.HandleFunc("/ui/api/discharge", s.requireAuth(s.requireOrigin(s.HandleDischarge)))
	mux.HandleFunc("/ui/api/promote", s.requireAuth(s.requireOrigin(s.HandlePromote)))
}

// HandleLogin validates the operator token and sets an HttpOnly session cookie.
func (s *Server) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Token != s.operatorToken {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "argus_session",
		Value:    s.operatorToken,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	w.WriteHeader(http.StatusOK)
}

// requireAuth checks cookie OR Bearer header.
func (s *Server) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("argus_session")
		if err == nil && cookie.Value == s.operatorToken {
			next(w, r)
			return
		}
		authHeader := r.Header.Get("Authorization")
		if authHeader == "Bearer "+s.operatorToken {
			next(w, r)
			return
		}
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	}
}

// requireOrigin validates Origin/Host headers for localhost.
func (s *Server) requireOrigin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			u, err := url.Parse(origin)
			if err != nil || (u.Hostname() != "localhost" && u.Hostname() != "127.0.0.1") {
				http.Error(w, "forbidden: invalid origin", http.StatusForbidden)
				return
			}
		}
		host := r.Host
		if host != "" {
			hostname := strings.Split(host, ":")[0]
			if hostname != "localhost" && hostname != "127.0.0.1" {
				http.Error(w, "forbidden: invalid host", http.StatusForbidden)
				return
			}
		}
		next(w, r)
	}
}

func (s *Server) HandleIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/ui/insights", http.StatusFound)
}

func (s *Server) HandleInsights(w http.ResponseWriter, r *http.Request) {
	dash, err := s.app.GetDashboard(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to load dashboard: %v", err), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Tasks":   dash.Tasks,
		"Beliefs": dash.Beliefs,
	}

	var buf bytes.Buffer
	if err := s.tmpl.ExecuteTemplate(&buf, "insights.html", data); err != nil {
		http.Error(w, "template error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	buf.WriteTo(w)
}

func (s *Server) HandleDebts(w http.ResponseWriter, r *http.Request) {
	dash, err := s.app.GetDashboard(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to load dashboard: %v", err), http.StatusInternalServerError)
		return
	}

	type beliefExplain struct {
		BeliefID      string
		Claim         string
		RemainingDebt string
		CanPromote    bool
	}

	var beliefs []beliefExplain
	for _, b := range dash.Beliefs {
		beliefs = append(beliefs, beliefExplain{
			BeliefID:      b.ID,
			Claim:         b.Claim,
			RemainingDebt: strings.Join(b.Debt, ", "),
			CanPromote:    len(b.Debt) == 0 && b.Status == "entered",
		})
	}

	data := map[string]interface{}{
		"Beliefs":     beliefs,
		"Retirements": []interface{}{},
	}

	var buf bytes.Buffer
	if err := s.tmpl.ExecuteTemplate(&buf, "debts.html", data); err != nil {
		http.Error(w, "template error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	buf.WriteTo(w)
}

func (s *Server) HandleDischarge(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var extReq application.ExternalDecisionRequest
	if err := json.NewDecoder(r.Body).Decode(&extReq); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Server-derived principal — never from request body
	principalID := application.LocalOperatorPrincipalID

	err := s.app.SubmitDecision(r.Context(), &application.AuthenticatedDecisionCommand{
		Type:          "discharge",
		ScenarioID:    extReq.ScenarioID,
		BeliefID:      extReq.BeliefID,
		ObligationKey: extReq.ObligationKey,
		InstrumentRef: extReq.InstrumentRef,
		PrincipalID:   principalID,
		EvidenceClass: extReq.EvidenceClass,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("discharge failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, `{"status":"ok"}`)
}

func (s *Server) HandlePromote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ScenarioID string `json:"scenario_id"`
		BeliefID   string `json:"belief_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err := s.app.SubmitDecision(r.Context(), &application.AuthenticatedDecisionCommand{
		Type:       "promote",
		ScenarioID: req.ScenarioID,
		BeliefID:   req.BeliefID,
		PrincipalID: application.LocalOperatorPrincipalID,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("promote failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, `{"status":"ok"}`)
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
