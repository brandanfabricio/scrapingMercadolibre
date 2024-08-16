package scraping

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func GetDataAdidas(w http.ResponseWriter, r *http.Request) []Items {
	proveedor := r.URL.Query().Get("proveedor")
	search := r.URL.Query().Get("search")
	var urlSearch string
	if proveedor == "" {
		urlSearch = fmt.Sprintf("https://www.adidas.com.ar/search?q=%s", search)
	} else {
		urlSearch = fmt.Sprintf("https://www.adidas.com.ar/search?q=%s", proveedor)
	}
	url, err := launcher.New().Headless(true).Launch()
	if err != nil {
		http.Error(w, "Error launching browser", http.StatusInternalServerError)
		return nil
	}
	browser := rod.New().ControlURL(url).MustConnect() //
	defer browser.Close()
	fmt.Println("entrando en Adidas ")
	fmt.Println(urlSearch)

	page := browser.MustPage(urlSearch)
	page.MustWaitLoad()
	time.Sleep(1 * time.Second)

	listItems := scrapingAdidas(page, proveedor)
	fmt.Println("fin adidas")
	return listItems

	// page.MustElement("#glass-gdpr-default-consent-accept-button").MustClick()

	// page.MustElement("._icon_1f3oz_44").MustClick()

	// err = page.Keyboard.Press(input.Enter)
	// if err != nil {
	// 	fmt.Println("Erorrrrrrrr al hacer enter")
	// 	fmt.Println(err)

	// }

	// time.Sleep(2 * time.Second)

	// aplicar filtro

	// var listItems []Items

}
func scrapingAdidas(page *rod.Page, proveedor string) []Items {
	page.MustWaitLoad()
	var listItems []Items
	fmt.Println("iniciando scraping")

	if proveedor == "" {

		listItems = scrapingList(page)

	} else {
		listItems = scrapingPage(page)
	}
	return listItems

}

func scrapingList(page *rod.Page) []Items {
	var listItems []Items
	page.MustWaitLoad()

	containerPage, err := page.Elements(".plp-grid___1FP1J")
	if err != nil {
		fmt.Println("contenido no encontrado")
		return []Items{}
	}
	if len(containerPage) <= 0 {
		fmt.Println("No hay datos")
		return []Items{}
	}

	listProduct, err := containerPage.First().Elements("div.grid-item")
	if err != nil {
		fmt.Println("contenido no encontrado")
		return []Items{}
	}
	for _, elemts := range listProduct {
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
			price = ""
		} else {
			price = prices.MustText()
		}
		item.Precio = price
		item.Marca = "Adidas"
		item.Vendedor = "Adidas"
		var listLinksImg []string
		var link string
		linksImg := elemts.MustElements("img.product-card-image")

		for _, img := range linksImg {
			linkImg, err := img.Attribute("src")
			if err != nil {
				link = ""
			} else {
				link = *linkImg
			}
			listLinksImg = append(listLinksImg, link)
		}
		item.Imagenes = listLinksImg
		listItems = append(listItems, item)
		// fmt.Println(Items)
	}

	return listItems

}

func scrapingPage(page *rod.Page) []Items {

	page.MustWaitLoad()
	var listItems []Items

	containerPage, err := page.Elements(".content-wrapper___3TFwT")
	var item Items
	if err != nil {
		fmt.Println("iniciando scraping")
		return []Items{}
	}
	if len(containerPage) <= 0 {
		fmt.Println("No hay datos")
		return []Items{}
	} else {
		container := containerPage.First()

		sidebar, err := container.Element(".product-description___1TLpA")
		if err != nil {
			fmt.Println("ERrrrrrrrrrrro")
		}

		title := sidebar.MustElement("h1.name___120FN").MustText()
		item.Title = title

		isExistCrooseed, err := sidebar.Element(".gl-price-item--crossed")

		var currentPrice string
		var oldPrice string

		if err != nil {
			currentPrice = container.MustElement("div.product-price___2Mip5").MustText()
			item.Precio = currentPrice[1:]

		} else {

			oldPrice = isExistCrooseed.MustText()
			item.PrecioAntiguo = oldPrice[3:]

			salePrice, err := sidebar.Element(".gl-price-item--sale")
			if err != nil {

				item.Precio = container.MustElement("div.product-price___2Mip5").MustText()

			} else {

				currentPrice := salePrice.MustText()

				item.Precio = currentPrice[3:]

				nuevo, err := strconv.ParseFloat(strings.Replace(item.Precio, ".", "", -1), 64)
				if err != nil {
					fmt.Println("err")
				}
				anti, err := strconv.ParseFloat(strings.Replace(item.PrecioAntiguo, ".", "", -1), 64)
				if err != nil {
					fmt.Println("err")
				}

				// num1, err := strconv.ParseInt(item.Precio, 10, 64)
				// if err != nil {
				// 	fmt.Println("error de conversion")
				// }
				// num2, err := strconv.ParseInt(item.PrecioAntiguo, 10, 64)

				// if err != nil {
				// 	fmt.Println("error de conversion")
				// }
				porcerntaje := int(((anti - nuevo) / anti) * 100)
				item.Porcentaje = fmt.Sprintf("%d%%", porcerntaje)

			}

		}
		// gl-price-item--sale

		//

		item.Url = page.MustInfo().URL
		item.Marca = "Adidas"
		item.Vendedor = "Adidas"
		var wg sync.WaitGroup
		var ListLinksImg []string
		constainerImages := container.MustElement(".image-grid___1JN2z")
		linksImgs := constainerImages.MustElements("img")
		wg.Add(4)
		var link string
		for i, img := range linksImgs {
			go func(img *rod.Element) {
				defer wg.Done()
				// if i+1 >= 4 {
				// 	break
				// }
				linkImg, err := img.Attribute("src")
				if err != nil {
					link = ""
				} else {
					link = *linkImg
				}

				ListLinksImg = append(ListLinksImg, link)
			}(img)
			if i+1 >= 4 {
				break
			}
		}
		wg.Wait()
		item.Imagenes = ListLinksImg
		listItems = append(listItems, item)
	}
	return listItems

}

