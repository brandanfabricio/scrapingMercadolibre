package scraping

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

/*
func WebScraping(w http.ResponseWriter, r *http.Request) {

	mercadolibreItem := GetDataMercadolibre(w, r)
	// pumaItem := GetDataPuma(w, r)
	// nikeItem := GetDataNike(w, r)
	adidaItem := GetDataAdidas(w, r)
	// // hola
	pumaItem := []Items{}
	// adidaItem := []Items{}
	nikeItem := []Items{}

	data := map[string]interface{}{
		"puma":         pumaItem,
		"adidas":       adidaItem,
		"nike":         nikeItem,
		"mercadoLibre": mercadolibreItem,
	}

	fmt.Println("fin scrapin")
	w.Header().Set("Content-Type", "application/json")

	// fmt.Fprintf(w, "Buscando %s", &jsonData)
	json.NewEncoder(w).Encode(data)

}
*/
func WebScraping(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	resultChan := make(chan map[string]interface{}, 4)

	// Lanzar las rutinas para cada función de scraping
	wg.Add(1)
	go func() {
		defer wg.Done()
		// mercadolibreItem := GetDataMercadolibre(w, r)
		mercadolibreItem := []Items{}
		resultChan <- map[string]interface{}{"mercadoLibre": mercadolibreItem}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		// pumaItem := GetDataPuma(w, r)
		pumaItem := []Items{}
		resultChan <- map[string]interface{}{"puma": pumaItem}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		nikeItem := GetDataNike(w, r)
		// nikeItem := []Items{}
		resultChan <- map[string]interface{}{"nike": nikeItem}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		// adidaItem := GetDataAdidas(w, r)
		adidaItem := []Items{}
		resultChan <- map[string]interface{}{"adidas": adidaItem}
	}()

	// Cerrar el canal después de que todas las rutinas hayan terminado
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Recopilar los resultados de todas las rutinas
	data := make(map[string]interface{})
	for result := range resultChan {
		for key, value := range result {
			data[key] = value
		}
	}

	fmt.Println("fin scrapin")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
