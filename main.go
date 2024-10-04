package main

import (
	"fmt"
	"log"
	"net/http"
	"webScraping/lib"
	"webScraping/scraping"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("######################################")
			err := fmt.Sprintf("Recuperado del pánico General: %v", r)
			lib.LoggerError(err)
			log.Printf("######################################")
		}
	}()

	http.HandleFunc("GET /api/scraping", scraping.WebScraping)
	http.HandleFunc("GET /api/scrapingMercadoLibre", scraping.WebScrapingMercadoLibre)
	fmt.Println("Servidor corriendo en el puerto 3000")
	http.ListenAndServe(":3000", AddCORSHeaders(http.DefaultServeMux))

}

func AddCORSHeaders(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		handler.ServeHTTP(w, r)
	})
}
