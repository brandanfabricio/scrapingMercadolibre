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
	time.Sleep(3 * time.Second)

	// iniciando scraping
	listItems := scrapingPuma(page)

	fmt.Println("fin scraping puma")
	return listItems

}

func scrapingPuma(page *rod.Page) []Items {
	page.MustWaitLoad()
	fmt.Println("iniciando scraping")
	var listItems []Items
	// rebisando si se obtuve contendoo buscado
	// containerPage, err := page.Elements(".ProductListPage")
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
