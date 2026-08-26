package middleware

import (
	"log"
	"net/http"
	"rune/internal/auth"
)

func AuthorizeNamespace(r *http.Request, namespace string) bool {

	value := r.Context().Value(TokenContextKey)
	if value == nil {
		log.Println("AUTHZ DEBUG -> TokenRecord NOT FOUND in context")
		return false
	}
	token, ok := value.(*auth.TokenRecord)
	if !ok {
		log.Printf(
			"AUTHZ DEBUG -> Invalid context type: %T",
			value,
		)
		return false
	}
	log.Printf(
		"AUTHZ DEBUG -> ID=%q Namespace=%q Requested=%q",
		token.ID,
		token.Namespace,
		namespace,
	)
	if IsRoot(r) {
		log.Println("AUTHZ DEBUG -> ROOT ACCESS GRANTED")
		return true
	}
	return token.Namespace == namespace

}

func IsRoot(r *http.Request) bool {

	value := r.Context().Value(TokenContextKey)
	if value == nil {
		log.Println("ROOT DEBUG -> TokenRecord NOT FOUND in context")
		return false
	}
	token, ok := value.(*auth.TokenRecord)
	if !ok {
		log.Printf(
			"ROOT DEBUG -> Invalid context type: %T",
			value,
		)
		return false
	}
	log.Printf(
		"ROOT DEBUG -> ID=%q Namespace=%q",
		token.ID,
		token.Namespace,
	)
	return token.ID == "root" && token.Namespace == "*"

}
