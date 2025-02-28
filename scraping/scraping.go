package scraping

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"sync"
	"time"
	"webScraping/lib"
)

var bm = lib.NewBrowserManager()

// handlePanic es una función para capturar y manejar pánicos
func handlePanic() {
	if r := recover(); r != nil {
		stringError := fmt.Sprintf("Recuperado del pánico: %v", r)
		lib.LoggerError(stringError)
	}
}
func WebScrapingMercadoLibre(w http.ResponseWriter, r *http.Request) {
	var mercadolibreItem []Items
	proveedor := r.URL.Query().Get("marca")
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()
	switch proveedor {
	case "PUMA":
		defer handlePanic()
		mercadolibreItem = GetDataMercadolibrePuma(ctx, r)
	case "ADIDAS":
		defer handlePanic()
		mercadolibreItem = GetDataMercadolibreAdidas(ctx, r)
	case "NIKE":
		defer handlePanic()
		mercadolibreItem = GetDataMercadolibreNike(ctx, r)
	default:
		defer handlePanic()
		mercadolibreItem = GetDataMercadolibre(ctx, r)
	}

	if len(mercadolibreItem) < 1 {
		mercadolibreItem = []Items{}
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

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
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
			if reflect.ValueOf(value).Len() <= 0 {
				data[key] = []Items{}
			} else {
				data[key] = value
			}
		}
	}

	fmt.Println("Fin scraping")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
