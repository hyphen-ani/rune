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

const (
	ansiReset  = "\033[0m"
	ansiBold   = "\033[1m"
	ansiDim    = "\033[2m"
	ansiCyan   = "\033[96m"
	ansiBlue   = "\033[94m"
	ansiYellow = "\033[33m"
	ansiGray   = "\033[90m"
)

func printBanner(dbPath string) {
	fmt.Println()
	fmt.Printf("  %s%srune%s  %slocal-first secrets vault%s\n", ansiBold, ansiCyan, ansiReset, ansiDim, ansiReset)
	fmt.Println()
	fmt.Printf("  %s➜%s  %-6s %shttp://localhost:%s%s\n",        ansiCyan, ansiReset, "api", ansiBlue, port, ansiReset)
	fmt.Printf("  %s➜%s  %-6s %shttp://localhost:%s/ui/%s\n",    ansiCyan, ansiReset, "ui",  ansiBold+ansiBlue, port, ansiReset)
	fmt.Printf("  %s➜%s  %-6s %s%s%s\n",                         ansiCyan, ansiReset, "db",  ansiGray, dbPath, ansiReset)
	fmt.Println()
	fmt.Printf("  %s⬡  vault sealed%s  —  run %srune unseal%s or open the web UI\n", ansiYellow, ansiReset, ansiBold, ansiReset)
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

	printBanner(dbPath)
	log.Fatal(http.ListenAndServe(":"+port, finalMux))
}
