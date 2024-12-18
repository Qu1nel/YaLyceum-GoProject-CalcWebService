package main

import (
	"CalcService/configs"
	"CalcService/internal/handlers"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc(configs.Pattern, handlers.CalculateHandler)

	if err := http.ListenAndServe(configs.Address, nil); err != nil {
		log.Fatal(err)
	}
}
