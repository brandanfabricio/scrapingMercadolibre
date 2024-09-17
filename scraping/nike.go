package scraping

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/go-rod/rod"
)

func GetDataNike(ctx context.Context, r *http.Request) []Items {

	proveedor := r.URL.Query().Get("proveedor")
	search := r.URL.Query().Get("search")

	var urlSearch string

	if proveedor == "" {
		urlSearch = fmt.Sprintf("https://www.nike.com.ar/%s?_q=%s&map=ft", search, search)
	} else {
		urlSearch = fmt.Sprintf("https://www.nike.com.ar/%s?_q=%s&map=ft", proveedor, proveedor)
	}

	fmt.Println("entrando en Nike ")
	fmt.Println(urlSearch)
	LoggerInfo(urlSearch)
	page, err := bm.GetPage(ctx, urlSearch)

	if err != nil {
		fmt.Println("Error al obtener la página:", err)
		return nil
	}
	defer page.Close()
	done := make(chan bool)

	go func() {
		page.MustWaitLoad()
		// checkbox, err := page.Element(`.no-js`)
		// Verificar si se ha encontrado un CAPTCHA
		for i := 0; i < 6; i++ {
			checkbox, err := page.Elements(`.no-js`)
			if err == nil {
				if len(checkbox) > 0 {
					time.Sleep(6 * time.Second)
					fmt.Println("CAPTCHA encontrado, cerrando página y reintentando...")
					LoggerWarning("CAPTCHA encontrado, cerrando página y reintentando...")
					// Cerrar la página y reabrir una nueva instancia
					page.Close()
					page.MustWaitLoad()
				} else {
					done <- true
					break
				}
			} else {
				done <- true
				break
			}
		}
	}()
	select {
	case <-done:
		listItems := scrapingNike(page, proveedor)
		if len(listItems) <= 0 {
			LoggerInfo("Utimo intento")
			page.Close()
			page.MustWaitLoad()
			listItems = scrapingNike(page, proveedor)
		}
		fmt.Println("fin nike")
		return listItems
	case <-ctx.Done():
		fmt.Println("Timeout o contexto cancelado en Puma ", ctx.Done())
		return []Items{}
	}
}

func scrapingNike(page *rod.Page, proveedor string) []Items {
	page.MustWaitLoad()
	fmt.Println("iniciando scraping")
	page.MustWaitLoad()
	time.Sleep(2 * time.Second)
	var listItems []Items
	containerPage, err := page.Elements("#gallery-layout-container")
	if err != nil {
		fmt.Println("contenido no encontrado")
		return []Items{}
	}
	if len(containerPage) <= 0 {
		fmt.Println("No hay datos")
		return []Items{}
	}
	listProduct, err := containerPage.First().Elements("div.nikear-search-result-4-x-galleryItem")
	if err != nil {
		fmt.Println("No hay datos")
		return []Items{}
	}
	for _, product := range listProduct {
		item := Items{}
		title := product.MustElement(".vtex-product-summary-2-x-nameContainer").MustText()
		item.Title = title

		price := product.MustElement(".vtex-product-price-1-x-sellingPrice").MustText()
		item.Precio = price
		isOldPrice, err := product.Element(".vtex-product-price-1-x-listPrice")
		if err == nil {
			item.PrecioAntiguo = isOldPrice.MustText()
		}
		isPorcentaje, err := product.Element(".vtex-product-price-1-x-savingsPercentage")

		if err == nil {

			item.Porcentaje = isPorcentaje.MustText()
		}

		item.Marca = "Nike"
		item.Vendedor = "Nike"

		var listLinkImage []string
		var link string
		url := product.MustElement("a.vtex-product-summary-2-x-clearLink").MustAttribute("href")
		locateUrl := fmt.Sprintf("https://www.nike.com.ar/%s", *url)
		item.Url = locateUrl

		item.CodProveedor = proveedor
		re := regexp.MustCompile(`\b\w{6}-\w{3}\b`)
		match := re.FindString(locateUrl)

		if match != "" {
			item.CodProveedor = match
		}

		listImage := product.MustElements("img.vtex-product-summary-2-x-imageNormal")

		for _, image := range listImage {
			linkImage, err := image.Attribute("src")
			if err != nil {
				link = ""
			} else {
				link = *linkImage
			}
			listLinkImage = append(listLinkImage, link)
		}
		item.Imagenes = listLinkImage

		listItems = append(listItems, item)

	}

	return listItems
}
