package scraping

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func WebScraping(w http.ResponseWriter, r *http.Request) {

	mercadolibreItem := GetDataMercadolibre(w, r)
	pumaItem := GetDataPuma(w, r)
	nikeItem := GetDataNike(w, r)
	adidaItem := GetDataAdidas(w, r)

	// pumaItem := []Items{}
	// adidaItem := []Items{}
	// mercadolibreItem := []Items{}

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

// func WebScraping(w http.ResponseWriter, r *http.Request) {
// 	url, err := launcher.New().Headless(false).Launch()
// 	if err != nil {
// 		fmt.Println("Error al lanzar el navegador:", err)
// 		return
// 	}

// 	browser := rod.New().ControlURL(url).MustConnect()
// 	defer browser.Close()

// 	var wg sync.WaitGroup
// 	resultChan := make(chan []Items, 3)

// 	// Crear y lanzar goroutines para cada página en una pestaña diferente
// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		tab := browser.MustPage("") // Crear una nueva pestaña vacía
// 		resultChan <- GetDataMercadolibre(w, r, tab)
// 		tab.Close() // Cerrar la pestaña después del scraping
// 	}()

// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		tab := browser.MustPage("") // Crear una nueva pestaña vacía
// 		resultChan <- GetDataPuma(w, r, tab)
// 		tab.Close() // Cerrar la pestaña después del scraping
// 	}()

// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		tab := browser.MustPage("") // Crear una nueva pestaña vacía
// 		resultChan <- GetDataAdidas(w, r, tab)
// 		tab.Close() // Cerrar la pestaña después del scraping
// 	}()

// 	// Cerrar el canal cuando todas las goroutines terminen
// 	go func() {
// 		wg.Wait()
// 		close(resultChan)
// 	}()

// 	// Almacenar los resultados en un mapa
// 	// data := make(map[string]ItemData)
// 	// for result := range resultChan {
// 	// 	data[result.Source] = result.Data
// 	// }

// 	fmt.Println("fin scraping")

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(map[string]it)
// }
