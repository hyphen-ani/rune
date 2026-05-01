package main

import (
	_ "log"
	_ "net/http"
	"rune/internal/server"
)

func main() {
	server.Start()
}
