package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
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

	"golang.org/x/term"
)

const port = constants.DefaultPort

const (
	ansiReset   = "\033[0m"
	ansiBold    = "\033[1m"
	ansiDim     = "\033[2m"
	ansiCyan    = "\033[96m"
	ansiBlue    = "\033[94m"
	ansiYellow  = "\033[33m"
	ansiGray    = "\033[90m"
	ansiMagenta = "\033[95m"
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

		fmt.Print("Set A Passphrase (Required For Sealing and Unsealing): ")
		passphraseBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		if err != nil {
			log.Fatal("Failed to read passphrase: ", err)
		}
		passphrase := string(passphraseBytes)

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

		sep := "  " + strings.Repeat("─", 72)
		fmt.Println()
		fmt.Println(sep)
		fmt.Printf("  %s%sroot token%s\n", ansiBold, ansiYellow, ansiReset)
		fmt.Println(sep)
		fmt.Println()
		fmt.Printf("  %s%s%s%s\n", ansiBold, ansiCyan, token, ansiReset)
		fmt.Println()
		fmt.Println(sep)
		fmt.Printf("  %sshown once  ·  not stored  ·  full vault access%s\n", ansiDim, ansiReset)
		fmt.Println(sep)
		fmt.Println()

		fmt.Println()
		fmt.Printf("  %s%s✦  You're one of the first people to run Rune — thank you!%s\n", ansiBold, ansiMagenta, ansiReset)
		fmt.Println()
		fmt.Printf("  %sRune is a passion project built for developers who take%s\n", ansiDim, ansiReset)
		fmt.Printf("  %ssecrets seriously. Your early trust means everything.%s\n", ansiDim, ansiReset)
		fmt.Println()
		fmt.Printf("  %sIf Rune is useful to you, a ★ on GitHub goes a long way%s\n", ansiDim, ansiReset)
		fmt.Printf("  %sand helps others discover the project.%s\n", ansiDim, ansiReset)
		fmt.Println()
		fmt.Printf("  %s→  https://github.com/runelock-io/rune%s\n", ansiCyan, ansiReset)
		fmt.Println()

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
	protected.HandleFunc("/token/create-system", handler.CreateSystemToken)
	protected.HandleFunc("/token/revoke", handler.RevokeToken)
	protected.HandleFunc("/token/list", handler.ListTokens)
	protected.HandleFunc("/namespace/create", handler.CreateNamespace)
	protected.HandleFunc("/namespace/list", handler.ListNamespaces)
	protected.HandleFunc("/namespace/delete", handler.DeleteNamespace)
	protected.HandleFunc("/secret/version", handler.GetSecretVersion)
	protected.HandleFunc("/secret/history", handler.ListSecretVersions)
	protected.HandleFunc("/secret/all", handler.GetAllSecrets)
	protected.HandleFunc("/stats", handler.GetStats)

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
