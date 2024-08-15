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
	if proveedor != "" {
		urlSearch = fmt.Sprintf("https://ar.puma.com/segmentifysearch?q=%s_*", proveedor)
	} else {
		urlSearch = fmt.Sprintf("https://ar.puma.com/segmentifysearch?q=%s_*", search)
	}
	// iniciar brouser
	url, err := launcher.New().Headless(true).Launch()
	if err != nil {
		http.Error(w, "Error launching browser", http.StatusInternalServerError)
		return nil
	}
	browser := rod.New().ControlURL(url).MustConnect()
	defer browser.Close()
	// navegando
	fmt.Println("entrando en Puma ")
	fmt.Println(urlSearch)
	page := browser.MustPage(urlSearch)

	page.MustWaitLoad()

	// iniciando scraping
	listItems := scrapingPuma(page)

	fmt.Println("fin scraping puma")
	return listItems

}

func scrapingPuma(page *rod.Page) []Items {
	time.Sleep(1 * time.Second)
	page.MustWaitLoad()
	fmt.Println("iniciando scraping")
	var listItems []Items
	// rebisando si se obtuve contendoo buscado
	containerPage, err := page.Elements(".ProductListPage")
	if err != nil {
		fmt.Println("contenido no encontrado")
		return []Items{}
	}
	if len(containerPage) <= 0 {
		fmt.Println("No hay datos")
		return []Items{}
	}
	// obtenidndo las card
	listProduct, err := containerPage.First().Elements("li.ProductCard")
	if err != nil {
		fmt.Println("No hay datos")
		return []Items{}
	}
	for _, product := range listProduct {
		item := Items{}
		// obtenedr descripcion
		title := product.MustElement(".ProductCard-Name").MustText()
		item.Title = title
		//obtener precio
		price := product.MustElement(".ProductPrice-CurrentPrice").MustText()
		item.Precio = price
		// ProductPrice-HighPrice

		existOldPrice, err := product.Element(".ProductPrice-HighPrice")
		var oldPrice string
		if err == nil {
			oldPrice = existOldPrice.MustText()
		} else {
			oldPrice = ""
		}
		item.PrecioAntiguo = oldPrice
		var porcentage string
		isPorcentage, err := product.Element(".ProductPrice-PercentageLabel")
		if err == nil {
			porcentage = isPorcentage.MustText()
		} else {
			porcentage = ""
		}
		item.Porcentaje = porcentage

		// obtener url para navegar
		url := product.MustElement(".ProductCard-Link").MustAttribute("href")
		item.Url = *url
		// obtener imagenes
		var listLinkImage []string
		listImage := product.MustElements("img.Image-Image")

		var link string
		linkImage, err := listImage.First().Attribute("src")

		if err != nil {
			link = ""
		}
		link = *linkImage
		listLinkImage = append(listLinkImage, link)

		item.Marca = "Puma"
		item.Vendedor = "Puma"
		item.Imagenes = listLinkImage
		listItems = append(listItems, item)
	}
	return listItems
}

// func appliFilter() {}

/*

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
		http.Error(w, "Error launching browser", http.StatusInternalServerError)
		return nil
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


*/
