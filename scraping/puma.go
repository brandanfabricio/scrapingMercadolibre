package scraping

import (
	"fmt"
	"log"
	"net/http"

	"github.com/playwright-community/playwright-go"
)

func GetDataPuma(w http.ResponseWriter, r *http.Request) []Items {
	proveedor := r.URL.Query().Get("proveedor")
	search := r.URL.Query().Get("search")

	var urlSearch string
	if proveedor != "" {
		urlSearch = fmt.Sprintf("https://ar.puma.com/segmentifysearch?q=%s_*", proveedor)
	} else {
		urlSearch = fmt.Sprintf("https://ar.puma.com/segmentifysearch?q=%s_*", search)
	}

	// Instala los navegadores necesarios
	if err := playwright.Install(); err != nil {
		log.Fatalf("Error al instalar Playwright: %v", err)
	}
	// Inicia Playwright
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("No se pudo iniciar Playwright: %v", err)
	}
	// Inicia el navegador en modo visible (no headless)
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true), // Cambia a false para ver el navegador
		Devtools: playwright.Bool(true), // Abre las herramientas de desarrollo (DevTools)
	})
	if err != nil {
		log.Fatalf("No se pudo lanzar el navegador: %v", err)
	}

	// Crea una nueva página
	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("No se pudo crear la página: %v", err)
	}
	// Navega a una SPA de ejemplo
	fmt.Println("Entrado en ", urlSearch)
	if _, err = page.Goto(urlSearch, playwright.PageGotoOptions{}); err != nil {
		log.Fatalf("No se pudo navegar a la SPA: %v", err)
	}

	page.WaitForTimeout(1000)
	fmt.Println("iniciando scraping")
	Litm, err := page.Locator(".ProductCard").All()

	if err != nil {
		log.Fatalf("Could not get the product node: %v", err)
	}
	// li.ProductCard"
	fmt.Println(Litm)
	var listItems []Items
	// var listLinkImage []string
	for _, product := range Litm {
		item := Items{}
		title, err := product.Locator(".ProductCard-Name").TextContent()
		fmt.Println(title)
		if err != nil {
			log.Fatalf("Could not get the product node: %v", err)
		}
		item.Title = title
		price, err := product.Locator(".ProductPrice-CurrentPrice").TextContent()
		if err != nil {
			log.Fatalf("Could not get the product node: %v", err)
		}
		item.Precio = price
		url, err := product.Locator(".ProductCard-Link").GetAttribute("href")
		if err != nil {
// 				log.Fatalf("Could not get the product node: %v", err)
		// }
		// 		item.Url = url

		// existOldPrice, err := product.Locator(".ProductPrice-HighPrice").TextContent()

		// if err == nil {
		// 	item.PrecioAntiguo = existOldPrice
		// }

		// isPorcentage, err := product.Locator(".ProductPrice-PercentageLabel").TextContent()

		// if err == nil {
		// 	item.Porcentaje = isPorcentage
		// }

		// // obtener imagenes
		// var listLinkImage []string
		// listImage, err := product.Locator("img.Image-Image").First().GetAttribute("src")
		// if err == nil {
		// 	fmt.Println("sin img")
		// }
		// listLinkImage = append(listLinkImage, listImage)
		// item.Imagenes = listLinkImage

		item.CodProveedor = proveedor
		item.Marca = "Puma"
		listItems = append(listItems, item)

	}

	// Cierra el navegador
	if err = browser.Close(); err != nil {
		log.Fatalf("No se pudo cerrar el navegador: %v", err)
	}
	if err = pw.Stop(); err != nil {
		log.Fatalf("No se pudo detener Playwright: %v", err)
	}

	return listItems
}
