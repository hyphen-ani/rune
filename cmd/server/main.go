package main

import (
	"log"
	_ "log"
	"net/http"
	_ "net/http"
	"rune/internal/api"
	"rune/internal/service"
	"rune/internal/storage"
)

func main() {
	log.Println("Starting Rune Server on: 8080")

	key := []byte("example key 1234example key 1234")

	store, err := storage.NewBoltStore("rune.db")
	if err != nil {
		log.Fatal(err)
	}

	secretService := service.NewSecretService(store, key)
	handler := api.NewHandler(secretService)

	//Routes For Rune
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	http.HandleFunc("/secret/put", handler.PutSecret)
	http.HandleFunc("/secret/get", handler.GetSecret)

	log.Fatal(http.ListenAndServe(":8080", nil))

}
