package main

import (
	"log"
	_ "log"
	"net/http"
	_ "net/http"
	"rune/internal/api"
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

	//salt, err := store.GetSalt()
	//if err != nil {
	//	salt, _ = crypto.GenerateSalt()
	//	store.SaveSalt(salt)
	//}
	//
	//var passphrase string
	//fmt.Print("Enter Passphrase: ")
	//fmt.Scanln(&passphrase)

	//key := crypto.DeriveKey(passphrase, salt)

	sealer := seal.NewManger()
	secretService := service.NewSecretService(store, sealer)
	handler := api.NewHandler(secretService, store, sealer)

	//Routes For Rune
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	http.HandleFunc("/secret/put", handler.PutSecret)
	http.HandleFunc("/secret/get", handler.GetSecret)
	http.HandleFunc("/seal", handler.Seal)
	http.HandleFunc("/unseal", handler.Unseal)
	http.HandleFunc("/status", handler.Status)

	log.Fatal(http.ListenAndServe(":8080", nil))

}
