package main

import (
	"fmt"
	"net/http"
	"webScraping/scraping"
)

func main() {

	// http.HandleFunc("GET /api/mercado-libre", scraping.GetDataMercadolibre)
	// http.HandleFunc("GET /api/puma", scraping.GetDataPuma)
	// http.HandleFunc("GET /api/adidas", scraping.GetDataAdidas)
	http.HandleFunc("GET /api/scraping", scraping.WebScraping)
	http.HandleFunc("GET /api/scrapingMercadoLibre", scraping.WebScrapingMercadoLibre)
	fmt.Println("Servidor corriendo en el puerto 8080")
	http.ListenAndServe(":8080", AddCORSHeaders(http.DefaultServeMux))

}

// func wr(n string, d string) {
// 	file, err := os.Create(n + ".html")
// 	if err != nil {
// 		fmt.Println("ewr")
// 		fmt.Println(err)
// 	}
// 	defer file.Close()
// 	file.WriteString(d)

// }

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
