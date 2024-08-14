package scraping

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/launcher"
)

func GetDataAdidas(w http.ResponseWriter, r *http.Request) []Items {

	proveedor := r.URL.Query().Get("proveedor")
	search := r.URL.Query().Get("search")
	var urlSearch string
	if proveedor != "" {

		urlSearch = fmt.Sprintf("https://www.adidas.com.ar/search?q=%s", search)
	} else {
		urlSearch = fmt.Sprintf("https://www.adidas.com.ar/search?q=%s", proveedor)

	}

	fmt.Println(urlSearch)

	url, err := launcher.New().Headless(true).Launch()
	if err != nil {
		fmt.Println("Erorrrrrrrr")
		fmt.Println(err)
	}
	browser := rod.New().ControlURL(url).MustConnect() //
	defer browser.Close()

	fmt.Println("entrando en Adidas ")

	page := browser.MustPage(urlSearch)

	// page.MustElement("#glass-gdpr-default-consent-accept-button").MustClick()

	// page.MustElement("._icon_1f3oz_44").MustClick()

	err = page.Keyboard.Press(input.Enter)
	if err != nil {
		fmt.Println("Erorrrrrrrr al hacer enter")
		fmt.Println(err)

	}

	page.MustWaitLoad()
	// time.Sleep(2 * time.Second)

	// aplicar filtro

	// var listItems []Items

	listItems := srapingAdidas(page)
	fmt.Println("fin adidas")

	return listItems

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
		containerPage, err := page.Elements(".content-wrapper___3TFwT")
		var item Items
		if err != nil {
			fmt.Println("No hay datos")
			return []Items{}
		}

		if len(containerPage) <= 0 {
			return []Items{}

		} else {
			container := containerPage[0]
			title := container.MustElement("h1.name___120FN").MustText()
			item.Title = title
			price := container.MustElement("div.product-price___2Mip5").MustText()
			item.Precio = price[1:]
			item.Url = page.MustInfo().URL
			item.Marca = "Adidas"
			item.Vendedor = "Adidas"

			var ListLinksImg []string
			constainerImages := container.MustElement(".image-grid___1JN2z")
			linksImgs := constainerImages.MustElements("img")

			for i, img := range linksImgs {
				if i+1 >= 4 {
					break
				}
				link := img.MustAttribute("src")
				ListLinksImg = append(ListLinksImg, *link)
			}

			item.Imagenes = ListLinksImg
			listItems = append(listItems, item)
		}

	} else {
		listProduct, err := containerPage.First().Elements("div.grid-item")
		if err != nil {
			fmt.Println("No hay datos")
			return []Items{}
		}
		// Wr("ht/prueba",)
		for _, elemts := range listProduct {

			fmt.Println(elemts)

			verifix := elemts.MustAttribute("data-index")

			if *verifix == "-1" {
				fmt.Println(*verifix)
				continue
			}

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
	}

	return listItems
}

// plp-grid___1FP1J
