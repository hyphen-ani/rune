package main

import (
	"fmt"
	"log"
	_ "log"
	"net/http"
	_ "net/http"
	"rune/internal/api"
	"rune/internal/auth"
	"rune/internal/crypto"
	"rune/internal/middleware"
	"rune/internal/seal"
	"rune/internal/service"
	"rune/internal/storage"
)

func main() {
	log.Println("Starting Rune Server on: 8080")

	store, err := storage.NewBoltStore("rune.db")
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
		err = store.SaveToken(hash)
		if err != nil {
			return
		}

		fmt.Println("====================================")
		fmt.Println("Root Token (SAVE THIS):", token)
		fmt.Println("====================================")

	}

	sealer := seal.NewManger()
	secretService := service.NewSecretService(store, sealer)
	handler := api.NewHandler(secretService, store, sealer)

	// PUBLIC
	public := http.NewServeMux()
	public.HandleFunc("/unseal", handler.Unseal)
	public.HandleFunc("/status", handler.Status)

	// PROTECTED
	protected := http.NewServeMux()
	protected.HandleFunc("/secret/put", handler.PutSecret)
	protected.HandleFunc("/secret/get", handler.GetSecret)
	protected.HandleFunc("/seal", handler.Seal)

	// APPLY MIDDLEWARE
	secured := middleware.AuthMiddleware(store)(protected)

	finalMux := http.NewServeMux()
	finalMux.Handle("/unseal", public)
	finalMux.Handle("/status", public)
	finalMux.Handle("/health", public)
	finalMux.Handle("/", secured)

	// HEALTH ROUTE
	public.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	log.Fatal(http.ListenAndServe(":8080", finalMux))
}
