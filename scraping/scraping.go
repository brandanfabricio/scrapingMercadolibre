package scraping

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
	"webScraping/scraping/lib"
)

var bm = lib.NewBrowserManager()

// handlePanic es una función para capturar y manejar pánicos
func handlePanic() {
	if r := recover(); r != nil {
		log.Printf("Recuperado del pánico: %v", r)
	}
}
func WebScrapingMercadoLibre(w http.ResponseWriter, r *http.Request) {
	var mercadolibreItem []Items
	proveedor := r.URL.Query().Get("marca")
	ctx, cancel := context.WithTimeout(r.Context(), 70*time.Second)
	defer cancel()
	switch proveedor {
	case "PUMA":
		mercadolibreItem = GetDataMercadolibrePuma(ctx, r)
	case "ADIDAS":
		mercadolibreItem = GetDataMercadolibreAdidas(ctx, r)
	case "NIKE":
		mercadolibreItem = GetDataMercadolibreNike(ctx, r)
	default:
		mercadolibreItem = GetDataMercadolibre(ctx, r)
	}

	data := map[string]interface{}{
		"mercadoLibre": mercadolibreItem,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func WebScraping(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	resultChan := make(chan map[string]interface{}, 4)
	data := map[string]interface{}{
		"puma":         []Items{},
		"adidas":       []Items{},
		"nike":         []Items{},
		"mercadoLibre": []Items{},
	}

	ctx, cancel := context.WithTimeout(r.Context(), 70*time.Second)
	defer cancel()
	// defer bm.KillChromeProcesses()

	proveedor := r.URL.Query().Get("marca")

	switch proveedor {
	case "PUMA":
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer handlePanic()
			pumaItem := GetDataPuma(ctx, r)
			resultChan <- map[string]interface{}{"puma": pumaItem}
		}()
	case "ADIDAS":
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer handlePanic()
			adidaItem := GetDataAdidas(ctx, r)
			resultChan <- map[string]interface{}{"adidas": adidaItem}
		}()
	case "NIKE":
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer handlePanic()
			nikeItem := GetDataNike(ctx, r)
			resultChan <- map[string]interface{}{"nike": nikeItem}
		}()
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for result := range resultChan {
		for key, value := range result {
			data[key] = value
		}
	}

	fmt.Println("Fin scraping")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
