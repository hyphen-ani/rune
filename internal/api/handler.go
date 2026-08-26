package api

import (
	"encoding/json"
	"net/http"
	"rune/internal/crypto"
	"rune/internal/middleware"
	"rune/internal/namespace"
	"rune/internal/seal"
	"rune/internal/service"
	"rune/internal/storage"
)

type Handler struct {
	namespaceService *service.NamespaceService
	secretService    *service.SecretService
	tokenService     *service.TokenService
	store            storage.Store
	sealer           *seal.Manager
}

func NewHandler(s *service.SecretService, t *service.TokenService, ns *service.NamespaceService,
	store storage.Store, sealer *seal.Manager) *Handler {
	return &Handler{
		secretService:    s,
		tokenService:     t,
		namespaceService: ns,
		store:            store,
		sealer:           sealer,
	}
}

// CORE VAULT

func (h *Handler) PutSecret(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Key       string `json:"key"`
		Value     string `json:"value"`
		Namespace string `json:"namespace"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	namespaceName := namespace.Normalize(req.Namespace)

	if !middleware.AuthorizeNamespace(r, namespaceName) {
		http.Error(w, "[AUTHORIZATION DENIED]: Token cannot access this namespace", http.StatusForbidden)
		return
	}

	err := h.secretService.Put(namespaceName, req.Key, req.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
func (h *Handler) GetSecret(w http.ResponseWriter, r *http.Request) {

	namespaceName := namespace.Normalize(r.URL.Query().Get("namespace"))
	key := r.URL.Query().Get("key")

	if !middleware.AuthorizeNamespace(r, namespaceName) {
		http.Error(w, "[AUTHORIZATION DENIED]: Token cannot access this namespace", http.StatusForbidden)
		return
	}

	val, err := h.secretService.Get(namespaceName, key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"value": val,
	})
}
func (h *Handler) ListSecrets(w http.ResponseWriter, r *http.Request) {
	namespaceName := namespace.Normalize(r.URL.Query().Get("namespace"))

	if !middleware.AuthorizeNamespace(r, namespaceName) {
		http.Error(w, "[AUTHORIZATION DENIED]: Token cannot access this namespace", http.StatusForbidden)
		return
	}

	keys, err := h.secretService.ListKeys(namespaceName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	json.NewEncoder(w).Encode(keys)
}
func (h *Handler) DeleteSecret(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	namespaceName := namespace.Normalize(r.URL.Query().Get("namespace"))

	if !middleware.AuthorizeNamespace(r, namespaceName) {
		http.Error(w, "[AUTHORIZATION DENIED]: Token cannot access this namespace", http.StatusForbidden)
		return
	}

	err := h.secretService.Delete(namespaceName, key)

	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
func (h *Handler) RotateSecret(w http.ResponseWriter, r *http.Request) {
	namespaceName := namespace.Normalize(r.URL.Query().Get("namespace"))
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	if !middleware.AuthorizeNamespace(r, namespaceName) {
		http.Error(w, "[AUTHORIZATION DENIED]: Token cannot access this namespace", http.StatusForbidden)
		return
	}

	value, err := h.secretService.Rotate(namespaceName, key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"value": value,
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
		http.Error(w, "Vault Not Initialized", http.StatusInternalServerError)
		return
	}

	plaintext, err := crypto.Decrypt(key, record.Ciphertext, record.Nonce)
	if err != nil || string(plaintext) != "rune-check" {
		http.Error(w, "Invalid Passphrase", http.StatusUnauthorized)
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

// TOKEN HANDLERS

func (h *Handler) CreateToken(w http.ResponseWriter, r *http.Request) {

	if !middleware.IsRoot(r) {
		http.Error(w, "[AUTHORIZATION DENIED]: Root Access Required", http.StatusForbidden)
		return
	}

	var req struct {
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	token, record, err := h.tokenService.Create(req.Name, req.Namespace)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	resp := map[string]interface{}{
		"token":  token,
		"record": record,
	}

	json.NewEncoder(w).Encode(resp)
}
func (h *Handler) ListTokens(w http.ResponseWriter, r *http.Request) {

	if !middleware.IsRoot(r) {
		http.Error(w, "[AUTHORIZATION DENIED]: Root Access Required", http.StatusForbidden)
		return
	}

	tokens, err := h.tokenService.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tokens)

}
func (h *Handler) RevokeToken(w http.ResponseWriter, r *http.Request) {

	if !middleware.IsRoot(r) {
		http.Error(w, "[AUTHORIZATION DENIED]: Root Access Required", http.StatusForbidden)
		return
	}

	if !middleware.IsRoot(r) {
		http.Error(w, "[AUTHORIZATION DENIED]: Root Access Required", http.StatusForbidden)
		return
	}

	var req struct {
		ID string `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	err := h.tokenService.Revoke(req.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	w.Write([]byte("revoked"))
}

// NAMESPACE HANDLERS

func (h *Handler) CreateNamespace(w http.ResponseWriter, r *http.Request) {

	if !middleware.IsRoot(r) {
		http.Error(w, "[AUTHORIZATION DENIED]: Root Access Required", http.StatusForbidden)
		return
	}

	if h.sealer.IsSealed() {
		http.Error(w, "vault is sealed", http.StatusForbidden)
		return
	}

	var req struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	err := h.store.CreateNamespace(req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("namespace created"))
}
func (h *Handler) ListNamespaces(w http.ResponseWriter, r *http.Request) {

	if !middleware.IsRoot(r) {
		http.Error(w, "[AUTHORIZATION DENIED]: Root Access Required", http.StatusForbidden)
		return
	}

	if h.sealer.IsSealed() {
		http.Error(w, "vault is sealed", http.StatusForbidden)
		return
	}

	namespaces, err := h.store.ListNamespaces()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(namespaces)
}
func (h *Handler) DeleteNamespace(w http.ResponseWriter, r *http.Request) {

	if h.sealer.IsSealed() {
		http.Error(w, "vault is sealed", http.StatusForbidden)
		return
	}

	var req struct {
		Name string `json:"name"`
	}

	json.NewDecoder(r.Body).Decode(&req)
	err := h.namespaceService.DeleteNamespace(req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	w.Write([]byte("namespace deleted"))
}
