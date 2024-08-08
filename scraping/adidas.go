package scraping

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/launcher"
)

func GetDataAdidas(w http.ResponseWriter, r *http.Request) {

	coditm := r.URL.Query().Get("search")
	// marca := r.URL.Query().Get("marca")
	categoria := r.URL.Query().Get("categoria")
	// genero := r.URL.Query().Get("genero")
	// talle := r.URL.Query().Get("Talle")
	// material := r.URL.Query().Get("material")

	search := fmt.Sprintf("%s %s", categoria, coditm)

	fmt.Println(search)
	url, err := launcher.New().Headless(false).Launch()
	if err != nil {
		fmt.Println("Erorrrrrrrr")
		fmt.Println(err)

	}
	browser := rod.New().ControlURL(url).MustConnect() //
	defer browser.Close()
	fmt.Println("entrando en Adidas ")
	page := browser.MustPage("https://www.adidas.com.ar/")
	// Llenar el formulario y hacer clic en el botón de búsqueda

	page.MustElement("#glass-gdpr-default-consent-accept-button").MustClick()

	page.MustElement("._input_1f3oz_13").MustInput(search)
	time.Sleep(1 * time.Second)

	// page.MustElement("._icon_1f3oz_44").MustClick()

	err = page.Keyboard.Press(input.Enter)
	if err != nil {
		fmt.Println("Erorrrrrrrr al hacer enter")
		fmt.Println(err)

	}

	page.MustWaitLoad()
	// time.Sleep(2 * time.Second)

	// var listItems []Items
	listItems := srapingAdidas(page)

	fmt.Println("fin")
	w.Header().Set("Content-Type", "application/json")

	// fmt.Fprintf(w, "Buscando %s", &jsonData)
	json.NewEncoder(w).Encode(listItems)

}

func srapingAdidas(page *rod.Page) []Items {
	page.MustWaitLoad()
	time.Sleep(2 * time.Second)

	var listItems []Items

	fmt.Println("iniciando scraping")

	containerPage, err := page.Elements(".plp-grid___1FP1J")

	if err != nil {
		fmt.Println("No hay datos")
		return []Items{}

	}

	if len(containerPage) <= 0 {
		return []Items{}
	}

	listProduct, err := containerPage.First().Elements("div.grid-item")

	if err != nil {
		fmt.Println("No hay datos")
		return []Items{}
	}
	// Wr("ht/prueba",)

	for _, elemts := range listProduct {

		var item Items
		urlLinks := elemts.MustElement(".glass-product-card__assets-link")

		links := urlLinks.MustAttribute("href")

		// fmt.Println(*links)
		item.Url = *links

		title := elemts.MustElement(".glass-product-card__title").MustText()
		item.Title = title

		prices, err := elemts.Element(".gl-price-item")

		price := ""
		if err != nil {
			fmt.Println("Erro")
			price = ""
		} else {
			price = prices.MustText()
		}
		item.Precio = price

		item.Marca = "Adidas"
		item.Vendedor = "Adidas"

		var listLinksImg []string
		linksImg := elemts.MustElements("img.product-card-image")

		for _, img := range linksImg {

			link := img.MustAttribute("src")

			listLinksImg = append(listLinksImg, *link)

		}

		item.Imagenes = listLinksImg

		listItems = append(listItems, item)

		// fmt.Println(Items)
	}

	return listItems
}

// plp-grid___1FP1J
