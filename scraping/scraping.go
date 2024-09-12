package scraping

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

// BrowserManager es una estructura para manejar la instancia del navegador.
type BrowserManager struct {
	browser *rod.Browser
	mu      sync.Mutex
	once    sync.Once
}

// NewBrowserManager crea una nueva instancia de BrowserManager.
func NewBrowserManager() *BrowserManager {
	return &BrowserManager{}
}

// initializeBrowser inicializa la instancia del navegador.
func (bm *BrowserManager) initializeBrowser() {
	url, err := launcher.New().Headless(true).NoSandbox(true).Launch()
	if err != nil {
		fmt.Println(err)
		return
	}
	bm.browser = rod.New().ControlURL(url).MustConnect()
}

// GetPage devuelve una nueva página utilizando la instancia compartida del navegador.
func (bm *BrowserManager) GetPage(ctx context.Context, url string) (*rod.Page, error) {
	// Inicializa el navegador solo una vez
	bm.once.Do(bm.initializeBrowser)

	// Controla el acceso concurrente a la instancia del navegador
	bm.mu.Lock()
	defer bm.mu.Unlock()

	page := bm.browser.MustPage(url)

	// Gestiona la cancelación de la página
	select {
	case <-ctx.Done(): // Si el contexto se cancela
		page.Close() // Cierra la página
		return nil, ctx.Err()
	default:
		return page, nil
	}
}

// Close cierra la instancia del navegador.
func (bm *BrowserManager) Close() {
	bm.browser.Close()
}

var bm = NewBrowserManager()

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
	data := map[string]interface{}{
		"puma":         []Items{},
		"adidas":       []Items{},
		"nike":         []Items{},
		"mercadoLibre": []Items{},
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	proveedor := r.URL.Query().Get("marca")
	fmt.Println("Proveedor:", proveedor)

	switch proveedor {
	case "PUMA":
		wg.Add(1)
		go func() {
			defer wg.Done()
			pumaItem := GetDataPuma(ctx, r)
			resultChan <- map[string]interface{}{"puma": pumaItem}
		}()
	case "ADIDAS":
		wg.Add(1)
		go func() {
			defer wg.Done()
			adidaItem := GetDataAdidas(ctx, r)
			resultChan <- map[string]interface{}{"adidas": adidaItem}
		}()
	case "NIKE":
		wg.Add(1)
		go func() {
			defer wg.Done()
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

// func WebScrapingMercadoLibre(w http.ResponseWriter, r *http.Request) {

// 	var mercadolibreItem []Items

// 	proveedor := r.URL.Query().Get("marca")

// 	switch proveedor {
// 	case "PUMA":
// 		mercadolibreItem = GetDataMercadolibrePuma(w, r)
// 	case "ADIDAS":
// 		mercadolibreItem = GetDataMercadolibreAdidas(w, r)
// 	case "NIKE":
// 		mercadolibreItem = GetDataMercadolibreNike(w, r)
// 	default:
// 		mercadolibreItem = GetDataMercadolibre(w, r)

// 	}

// 	data := map[string]interface{}{
// 		"mercadoLibre": mercadolibreItem,
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(data)
// }
// func WebScraping(w http.ResponseWriter, r *http.Request) {

// 	var wg sync.WaitGroup
// 	resultChan := make(chan map[string]interface{}, 4)

// 	// Inicializar las claves con valores vacíos
// 	data := map[string]interface{}{
// 		"puma":         []Items{},
// 		"adidas":       []Items{},
// 		"nike":         []Items{},
// 		"mercadoLibre": []Items{},
// 	}

// 	// Lanzar la rutina para siempre ejecutar el scraping de MercadoLibre
// 	// wg.Add(1)
// 	// go func() {
// 	// 	defer wg.Done()
// 	// 	mercadolibreItem := GetDataMercadolibre(w, r)
// 	// 	resultChan <- map[string]interface{}{"mercadoLibre": mercadolibreItem}
// 	// }()

// 	// Obtener el proveedor de la query string
// 	// proveedor := r.URL.Query().Get("marca")
// 	proveedor := r.URL.Query().Get("marca")
// 	fmt.Println("Proveedor:", proveedor)

// 	// Usar un switch para determinar qué otras funciones de scraping ejecutar
// 	switch proveedor {
// 	case "PUMA":
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			pumaItem := GetDataPuma(w, r)
// 			resultChan <- map[string]interface{}{"puma": pumaItem}
// 		}()
// 	case "ADIDAS":
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			adidaItem := GetDataAdidas(w, r)
// 			resultChan <- map[string]interface{}{"adidas": adidaItem}
// 		}()
// 	case "NIKE":
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			nikeItem := GetDataNike(w, r)
// 			resultChan <- map[string]interface{}{"nike": nikeItem}
// 		}()
// 	}

// 	// Cerrar el canal después de que todas las rutinas hayan terminado
// 	go func() {
// 		wg.Wait()
// 		close(resultChan)
// 	}()

// 	// Recopilar los resultados de todas las rutinas
// 	for result := range resultChan {
// 		for key, value := range result {
// 			data[key] = value
// 		}
// 	}

// 	fmt.Println("Fin scraping")
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(data)
// }
