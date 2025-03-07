package scraping

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
	"unicode"
	"webScraping/lib"

	"github.com/go-rod/rod"
	"golang.org/x/text/unicode/norm"
)

var isCoincidence = false

func GetDataMercadolibreNike(ctx context.Context, r *http.Request) []Items {
	isCoincidence = true
	proveedor := r.URL.Query().Get("proveedor")
	producSearch := r.URL.Query().Get("marca")
	search := proveedor
	urlSearch := fmt.Sprintf("https://listado.mercadolibre.com.ar/%s", search)
	fmt.Println("entrando en mercado libre ")

	page, err := bm.GetPage(ctx, urlSearch)
	if err != nil {
		fmt.Println("Error al obtener la página:", err)
		return nil
	}
	done := make(chan bool)
	go func() {
		lib.HandlePanicScraping(done, page)
		LoggerInfo("Bucando " + search)
		// page.MustElement("#cb1-edit").MustInput(search)
		// page.MustElement(".nav-search-btn").MustClick()
		page.MustWaitLoad()
		done <- true
	}()
	defer page.Close()
	select {
	case success := <-done:
		var listItems []Items
		if success {
			listItems = scraping(page, proveedor, producSearch)
			if len(listItems) <= 0 {
				listItems = GetDataMercadolibre(ctx, r)
			}
		} else {
			listItems = []Items{}
		}

		fmt.Println("Fin scraping Mercado Libre ")
		return listItems
	case <-ctx.Done():
		// fmt.Println("Timeout o contexto cancelado en Puma ", ctx.Done())
		fmt.Println("Timeout o contexto cancelado en Puma ")
		stringError := fmt.Sprintf("Timeout o contexto cancelado en Puma  %v", ctx.Done())
		LoggerWarning(stringError)
		return []Items{}
	}

}

func GetDataMercadolibreAdidas(ctx context.Context, r *http.Request) []Items {
	isCoincidence = true
	proveedor := r.URL.Query().Get("proveedor")
	producSearch := r.URL.Query().Get("marca")
	search := proveedor
	urlSearch := fmt.Sprintf("https://listado.mercadolibre.com.ar/%s", search)
	fmt.Println("entrando en mercado libre ")

	page, err := bm.GetPage(ctx, urlSearch)

	if err != nil {
		fmt.Println("Error al obtener la página:", err)
		return nil
	}
	defer page.Close()

	done := make(chan bool)

	go func() {
		// Llenar el formulario y hacer clic en el botón de búsqueda
		lib.HandlePanicScraping(done, page)
		// page.MustElement("#cb1-edit").MustInput(search)
		// page.MustElement(".nav-search-btn").MustClick()
		page.MustWaitLoad()
		done <- true
	}()

	select {
	case success := <-done:
		var listItems []Items
		if success {
			listItems = scraping(page, proveedor, producSearch)
			if len(listItems) <= 0 {
				listItems = GetDataMercadolibre(ctx, r)
			}
		} else {
			listItems = []Items{}
		}

		fmt.Println("Fin scraping Mercado Libre ")
		return listItems
	case <-ctx.Done():
		// fmt.Println("Timeout o contexto cancelado en Puma ", ctx.Done())
		fmt.Println("Timeout o contexto cancelado en Puma ")
		stringError := fmt.Sprintf("Timeout o contexto cancelado en Puma  %v", ctx.Done())
		LoggerWarning(stringError)
		return []Items{}
	}
}

func GetDataMercadolibrePuma(ctx context.Context, r *http.Request) []Items {
	isCoincidence = true

	proveedor := r.URL.Query().Get("proveedor")
	producSearch := r.URL.Query().Get("marca")
	search := proveedor

	urlSearch := fmt.Sprintf("https://listado.mercadolibre.com.ar/%s", search)

	fmt.Println("entrando en mercado libre ")
	page, err := bm.GetPage(ctx, urlSearch)
	if err != nil {
		fmt.Println("Error al obtener la página:", err)
		return nil
	}
	defer page.Close()
	done := make(chan bool)
	go func() {
		lib.HandlePanicScraping(done, page)
		// page.MustElement("#cb1-edit").MustInput(search)
		// page.MustElement(".nav-search-btn").MustClick()
		page.MustWaitLoad()
		done <- true
	}()

	select {
	case success := <-done:
		var listItems []Items
		if success {
			listItems = scraping(page, proveedor, producSearch)
			if len(listItems) <= 0 {
				listItems = GetDataMercadolibre(ctx, r)
			}
		} else {
			listItems = []Items{}

		}

		fmt.Println("Fin scraping Mercado Libre ")
		return listItems
	case <-ctx.Done():
		// fmt.Println("Timeout o contexto cancelado en Puma ", ctx.Done())
		fmt.Println("Timeout o contexto cancelado en Puma ")
		stringError := fmt.Sprintf("Timeout o contexto cancelado en Puma  %v", ctx.Done())
		LoggerWarning(stringError)
		return []Items{}
	}
}

