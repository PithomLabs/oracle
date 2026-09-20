package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	coordinator "github.com/PithomLabs/oracle/coordinator"
)

// Handler is the HTTP handler for the Coordinator.
type Handler struct {
	coordinator *coordinator.Coordinator
}

// NewHandler creates a new HTTP handler.
func NewHandler(c *coordinator.Coordinator) *Handler {
	return &Handler{coordinator: c}
}

// ServeHTTP implements the http.Handler interface.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "POST" && r.URL.Path == "/decisions":
		h.handleSubmitDecision(w, r)
	case r.Method == "GET" && strings.HasPrefix(r.URL.Path, "/packets/"):
		h.handlePacketStatus(w, r)
	case r.Method == "GET" && strings.HasPrefix(r.URL.Path, "/beliefs/") && strings.HasSuffix(r.URL.Path, "/decision-context"):
		h.handleDecisionContext(w, r)
	case r.Method == "GET" && strings.HasPrefix(r.URL.Path, "/authorization-context/"):
		h.handleAuthorizationContext(w, r)
	case r.Method == "GET" && strings.HasPrefix(r.URL.Path, "/context/"):
		h.handleGetContext(w, r)
	case r.Method == "GET" && r.URL.Path == "/packs/rules":
		h.handleGetPackRules(w, r)
	default:
		http.Error(w, "not found", http.StatusNotFound)
	}
}

func (h *Handler) requireAuth(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("missing Authorization header")
	}
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", fmt.Errorf("invalid Authorization header: must use Bearer scheme")
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		return "", fmt.Errorf("empty bearer token")
	}
	operatorID, err := h.coordinator.AuthenticateOperator(token)
	if err != nil {
		return "", fmt.Errorf("authentication failed: %w", err)
	}
	return operatorID, nil
}

func (h *Handler) handleSubmitDecision(w http.ResponseWriter, r *http.Request) {
	// Authenticate operator before processing decision
	if _, err := h.requireAuth(r); err != nil {
		writeError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var req SubmitDecisionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	decisionReq := coordinator.DecisionRequest{
		Type:          coordinator.DecisionType(req.Type),
		BeliefID:      req.BeliefID,
		ScenarioID:    req.ScenarioID,
		DebtItem:      req.DebtItem,
		EvidenceClass: req.EvidenceClass,
		Reason:        req.Reason,
	}

	record, err := h.coordinator.SubmitDecision(decisionReq)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, SubmitDecisionResponse{
		ID:            record.ID,
		Type:          string(record.Type),
		BeliefID:      record.BeliefID,
		Result:        record.Result,
		RefusalReason: record.RefusalReason,
		CreatedAt:     record.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, http.StatusOK)
}

func (h *Handler) handlePacketStatus(w http.ResponseWriter, r *http.Request) {
	packetID := strings.TrimPrefix(r.URL.Path, "/packets/")
	packetID = strings.TrimSuffix(packetID, "/status")

	status, err := h.coordinator.PacketStatus(packetID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, PacketStatusResponse{
		PacketID:  status.PacketID,
		Status:    status.Status,
		BeliefIDs: status.BeliefIDs,
	}, http.StatusOK)
}

func (h *Handler) handleDecisionContext(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/beliefs/")
	beliefID := strings.TrimSuffix(path, "/decision-context")

	ctx, err := h.coordinator.DecisionContext(beliefID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, DecisionContextResponse{
		BeliefID:    ctx.BeliefID,
		Claim:       ctx.Claim,
		Status:      ctx.Status,
		Debt:        ctx.Debt,
		EvidenceIDs: ctx.EvidenceIDs,
	}, http.StatusOK)
}

func (h *Handler) handleAuthorizationContext(w http.ResponseWriter, r *http.Request) {
	targetID := strings.TrimPrefix(r.URL.Path, "/authorization-context/")

	projection, err := h.coordinator.AuthorizationContext(targetID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, AuthorizationContextResponse{
		Target:         projection.Target,
		BeliefStatus:   projection.BeliefStatus,
		AuthorityState: projection.AuthorityState,
		CurrentResult:  projection.CurrentResult,
		RefusalReason:  projection.RefusalReason,
	}, http.StatusOK)
}

func (h *Handler) handleGetContext(w http.ResponseWriter, r *http.Request) {
	taskID := strings.TrimPrefix(r.URL.Path, "/context/")
	if taskID == "" {
		writeError(w, "task_id required", http.StatusBadRequest)
		return
	}

	ctx, err := h.coordinator.GetContext(taskID)
	if err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}

	writeJSON(w, ctx, http.StatusOK)
}

func (h *Handler) handleGetPackRules(w http.ResponseWriter, r *http.Request) {
	// Get retirement rules from all registered packs
	rules := h.coordinator.GetPackRules()
	writeJSON(w, rules, http.StatusOK)
}

func writeJSON(w http.ResponseWriter, v interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, message string, status int) {
	writeJSON(w, ErrorResponse{Error: message}, status)
}
