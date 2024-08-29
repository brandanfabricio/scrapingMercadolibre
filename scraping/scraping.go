package scraping

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

func WebScrapingMercadoLibre(w http.ResponseWriter, r *http.Request) {

	var mercadolibreItem []Items

	proveedor := r.URL.Query().Get("marca")

	switch proveedor {
	case "PUMA":
		mercadolibreItem = GetDataMercadolibrePuma(w, r)
	case "ADIDAS":
		mercadolibreItem = GetDataMercadolibreAdidas(w, r)
	case "NIKE":
		mercadolibreItem = GetDataMercadolibreNike(w, r)
	default:
		mercadolibreItem = GetDataMercadolibre(w, r)

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

	// Inicializar las claves con valores vacíos
	data := map[string]interface{}{
		"puma":         []Items{},
		"adidas":       []Items{},
		"nike":         []Items{},
		"mercadoLibre": []Items{},
	}

	// Lanzar la rutina para siempre ejecutar el scraping de MercadoLibre
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	mercadolibreItem := GetDataMercadolibre(w, r)
	// 	resultChan <- map[string]interface{}{"mercadoLibre": mercadolibreItem}
	// }()

	// Obtener el proveedor de la query string
	// proveedor := r.URL.Query().Get("marca")
	proveedor := r.URL.Query().Get("marca")
	fmt.Println("Proveedor:", proveedor)

	// Usar un switch para determinar qué otras funciones de scraping ejecutar
	switch proveedor {
	case "PUMA":
		wg.Add(1)
		go func() {
			defer wg.Done()
			pumaItem := GetDataPuma(w, r)
			resultChan <- map[string]interface{}{"puma": pumaItem}
		}()
	case "ADIDAS":
		wg.Add(1)
		go func() {
			defer wg.Done()
			adidaItem := GetDataAdidas(w, r)
			resultChan <- map[string]interface{}{"adidas": adidaItem}
		}()
	case "NIKE":
		wg.Add(1)
		go func() {
			defer wg.Done()
			nikeItem := GetDataNike(w, r)
			resultChan <- map[string]interface{}{"nike": nikeItem}
		}()
	}

	// Cerrar el canal después de que todas las rutinas hayan terminado
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Recopilar los resultados de todas las rutinas
	for result := range resultChan {
		for key, value := range result {
			data[key] = value
		}
	}

	fmt.Println("Fin scraping")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
