package main

import (
	_ "KDZ/docs"
	"KDZ/internal/api-gateway"
	"fmt"
	"log"
	"net/http"
)

func main() {
	forwarder := &api_gateway.DefaultRequestForwarder{}
	handler := &api_gateway.APIHandler{Forwarder: forwarder}

	router := api_gateway.NewRouter(handler)

	port := ":8080"
	fmt.Printf("\nAPI Gateway запущен на порту %s\n", port)

	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
