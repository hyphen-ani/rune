package middleware

import (
	"net/http"
	"rune/internal/auth"
)

// AuthorizeOperation checks whether the token in the request may perform
// `permission` (one of auth.PermRead / PermWrite / PermDelete / PermRotate)
// on `namespace`.
//
// Decision table:
//   - Root token → always allowed.
//   - Token with Namespaces set (new format) → checks namespace + permission in the slice.
//   - Legacy token (Namespaces nil, Namespace field set) → full access except ReadOnly tokens get read-only.
func AuthorizeOperation(r *http.Request, namespace, permission string) bool {
	if IsRoot(r) {
		return true
	}
	token := GetTokenRecord(r)
	if token == nil {
		return false
	}

	// New fine-grained format
	if len(token.Namespaces) > 0 {
		for _, np := range token.Namespaces {
			if np.Namespace == namespace {
				for _, p := range np.Permissions {
					if p == permission {
						return true
					}
				}
				return false
			}
		}
		return false
	}

	// Legacy single-namespace format
	if token.Namespace != namespace {
		return false
	}
	if token.ReadOnly {
		return permission == auth.PermRead
	}
	return true
}

// AuthorizeNamespace returns true if the token has any access to namespace.
// Kept for call sites that only need namespace-level gating (not operation-specific).
func AuthorizeNamespace(r *http.Request, namespace string) bool {
	if IsRoot(r) {
		return true
	}
	token := GetTokenRecord(r)
	if token == nil {
		return false
	}
	if len(token.Namespaces) > 0 {
		for _, np := range token.Namespaces {
			if np.Namespace == namespace && len(np.Permissions) > 0 {
				return true
			}
		}
		return false
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

// IsReadOnly is kept for existing call sites; prefer AuthorizeOperation for new code.
func IsReadOnly(r *http.Request) bool {
	token := GetTokenRecord(r)
	return token != nil && token.ReadOnly
}

func GetPlainToken(r *http.Request) string {
	value := r.Context().Value(PlainTokenContextKey)
	if value == nil {
		return ""
	}
	s, _ := value.(string)
	return s
}
