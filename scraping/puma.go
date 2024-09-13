package scraping

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/go-rod/rod"
)

func GetDataPuma(ctx context.Context, r *http.Request) []Items {
	proveedor := r.URL.Query().Get("proveedor")
	search := r.URL.Query().Get("search")

	var urlSearch string
	if proveedor != "" {
		urlSearch = fmt.Sprintf("https://ar.puma.com/segmentifysearch?q=%s_*", proveedor)
	} else {
		urlSearch = fmt.Sprintf("https://ar.puma.com/segmentifysearch?q=%s_*", search)
	}
	fmt.Println("entrando en Puma ")
	fmt.Println(urlSearch)
	LoggerInfo(urlSearch)

	page, err := bm.GetPage(ctx, urlSearch)
	if err != nil {
		fmt.Println("Error al obtener la página:", err)
		return nil
	}
	defer page.Close()

	page.MustWaitLoad()
	time.Sleep(2 * time.Second)

	listItems := scrapingPuma(page, proveedor)

	fmt.Println("fin scraping puma")
	return listItems
}

func scrapingPuma(page *rod.Page, proveedor string) []Items {
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

		re := regexp.MustCompile(fmt.Sprintf(`%s-(\d+)`, proveedor))
		match := re.FindStringSubmatch(*url)

		if len(match) > 1 {
			codigoFinal := match[0] // 107993-01

			if codigoFinal != "" {

				item.CodProveedor = codigoFinal
			}

		}

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
