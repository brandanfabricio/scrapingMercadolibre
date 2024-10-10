package scraping

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"webScraping/lib"

	"github.com/go-rod/rod"
)

func GetDataAdidas(ctx context.Context, r *http.Request) []Items {
	proveedor := r.URL.Query().Get("proveedor")
	search := r.URL.Query().Get("search")
	var urlSearch string
	if proveedor == "" {
		urlSearch = fmt.Sprintf("https://www.adidas.com.ar/search?q=%s", search)
	} else {

		urlSearch = fmt.Sprintf("https://www.adidas.com.ar/search?q=%s", proveedor)
	}
	fmt.Println("entrando en Adidas ")
	fmt.Println(urlSearch)
	LoggerInfo(urlSearch)
	page, err := bm.GetPage(ctx, urlSearch)

	if err != nil {
		fmt.Println("Error al obtener la página:", err)
		return nil
	}
	defer page.Close()
	done := make(chan bool)

	go func() {
		lib.HandlePanicScraping(done, page)
		page.MustWaitLoad()
		done <- true
	}()

	select {
	case success := <-done:
		var listItems []Items
		if success {
			listItems = scrapingAdidas(page, proveedor)
		} else {
			listItems = []Items{}
		}
		fmt.Println("fin adidas")
		return listItems
	case <-ctx.Done():
		// fmt.Println("Timeout o contexto cancelado en Puma ", ctx.Done())
		fmt.Println("Timeout o contexto cancelado en Puma ")
		stringError := fmt.Sprintf("Timeout o contexto cancelado en Puma  %v", ctx.Done())
		LoggerWarning(stringError)
		return []Items{}
	}

}
func scrapingAdidas(page *rod.Page, proveedor string) []Items {
	page.MustWaitLoad()
	var listItems []Items
	fmt.Println("iniciando scraping")

	if proveedor == "" {

		listItems = scrapingList(page)

	} else {
		listItems = scrapingPage(page, proveedor)
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

func scrapingPage(page *rod.Page, proveedor string) []Items {

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

				porcerntaje := int(((anti - nuevo) / anti) * 100)
				item.Porcentaje = fmt.Sprintf("%d%%", porcerntaje)

			}

		}
		// gl-price-item--sale

		hola
		//
		item.CodProveedor = proveedor

		item.Url = page.MustInfo().URL
		item.Marca = "Adidas"
		item.Vendedor = "Adidas"
		var wg sync.WaitGroup
		var ListLinksImg []string
		constainerImages := container.MustElement(".image-grid___1JN2z")
		linksImgs := constainerImages.MustElements("img")

		var link string
		for i, img := range linksImgs {
			wg.Add(1)
			lib.HandlePanic()
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