func GetDataMercadolibre(ctx context.Context, r *http.Request) []Items {
	coditm := strings.ToLower(strings.Join(strings.Split(r.URL.Query().Get("search"), " "), "-"))
	marca := r.URL.Query().Get("marca")
	categoria := r.URL.Query().Get("categoria")
	genero := r.URL.Query().Get("genero")
	talle := r.URL.Query().Get("Talle")
	material := r.URL.Query().Get("material")

	search := fmt.Sprintf("https://listado.mercadolibre.com.ar/%s-%s-%s-%s", categoria, marca, coditm, genero)
	if material == "SINTETICO" {
		material = "Sintético"
	}
	if marca == "DISTRINANDO (CHOCOLATE)" {
		material = "CHOCOLATE"
	}
	fmt.Println(search)
	// page := rod.New().MustConnect().MustPage("https://listado.mercadolibre.com.ar/mochilas-hombre#D[A:mochilas%20hombre%20]")
	// Launch a headless browser
	fmt.Println("entrando en mercado libre ")
	// page, err := bm.GetPage(ctx, "https://www.mercadolibre.com.ar/")
	page, err := bm.GetPage(ctx, search)
	if err != nil {
		fmt.Println("Error al obtener la página:", err)
		return nil
	}
	defer page.Close()
	done := make(chan bool)
	go func() {
		defer lib.HandlePanicScraping(done, page)
		page.MustWaitLoad()
		// page.MustElement("#cb1-edit").MustInput(search)
		// page.MustElement(".nav-search-btn").MustClick()
		done <- true
	}()

	select {
	case success := <-done:
		var listItems []Items
		if success {
			fils := []string{"Marca:" + marca, "Género:" + genero, "Categorías:" + categoria, "Talle:" + talle, "Material principal:" + material, "Condición:Nuevo", "Tiendas oficiales:Solo tiendas oficiales"}
			getMarc(page, fils)
			listItems = scraping(page, "", "")
			// // Guardar los datos en un archivo JSON
			fmt.Println("fin")
		} else {
			listItems = []Items{}
		}
		return listItems
	case <-ctx.Done():
		// fmt.Println("Timeout o contexto cancelado en Puma ", ctx.Done())
		fmt.Println("Timeout o contexto cancelado en Puma ")
		stringError := fmt.Sprintf("Timeout o contexto cancelado en Puma  %v", ctx.Done())
		LoggerWarning(stringError)
		return []Items{}
	}
}
func scraping(page *rod.Page, proveedor string, producSearch string) []Items {
	fmt.Println("Iniciado scraping")
	page.MustWaitLoad()
	listItems := []Items{}
	containe, err := page.Elements(".ui-search-layout")
	if err != nil {
		fmt.Println("contenido no encontrado")
		return []Items{}
	}
	if len(containe) <= 0 {
		fmt.Println("contenido no encontrado")
		return []Items{}
	}
	element := containe.First().MustElements(".ui-search-layout__item")

	for _, elme := range element {
		elme.WaitVisible()
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Error en scraping MercadoLibre ", r)
			}
		}()

		item := Items{}
		var saller string

		marca, err := elme.Element("span.poly-component__brand")
		if err != nil {
			item.Marca = ""
		} else {
			item.Marca = marca.MustText()
		}

		if isCoincidence {
			marcCompar := strings.ToLower(item.Marca)
			producCompar := strings.ToLower(producSearch)
			coincidence := false
			if strings.Contains(marcCompar, producCompar) {
				coincidence = true
			}
			if !coincidence {
				return []Items{}
			} else {
				isCoincidence = false
			}
		}
		// title := elme.MustElement("h2").MustText()
		// Stitle, err := elme.Element("h2")
		Stitle, err := elme.Element(".poly-component__title-wrapper")
		var title string
		if err != nil {
			title = "---"
		} else {
			title = Stitle.MustText()
		}

		item.Title = title
		isSaller, err := elme.Element("div.poly-card__content > span.poly-component__seller")
		if err != nil {
			saller = ""
		} else {
			saller, err = isSaller.Text()
			if err != nil {
				saller = ""
			}
		}
		var linksImg []string

		imgs := elme.MustElements("img")
		for _, img := range imgs {
			linkimg := img.MustAttribute("src")
			links := *linkimg
			link := strings.Split(links, ".")
			if link[len(link)-1] == "webp" {
				linksImg = append(linksImg, links)
			}
		}
		// links
		links := elme.MustElement("a")
		link := links.MustAttribute("href")
		isCuttentPrice, err := elme.Element(".poly-price__current")
		if err != nil {
			price := elme.MustElement("span.andes-money-amount__fraction").MustText()
			item.Precio = price
		} else {
			item.Precio = isCuttentPrice.MustElement("span.andes-money-amount__fraction").MustText()
			item.PrecioAntiguo = elme.MustElement("span.andes-money-amount__fraction").MustText()
			isExitProcentaje, err := isCuttentPrice.Element(".andes-money-amount__discount")
			if err == nil {
				item.Porcentaje = isExitProcentaje.MustText()
			}
		}
		item.Imagenes = linksImg
		item.Url = *link
		item.Vendedor = saller
		if proveedor != "" {
			item.CodProveedor = proveedor
		}
		if item.Title != "" {
			listItems = append(listItems, item)
		}
	}
	return listItems
}

