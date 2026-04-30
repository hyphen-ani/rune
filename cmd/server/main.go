package main

import (
	"log"
	_ "log"
	"net/http"
	_ "net/http"
	"rune/internal/api"
	"rune/internal/service"
)

func main() {
	log.Println("Starting Rune Server on: 8080")

	secretService := service.NewSecretService()
	handler := api.NewHandler(secretService)

	//Routes For Rune
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	http.HandleFunc("/secret/put", handler.PutSecret)
	http.HandleFunc("/secret/get", handler.GetSecret)

	log.Fatal(http.ListenAndServe(":8080", nil))

}
