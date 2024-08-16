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
// hola, como estas ?
/*
func WebScraping(w http.ResponseWriter, r *http.Request) {

	// proveedor := r.URL.Query().Get("marca")
	// fmt.Println(proveedor)
	// switch proveedor {
	// case "PUMA":
	// 	fmt.Println("aki")
	// case "ADIDAS":
	// 	fmt.Println("aya")
	// }

	var wg sync.WaitGroup
	resultChan := make(chan map[string]interface{}, 4)

	// Lanzar las rutinas para cada función de scraping
	wg.Add(1)
	go func() {
		defer wg.Done()
		mercadolibreItem := GetDataMercadolibre(w, r)
		// mercadolibreItem := []Items{}
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
**/

func WebScraping(w http.ResponseWriter, r *http.Request) {

	var wg sync.WaitGroup
	resultChan := make(chan map[string]interface{}, 4)

	// Inicializar las claves con valores vacíos
	data := map[string]interface{}{
		"puma":   []Items{},
		"adidas": []Items{},
		"nike":   []Items{},
	}

	// Lanzar la rutina para siempre ejecutar el scraping de MercadoLibre
	wg.Add(1)
	go func() {
		defer wg.Done()
		mercadolibreItem := GetDataMercadolibre(w, r)
		resultChan <- map[string]interface{}{"mercadoLibre": mercadolibreItem}
	}()

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
	default:
		// En caso de que no se proporcione un proveedor válido o reconocido
		fmt.Println("Proveedor no reconocido, no se ejecutará scraping adicional.")
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