func getMarc(page *rod.Page, fils []string) {
	for _, fil := range fils {
		fil := strings.Split(fil, ":")
		key, search := fil[0], fil[1]
		if search != "" {
			page = applyFilter(page, key, search)
		}
	}
}

func applyFilter(page *rod.Page, key, filter string) *rod.Page {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error en applyFilter: ", r)
		}
	}()
	for i := 0; i < 1; i++ { // Intentar aplicar el filtro hasta 3 veces
		page.MustWaitLoad()
		time.Sleep(2 * time.Millisecond)
		_, err := page.Eval(`window.scrollTo(0, document.body.scrollHeight)`)
		if err != nil {
			_, err := page.Eval(`window.scrollTo(0, document.body.scrollHeight)`)
			if err != nil {
				fmt.Println("Error al hacer scrol")
			}
		}
		containerFilters, err := page.Elements(".ui-search-filter-groups")
		if err != nil {
			return page
		}
		if len(containerFilters) <= 0 {
			return page
		}
		filters := containerFilters.First().MustElements(".ui-search-filter-dl")
		for _, item := range filters {
			if len(item.MustElements("h3.ui-search-filter-dt-title")) > 0 {
				filtro := item.MustElements("h3.ui-search-filter-dt-title")[0].MustText()
				// fmt.Println(filtro)
				if filtro == key {
					listLink := item.MustElements("a.ui-search-link")
					for i, linkgenero := range listLink {
						if i+1 >= len(listLink) {
							modalsLink, err := item.Element("a.ui-search-modal__link")
							if err == nil {
								if applyFilterWithModal(page, modalsLink, filter) {
									return page
								}
							} else {
								if applyFilterToElement(linkgenero, filter) {
									return page
								}
							}
						}
						if i+1 < len(listLink) {
							if applyFilterToElement(linkgenero, filter) {
								return page
							}
						}
					}
				}
			}
		}
	}
	fmt.Printf("No se pudo aplicar el filtro: %s = %s\n", key, filter)
	return page
}

func applyFilterToElement(linkgenero *rod.Element, filter string) bool {
	targ := normalizeFilter(linkgenero.MustElement("span.ui-search-filter-name").MustText())
	filter = normalizeFilter(filter)
	if filter == "damas" {
		filter = "mujer"
	}
	if strings.Contains(targ, filter) || strings.Contains(targ+"s", filter) {
		linkgenero.MustClick().Page()
		fmt.Println("Se aplico el filtro = ", filter)
		return true
	}
	return false
}

func applyFilterWithModal(page *rod.Page, modalsLink *rod.Element, filter string) bool {
	modalsLink.MustClick()
	time.Sleep(1 * time.Second)
	fmt.Println("abriendo modal")
	modal := page.MustElement("#modal")
	modalItem := modal.MustElements("a.ui-search-link")
	for i, linkgenero := range modalItem {
		targName, err := linkgenero.Text()
		if err == nil {
			targName = normalizeFilter(targName)
			filter = normalizeFilter(filter)

			if filter == "damas" {
				filter = "mujer"
			}
			if strings.Contains(targName, filter) || strings.Contains(targName+"s", filter) {
				linkgenero.MustClick().Page()
				fmt.Println("Se aplico el filtro = ", filter)

				return true
			} else {
				if i+1 == len(modalItem) {
					page.MustElement("button.andes-modal__close-button").MustClick()
					time.Sleep(1 * time.Second)
					fmt.Println("cerrando modal")
				}
				continue
			}

		} else {
			fmt.Println("errr")
			fmt.Println(err)
		}
	}
	return false
}

func normalizeFilter(filter string) string {
	filter = removeAccents(strings.ToLower(filter))
	Splifiltro := strings.Split(filter, " ")
	if len(Splifiltro) > 0 {
		filter = Splifiltro[0]
	}
	return filter
}

func removeAccents(s string) string {
	t := norm.NFD.String(s)
	t = strings.Map(func(r rune) rune {
		if unicode.Is(unicode.Mn, r) {
			return -1
		}
		return r
	}, t)
	return t
}
