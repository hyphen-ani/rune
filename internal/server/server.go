package server

import (
	"fmt"
	"log"
	"net/http"
	"rune/internal/api"
	"rune/internal/auth"
	"rune/internal/config"
	"rune/internal/crypto"
	"rune/internal/middleware"
	"rune/internal/seal"
	"rune/internal/service"
	"rune/internal/storage"
	"time"
)

func Start() {
	log.Println("Starting Rune Server on: 8080")

	dbPath, err := config.GetDBPath()
	if err != nil {
		log.Fatal(err)
	}

	store, err := storage.NewBoltStore(dbPath)
	if err != nil {
		log.Fatal(err)
	}

	salt, err := store.GetSalt()
	if err != nil {
		// Initialization
		salt, _ = crypto.GenerateSalt()
		store.SaveSalt(salt)

		var passphrase string
		fmt.Print("Enter Passphrase: ")
		fmt.Scanln(&passphrase)

		key := crypto.DeriveKey(passphrase, salt)

		ciphertext, nonce, _ := crypto.Encrypt(key, []byte("rune-check"))

		store.Put(storage.VerifyKey, storage.SecretRecord{
			Ciphertext: ciphertext,
			Nonce:      nonce,
		})

		token, hash := auth.GenerateToken()
		record := auth.TokenRecord{
			ID:        "root",
			Hash:      hash,
			Name:      "root",
			CreatedAt: time.Now().Format(time.RFC3339),
			Revoked:   false,
		}
		err = store.SaveToken(record)
		if err != nil {
			return
		}

		fmt.Println("====================================")
		fmt.Println("Root Token (SAVE THIS):", token)
		fmt.Println("====================================")
		fmt.Println("Store this token securely. It will not be shown again.")

	}

	sealer := seal.NewManger()
	secretService := service.NewSecretService(store, sealer)
	tokenService := service.NewTokenService(store)
	handler := api.NewHandler(secretService, tokenService, store, sealer)

	// PUBLIC
	public := http.NewServeMux()
	public.HandleFunc("/unseal", handler.Unseal)
	public.HandleFunc("/status", handler.Status)
	public.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	// PROTECTED
	protected := http.NewServeMux()
	protected.HandleFunc("/secret/put", handler.PutSecret)
	protected.HandleFunc("/secret/get", handler.GetSecret)
	protected.HandleFunc("/secret/list", handler.ListSecrets)
	protected.HandleFunc("/seal", handler.Seal)
	protected.HandleFunc("/token/create", handler.CreateToken)
	protected.HandleFunc("/token/revoke", handler.RevokeToken)
	protected.HandleFunc("/token/list", handler.ListTokens)

	// APPLY MIDDLEWARE
	secured := middleware.AuthMiddleware(store)(protected)

	finalMux := http.NewServeMux()
	finalMux.Handle("/unseal", public)
	finalMux.Handle("/status", public)
	finalMux.Handle("/health", public)
	finalMux.Handle("/", secured)

	log.Fatal(http.ListenAndServe(":8080", finalMux))
}
