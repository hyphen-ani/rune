package middleware

import (
	"net/http"
	"rune/internal/auth"
)

func AuthorizeNamespace(r *http.Request, namespace string) bool {

	value := r.Context().Value(TokenContextKey)
	if value == nil {
		return false
	}
	token, ok := value.(*auth.TokenRecord)
	if !ok {
		return false
	}
	if IsRoot(r) {
		return true
	}
	return token.Namespace == namespace

}

func GetTokenRecord(r *http.Request) *auth.TokenRecord {
	value := r.Context().Value(TokenContextKey)
	if value == nil {
		return nil
	}
	token, _ := value.(*auth.TokenRecord)
	return token
}

func IsRoot(r *http.Request) bool {

	value := r.Context().Value(TokenContextKey)
	if value == nil {
		return false
	}
	token, ok := value.(*auth.TokenRecord)
	if !ok {
		return false
	}
	return token.ID == "root" && token.Namespace == "*"

}

func IsReadOnly(r *http.Request) bool {
	value := r.Context().Value(TokenContextKey)
	if value == nil {
		return false
	}
	token, ok := value.(*auth.TokenRecord)
	if !ok {
		return false
	}
	return token.ReadOnly
}

func GetPlainToken(r *http.Request) string {
	value := r.Context().Value(PlainTokenContextKey)
	if value == nil {
		return ""
	}
	s, _ := value.(string)
	return s
}
