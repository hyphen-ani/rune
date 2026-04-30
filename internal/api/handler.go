package api

import (
	"encoding/json"
	"net/http"
	"rune/internal/crypto"
	"rune/internal/seal"
	"rune/internal/service"
	"rune/internal/storage"
)

type Handler struct {
	secretService *service.SecretService
	store         storage.Store
	sealer        *seal.Manager
}

func NewHandler(s *service.SecretService, store storage.Store, sealer *seal.Manager) *Handler {
	return &Handler{
		secretService: s,
		store:         store,
		sealer:        sealer,
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

	err := h.secretService.Put(req.Key, req.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
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

func (h *Handler) Unseal(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Passphrase string `json:"passphrase"`
	}

	json.NewDecoder(r.Body).Decode(&req)
	salt, _ := h.store.GetSalt()
	key := crypto.DeriveKey(req.Passphrase, salt)

	record, err := h.store.Get(storage.VerifyKey)
	if err != nil {
		http.Error(w, "vault not initialized", http.StatusInternalServerError)
		return
	}

	plaintext, err := crypto.Decrypt(key, record.Ciphertext, record.Nonce)
	if err != nil || string(plaintext) != "rune-check" {
		http.Error(w, "invalid passphrase", http.StatusUnauthorized)
		return
	}

	h.sealer.Unseal(key)
	w.Write([]byte("unsealed"))

}

func (h *Handler) Seal(w http.ResponseWriter, r *http.Request) {
	h.sealer.Seal()
	w.Write([]byte("sealed"))
}

func (h *Handler) Status(w http.ResponseWriter, r *http.Request) {
	if h.sealer.IsSealed() {
		w.Write([]byte("sealed"))
	} else {
		w.Write([]byte("unsealed"))
	}
}
