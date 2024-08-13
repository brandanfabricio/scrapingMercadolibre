package scraping

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

func GetDataNike(w http.ResponseWriter, r *http.Request) []Items {

	proveedor := r.URL.Query().Get("proveedor")

	urlSearch := fmt.Sprintf("https://www.nike.com.ar/%s?_q=%s&map=ft", proveedor, proveedor)

	fmt.Println(urlSearch)

	// url, err := launcher.New().Headless(true).Launch()
	// if err != nil {
	// 	fmt.Println("Erorrrrrrrr")
	// 	fmt.Println(err)

	// }
	url := launcher.New().
		Headless(true).  // Ejecutar en modo no-headless para ser menos detectable
		NoSandbox(true). // Omitir la caja de arena para evitar detección
		Leakless(false). // Desactivar los argumentos que revelan el modo headless
		Devtools(true).  // Permitir herramientas de desarrollador para parecer más real

		MustLaunch()

	browser := rod.New().ControlURL(url).MustConnect()
	defer browser.Close()

	incognitoContext := browser.MustIncognito()
	defer incognitoContext.Close()

	fmt.Println("entrando en nike ")

	userAgent := &proto.NetworkSetUserAgentOverride{
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
	}
	// checkbox
	page := incognitoContext.MustPage(urlSearch)

	page.MustSetUserAgent(userAgent)

	// Ejecutar script para eliminar propiedades detectables

	// userAgent := &proto.NetworkSetUserAgentOverride{
	// 	UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
	// }
	// page.MustSetUserAgent(userAgent)

	time.Sleep(10 * time.Second) // Esperar antes de intentar nuevamente

	page.MustWaitLoad()

	// aplicar filtros

	// scrapin
	listItems := srapingNike(page)

	// fmt.Println("fin")
	// w.Header().Set("Content-Type", "application/json")

	// // fmt.Fprintf(w, "Buscando %s", &jsonData)
	// json.NewEncoder(w).Encode(listItems)

	// //
	fmt.Println("fin")
	return listItems

}

func srapingNike(page *rod.Page) []Items {
	// ProductListPage
	var listItems []Items

	fmt.Println("iniciando scraping")

	containerPage, err := page.Elements("#gallery-layout-container")
	// fmt.Println(listProduct)

	if err != nil {
		fmt.Println("No hay datos")
		return []Items{}

	}
	if len(containerPage) <= 0 {
		return []Items{}
	}

	listProduct, err := containerPage.First().Elements("div.nikear-search-result-4-x-galleryItem")

	if err != nil {
		fmt.Println("No hay datos")
		return []Items{}
	}
	// Wr("ht/prueba",)

	for _, product := range listProduct {
		// go wr("ht/puma", product.MustHTML())
		item := Items{}

		title := product.MustElement(".vtex-product-summary-2-x-nameContainer").MustText()
		item.Title = title

		price := product.MustElement(".vtex-product-price-1-x-sellingPrice").MustText()
		item.Precio = price

		item.Marca = "Nike"
		item.Vendedor = "Nike"

		url := product.MustElement("a.vtex-product-summary-2-x-clearLink").MustAttribute("href")
		locateUrl := fmt.Sprintf("https://www.nike.com.ar/%s", *url)
		item.Url = locateUrl

		var listLinkImage []string
		listImage := product.MustElements("img.vtex-product-summary-2-x-imageNormal")
		// fmt.Println(listImage)

		for _, image := range listImage {
			linkImage := image.MustAttribute("src")
			listLinkImage = append(listLinkImage, *linkImage)
		}

		item.Imagenes = listLinkImage

		listItems = append(listItems, item)

	}

	return listItems
}

// func appliFilter() {}
