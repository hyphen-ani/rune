package api

import (
	"encoding/json"
	"net/http"
	"rune/internal/service"
)

type Handler struct {
	secretService *service.SecretService
}

func NewHandler(s *service.SecretService) *Handler {
	return &Handler{
		secretService: s,
	}
}

func (h *Handler) PutSecret(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
	}

	h.secretService.Put(req.Key, req.Value)
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) GetSecret(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")

	val, err := h.secretService.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"value": val,
	})
}
