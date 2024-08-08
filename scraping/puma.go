package scraping

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func GetDataPuma(w http.ResponseWriter, r *http.Request) {

	coditm := r.URL.Query().Get("search")
	// marca := r.URL.Query().Get("marca")
	categoria := r.URL.Query().Get("categoria")
	genero := r.URL.Query().Get("genero")
	// talle := r.URL.Query().Get("Talle")
	// material := r.URL.Query().Get("material")

	// search := fmt.Sprintf("%s %s %s %s %s %s", coditm, marca, categoria, material, genero, talle)

	if genero == "DAMAS" {
		genero = "Mujer"
	}
	if genero == "UNISEX" {
		genero = ""
	}

	search := fmt.Sprintf("%s %s %s", categoria, coditm, genero)

	fmt.Println(search)

	url, err := launcher.New().Headless(false).Launch()
	if err != nil {
		fmt.Println("Erorrrrrrrr")
		fmt.Println(err)

	}

	browser := rod.New().ControlURL(url).MustConnect()
	defer browser.Close()

	fmt.Println("entrando en Puma ")
	page := browser.MustPage("https://ar.puma.com/")

	// Llenar el formulario y hacer clic en el botón de búsqueda
	page.MustElement("#search-field").MustInput(search)
	page.MustElement(".SearchField-SearchFieldIcon").MustClick()
	page.MustWaitLoad()

	// aplicar filtros

	// scrapin
	listItems := srapingPuma(page)

	fmt.Println("fin")
	w.Header().Set("Content-Type", "application/json")

	// fmt.Fprintf(w, "Buscando %s", &jsonData)
	json.NewEncoder(w).Encode(listItems)

	//

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
