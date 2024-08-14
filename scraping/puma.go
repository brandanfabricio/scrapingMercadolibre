package scraping

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func GetDataPuma(w http.ResponseWriter, r *http.Request) []Items {

	proveedor := r.URL.Query().Get("proveedor")
	search := r.URL.Query().Get("search")

	var urlSearch string
	if search != "" {
		urlSearch = fmt.Sprintf("https://ar.puma.com/segmentifysearch?q=%s_*", proveedor)

	} else {

		urlSearch = fmt.Sprintf("https://ar.puma.com/segmentifysearch?q=%s_*", search)

	}

	// color := r.URL.Query().Get("color")

	// color = color[1:]

	fmt.Println(urlSearch)

	url, err := launcher.New().Headless(true).Launch()
	if err != nil {
		fmt.Println("Erorrrrrrrr")
		fmt.Println(err)

	}

	browser := rod.New().ControlURL(url).MustConnect()
	defer browser.Close()

	fmt.Println("entrando en Puma ")
	page := browser.MustPage(urlSearch)
	time.Sleep(5 * time.Second) // Esperar antes de intentar nuevamente

	page.MustWaitLoad()

	// aplicar filtros

	// scrapin
	listItems := srapingPuma(page)

	// fmt.Println("fin")
	// w.Header().Set("Content-Type", "application/json")

	// // fmt.Fprintf(w, "Buscando %s", &jsonData)
	// json.NewEncoder(w).Encode(listItems)

	// //
	fmt.Println("fin")
	return listItems

}

func srapingPuma(page *rod.Page) []Items {
	// ProductListPage
	var listItems []Items

	fmt.Println("iniciando scraping")

	containerPage, err := page.Elements(".ProductListPage")
	// fmt.Println(listProduct)

	if err != nil {
		fmt.Println("No hay datos")
		return []Items{}

	}
	if len(containerPage) <= 0 {
		return []Items{}
	}

	listProduct, err := containerPage.First().Elements("li.ProductCard")

	if err != nil {
		fmt.Println("No hay datos")
		return []Items{}
	}
	// Wr("ht/prueba",)

	for _, product := range listProduct {
		// go wr("ht/puma", product.MustHTML())
		item := Items{}

		title := product.MustElement(".ProductCard-Name").MustText()
		item.Title = title

		price := product.MustElement(".ProductPrice-CurrentPrice").MustText()
		item.Precio = price

		item.Marca = "Puma"
		item.Vendedor = "Puma"

		url := product.MustElement(".ProductCard-Link").MustAttribute("href")
		item.Url = *url

		var listLinkImage []string
		listImage := product.MustElements("img.Image-Image")
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
