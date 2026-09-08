package handlers

import (
	"encoding/json"
	"net/http"
	"net/mail"
	"strings"

	"portfolio-backend/internal/data"
	"portfolio-backend/internal/mailer"
)

type Handler struct {
	store  *data.Store
	mailer mailer.Mailer
}

func New(store *data.Store, m mailer.Mailer) *Handler {
	return &Handler{store: store, mailer: m}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) Profile(w http.ResponseWriter, r *http.Request) {
	writeRaw(w, h.store.Profile)
}

func (h *Handler) Skills(w http.ResponseWriter, r *http.Request) {
	writeRaw(w, h.store.Skills)
}

func (h *Handler) Projects(w http.ResponseWriter, r *http.Request) {
	writeRaw(w, h.store.Projects)
}

func (h *Handler) Experience(w http.ResponseWriter, r *http.Request) {
	writeRaw(w, h.store.Experience)
}

type contactRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

const maxContactBodyBytes = 1 << 12 // 4KB is plenty for a contact form

func (h *Handler) Contact(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxContactBodyBytes)

	var req contactRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)
	req.Message = strings.TrimSpace(req.Message)

	if req.Name == "" || req.Message == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "name and message are required"})
		return
	}
	if len(req.Name) > 200 || len(req.Message) > 5000 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "name or message too long"})
		return
	}
	if _, err := mail.ParseAddress(req.Email); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid email address"})
		return
	}

	if err := h.mailer.Send(req.Name, req.Email, req.Message); err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": "failed to send message"})
		return
	}

	writeJSON(w, http.StatusAccepted, map[string]string{"status": "sent"})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeRaw(w http.ResponseWriter, raw []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(raw)
}
