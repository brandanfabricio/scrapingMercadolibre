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
		Headless(true).  // Ejecutar en modo no-headless para ser menos detectable
		NoSandbox(true). // Omitir la caja de arena para evitar detección

		Leakless(true). // Desactivar los argumentos que revelan el modo headless
		Devtools(true). // Permitir herramientas de desarrollador para parecer más real
		Launch()
	if err != nil {
		fmt.Println(err)
		LoggerError(err.Error())
		http.Error(w, "Error launching browser", http.StatusInternalServerError)
		return nil
	}
	browser := rod.New().ControlURL(url).MustConnect()
	defer browser.Close()
	incognitoContext := browser.MustIncognito()
	defer incognitoContext.Close()
	fmt.Println("entrando en nike ")
	fmt.Println(urlSearch)
	LoggerInfo(urlSearch)
	userAgent := &proto.NetworkSetUserAgentOverride{
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
	}
	// checkbox
	page := incognitoContext.MustPage(urlSearch)
	page.MustSetUserAgent(userAgent)
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
	listItems := scrapingNike(page, proveedor)
	if len(listItems) <= 0 {
		LoggerInfo("Utimo intento")
		page.Close()
		page = incognitoContext.MustPage(urlSearch)
		page.MustSetUserAgent(userAgent)
		page.MustWaitLoad()
		listItems = scrapingNike(page, proveedor)
	}
	fmt.Println("fin nike")
	return listItems
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
