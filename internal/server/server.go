package server

import (
	"fmt"
	"log"
	"net/http"
	"rune/internal/api"
	"rune/internal/auth"
	"rune/internal/config"
	"rune/internal/constants"
	"rune/internal/crypto"
	"rune/internal/middleware"
	"rune/internal/seal"
	"rune/internal/service"
	"rune/internal/storage"
	"rune/internal/ui"
	"time"
)

const port = "8080"

func printBanner() {
	fmt.Println()
	fmt.Println("  ┌─────────────────────────────────────────────┐")
	fmt.Println("  │                  rune vault                  │")
	fmt.Println("  ├─────────────────────────────────────────────┤")
	fmt.Printf("  │  api  →  http://localhost:%s              │\n", port)
	fmt.Printf("  │  ui   →  http://localhost:%s/ui/          │\n", port)
	fmt.Println("  └─────────────────────────────────────────────┘")
	fmt.Println()
}

func Start() {
	log.SetFlags(0)

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
		fmt.Print("Set A Passphrase (Required For Sealing and Unsealing): ")
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
			Namespace: "*",
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

	store.CreateNamespace(constants.DefaultNamespace)

	sealer := seal.NewManger()
	secretService := service.NewSecretService(store, sealer)
	tokenService := service.NewTokenService(store, sealer)
	namespaceService := service.NewNamespaceService(store)
	handler := api.NewHandler(secretService, tokenService, namespaceService, store, sealer)

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
	protected.HandleFunc("/secret/delete", handler.DeleteSecret)
	protected.HandleFunc("/secret/list", handler.ListSecrets)
	protected.HandleFunc("/secret/rotate", handler.RotateSecret)
	protected.HandleFunc("/seal", handler.Seal)
	protected.HandleFunc("/token/me", handler.GetCurrentToken)
	protected.HandleFunc("/token/create", handler.CreateToken)
	protected.HandleFunc("/token/revoke", handler.RevokeToken)
	protected.HandleFunc("/token/list", handler.ListTokens)
	protected.HandleFunc("/namespace/create", handler.CreateNamespace)
	protected.HandleFunc("/namespace/list", handler.ListNamespaces)
	protected.HandleFunc("/namespace/delete", handler.DeleteNamespace)
	protected.HandleFunc("/secret/version", handler.GetSecretVersion)
	protected.HandleFunc("/secret/history", handler.ListSecretVersions)

	// APPLY MIDDLEWARE
	secured := middleware.AuthMiddleware(store)(protected)

	finalMux := http.NewServeMux()
	finalMux.Handle("/unseal", public)
	finalMux.Handle("/status", public)
	finalMux.Handle("/health", public)
	finalMux.Handle("/ui/", http.StripPrefix("/ui", ui.Handler()))
	finalMux.Handle("/", secured)

	printBanner()
	log.Fatal(http.ListenAndServe(":"+port, finalMux))
}