// plp-grid___1FP1J
/*

func scrapingAdidas(page *rod.Page) []Items {
	page.MustWaitLoad()
	var listItems []Items
	fmt.Println("iniciando scraping")
	containerPage, err := page.Elements(".plp-grid___1FP1J")
	if err != nil {
		fmt.Println("contenido no encontrado")
		return []Items{}
	}
	if len(containerPage) <= 0 {
		page.MustWaitLoad()

		containerPage, err := page.Elements(".content-wrapper___3TFwT")
		var item Items
		if err != nil {
			fmt.Println("iniciando scraping")
			return []Items{}
		}
		if len(containerPage) <= 0 {
			fmt.Println("No hay datos")
			return []Items{}
		} else {
			container := containerPage.First()
			title := container.MustElement("h1.name___120FN").MustText()
			item.Title = title
			price := container.MustElement("div.product-price___2Mip5").MustText()
			item.Precio = price[1:]
			item.Url = page.MustInfo().URL
			item.Marca = "Adidas"
			item.Vendedor = "Adidas"
			var wg sync.WaitGroup
			var ListLinksImg []string
			constainerImages := container.MustElement(".image-grid___1JN2z")
			linksImgs := constainerImages.MustElements("img")
			wg.Add(4)
			var link string
			for i, img := range linksImgs {
				go func(img *rod.Element) {
					defer wg.Done()
					// if i+1 >= 4 {
					// 	break
					// }
					linkImg, err := img.Attribute("src")
					if err != nil {
						link = ""
					} else {
						link = *linkImg
					}

					ListLinksImg = append(ListLinksImg, link)
				}(img)
				if i+1 >= 4 {
					break
				}
			}
			wg.Wait()
			item.Imagenes = ListLinksImg
			listItems = append(listItems, item)
		}

	} else {
		page.MustWaitLoad()

		listProduct, err := containerPage.First().Elements("div.grid-item")
		if err != nil {
			fmt.Println("contenido no encontrado")
			return []Items{}
		}
		for _, elemts := range listProduct {
			verifix := elemts.MustAttribute("data-index")
			if *verifix == "-1" {
				fmt.Println(*verifix)
				continue
			}
			var item Items
			fmt.Println("aju")
			urlLinks := elemts.MustElement(".glass-product-card__assets-link")
			fmt.Println("ay")

			links := urlLinks.MustAttribute("href")
			// fmt.Println(*links)
			item.Url = *links
			title := elemts.MustElement(".glass-product-card__title").MustText()
			item.Title = title
			prices, err := elemts.Element(".gl-price-item")
			price := ""
			if err != nil {
				price = ""
			} else {
				price = prices.MustText()
			}
			item.Precio = price
			item.Marca = "Adidas"
			item.Vendedor = "Adidas"
			var listLinksImg []string
			var link string
			linksImg := elemts.MustElements("img.product-card-image")

			for _, img := range linksImg {
				linkImg, err := img.Attribute("src")
				if err != nil {
					link = ""
				} else {
					link = *linkImg
				}
				listLinksImg = append(listLinksImg, link)
			}
			item.Imagenes = listLinksImg
			listItems = append(listItems, item)
			// fmt.Println(Items)
		}
	}

	return listItems
}

*/
