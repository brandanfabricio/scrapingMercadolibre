package scraping

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"unicode"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"golang.org/x/text/unicode/norm"
)

func GetDataMercadolibreNike(w http.ResponseWriter, r *http.Request) []Items {
	proveedor := r.URL.Query().Get("proveedor")
	search := proveedor
	url, err := launcher.New().
		Headless(true).
		NoSandbox(true).
		Launch()
	if err != nil {
		http.Error(w, "Error launching browser", http.StatusInternalServerError)
		return nil
	}
	browser := rod.New().ControlURL(url).
		MustConnect().
		MustIgnoreCertErrors(false)
	defer browser.Close()
	fmt.Println("entrando en mercado libre ")
	page := browser.MustPage("https://www.mercadolibre.com.ar/")
	// Llenar el formulario y hacer clic en el botón de búsqueda
	LoggerInfo("Bucando " + search)
	page.MustElement("#cb1-edit").MustInput(search)
	page.MustElement(".nav-search-btn").MustClick()
	listItems := scraping(page, proveedor)
	if len(listItems) <= 0 {
		LoggerInfo(fmt.Sprintf("No se encontro producto de nike por codigo de proveedor  %s ,Buscando por descripcion ", proveedor))
		return GetDataMercadolibre(w, r)
	}
	fmt.Println("Fin scraping Mercado Libre ")
	return listItems

}

func GetDataMercadolibreAdidas(w http.ResponseWriter, r *http.Request) []Items {
	proveedor := r.URL.Query().Get("proveedor")
	search := proveedor
	url, err := launcher.New().
		Headless(true).
		NoSandbox(true).
		Launch()
	if err != nil {
		http.Error(w, "Error launching browser", http.StatusInternalServerError)
		return nil
	}
	browser := rod.New().ControlURL(url).
		MustConnect().
		MustIgnoreCertErrors(false)
	defer browser.Close()
	fmt.Println("entrando en mercado libre ")
	page := browser.MustPage("https://www.mercadolibre.com.ar/")
	// Llenar el formulario y hacer clic en el botón de búsqueda
	page.MustElement("#cb1-edit").MustInput(search)
	page.MustElement(".nav-search-btn").MustClick()
	listItems := scraping(page, proveedor)
	fmt.Println("Fin scraping Mercado Libre ")
	return listItems

}

func GetDataMercadolibrePuma(w http.ResponseWriter, r *http.Request) []Items {
	proveedor := r.URL.Query().Get("proveedor")
	search := proveedor
	url, err := launcher.New().
		Headless(true).
		NoSandbox(true).
		Launch()
	if err != nil {
		http.Error(w, "Error launching browser", http.StatusInternalServerError)
		return nil
	}
	browser := rod.New().ControlURL(url).
		MustConnect().
		MustIgnoreCertErrors(false)
	defer browser.Close()
	fmt.Println("entrando en mercado libre ")
	page := browser.MustPage("https://www.mercadolibre.com.ar/")
	// Llenar el formulario y hacer clic en el botón de búsqueda
	page.MustElement("#cb1-edit").MustInput(search)
	page.MustElement(".nav-search-btn").MustClick()
	listItems := scraping(page, proveedor)
	fmt.Println("Fin scraping Mercado Libre ")
	return listItems

}

func GetDataMercadolibre(w http.ResponseWriter, r *http.Request) []Items {
	coditm := r.URL.Query().Get("search")
	marca := r.URL.Query().Get("marca")
	categoria := r.URL.Query().Get("categoria")
	genero := r.URL.Query().Get("genero")
	talle := r.URL.Query().Get("Talle")
	material := r.URL.Query().Get("material")
	search := fmt.Sprintf("%s %s %s", categoria, marca, coditm)
	if material == "SINTETICO" {
		material = "Sintético"
	}
	if marca == "DISTRINANDO (CHOCOLATE)" {
		material = "CHOCOLATE"
	}
	fmt.Println(search)
	// page := rod.New().MustConnect().MustPage("https://listado.mercadolibre.com.ar/mochilas-hombre#D[A:mochilas%20hombre%20]")
	// Launch a headless browser
	url, err := launcher.New().
		Headless(true).
		NoSandbox(true).
		Launch()
	if err != nil {
		fmt.Println("Erorrrrrrrr")
		fmt.Println(err)
	}
	browser := rod.New().ControlURL(url).
		MustConnect().
		MustIgnoreCertErrors(false)
	defer browser.Close()
	fmt.Println("entrando en mercado libre ")
	page := browser.MustPage("https://www.mercadolibre.com.ar/")
	// Llenar el formulario y hacer clic en el botón de búsqueda
	page.MustElement("#cb1-edit").MustInput(search)
	page.MustElement(".nav-search-btn").MustClick()
	// class="ui-search-layout__item"
	// ui-search-filter-groups
	var listItems []Items
	fils := []string{"Marca:" + marca, "Género:" + genero, "Categorías:" + categoria, "Talle:" + talle, "Material principal:" + material, "Condición:Nuevo", "Tiendas oficiales:Solo tiendas oficiales"}
	getMarc(page, fils)
	listItems = scraping(page, "")
	// // Guardar los datos en un archivo JSON
	fmt.Println("fin")
	return listItems
}

func scraping(page *rod.Page, proveedor string) []Items {
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
				fmt.Println("Recover en scraping")
			}
		}()
		item := Items{}
		var saller string
		title := elme.MustElement("h2").MustText()
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
		marca, err := elme.Element("span.poly-component__brand")
		if err != nil {
			item.Marca = ""
		} else {
			item.Marca = marca.MustText()
		}
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
			fmt.Println("Recovered en applyFilter: ", r)
		}
	}()
	for i := 0; i < 1; i++ { // Intentar aplicar el filtro hasta 3 veces
		page.MustWaitLoad()
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
