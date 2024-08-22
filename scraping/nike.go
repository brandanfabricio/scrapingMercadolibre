package scraping

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

func GetDataNike(w http.ResponseWriter, r *http.Request) []Items {

	proveedor := r.URL.Query().Get("proveedor")
	search := r.URL.Query().Get("search")

	var urlSearch string

	if proveedor == "" {
		urlSearch = fmt.Sprintf("https://www.nike.com.ar/%s?_q=%s&map=ft", search, search)
	} else {
		urlSearch = fmt.Sprintf("https://www.nike.com.ar/%s?_q=%s&map=ft", proveedor, proveedor)
	}
	url, err := launcher.New().
		Headless(false). // Ejecutar en modo no-headless para ser menos detectable
		// NoSandbox(true). // Omitir la caja de arena para evitar detección
		// Leakless(false). // Desactivar los argumentos que revelan el modo headless
		Devtools(true). // Permitir herramientas de desarrollador para parecer más real
		Launch()
	if err != nil {
		http.Error(w, "Error launching browser", http.StatusInternalServerError)
		return nil
	}
	browser := rod.New().ControlURL(url).MustConnect()
	defer browser.Close()
	incognitoContext := browser.MustIncognito()
	defer incognitoContext.Close()
	fmt.Println("entrando en nike ")
	fmt.Println(urlSearch)
	userAgent := &proto.NetworkSetUserAgentOverride{
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
	}
	// checkbox
	page := incognitoContext.MustPage(urlSearch)

	page.MustSetUserAgent(userAgent)
	page.MustWaitLoad()

	// checkbox, err := page.Element(`.no-js`)
	// Verificar si se ha encontrado un CAPTCHA

	for i := 0; i < 5; i++ {
		checkbox, err := page.Elements(`.no-js`)
		if err == nil {
			if len(checkbox) > 0 {

				time.Sleep(5 * time.Second)
				fmt.Println("CAPTCHA encontrado, cerrando página y reintentando...")

				// Cerrar la página y reabrir una nueva instancia
				page.Close()
				page = incognitoContext.MustPage(urlSearch)
				page.MustSetUserAgent(userAgent)
				page.MustWaitLoad()
			} else {
				break
			}
		} else {
			break
		}

	}

	// time.Sleep(10 * time.Second)
	// fmt.Println(checkbox)
	// Intentar encontrar el checkbox

	// if err == nil {

	// 	ckeck, err := checkbox.Element("#RlquG0")
	// 	fmt.Println("akii")
	// 	fmt.Println(ckeck)
	// 	if err == nil {

	// 		ckeck2 := ckeck.MustElement(`input`)
	// 		fmt.Println(ckeck2)
	// 		ckeck2.MustHover()
	// 		ckeck2.MustFocus()
	// 		ckeck2.Click(proto.InputMouseButtonRight, 1)
	// 		page.MustWaitLoad()
	// 	}

	// 	// content
	// 	// fmt.Println("ckeck book")
	// 	// checkbox := checkbox.MustElement(`input[type="checkbox"]`)
	// 	// checkbox.MustHover()
	// 	// checkbox.MustFocus()
	// 	// checkbox.Click(proto.InputMouseButtonRight, 1)
	// }

	// time.Sleep(20 * time.Second) // Esperar antes de intentar nuevamente

	listItems := scrapingNike(page, proveedor)
	fmt.Println("fin nike")
	return listItems
}

func scrapingNike(page *rod.Page, proveedor string) []Items {
	page.MustWaitLoad()
	fmt.Println("iniciando scraping")
	page.MustWaitLoad()
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
