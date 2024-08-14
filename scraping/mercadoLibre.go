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

func GetDataMercadolibre(w http.ResponseWriter, r *http.Request) []Items {

	coditm := r.URL.Query().Get("search")
	marca := r.URL.Query().Get("marca")
	categoria := r.URL.Query().Get("categoria")
	genero := r.URL.Query().Get("genero")
	talle := r.URL.Query().Get("Talle")
	material := r.URL.Query().Get("material")

	search := fmt.Sprintf("%s %s %s %s", categoria, coditm, marca, genero)

	if material == "SINTETICO" {
		material = "Sintético"
	}
	if marca == "DISTRINANDO (CHOCOLATE)" {
		material = "CHOCOLATE"
	}

	fmt.Println(search)

	// page := rod.New().MustConnect().MustPage("https://listado.mercadolibre.com.ar/mochilas-hombre#D[A:mochilas%20hombre%20]")

	// Launch a headless browser
<<<<<<< HEAD
	url, err := launcher.New().Headless(true).Launch()
=======
	url, err := launcher.New().
		Headless(true).
		NoSandbox(true).
		Launch()
>>>>>>> 35c02612185f503e82e5bbfadc5a806d09462b12
	if err != nil {
		fmt.Println("Erorrrrrrrr")
		fmt.Println(err)

	}

	browser := rod.New().ControlURL(url).
		MustConnect().
		MustIgnoreCertErrors(false)
	// 	.
	// 	DefaultDevice(proto.EmulationSetUserAgentOverride{
	// 		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36",
	// 	})

	// // Desactivar la carga de imágenes
	// browser.MustSetExtraHTTPHeaders(proto.NetworkSetExtraHTTPHeaders{
	// 	Headers: map[string]string{
	// 		"Sec-Fetch-Dest": "document",
	// 		"Sec-Fetch-Mode": "navigate",
	// 		"Sec-Fetch-Site": "same-origin",
	// 		"Accept":         "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
	// 	},
	// })
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
	fmt.Println("Iniciado scraping")
	listItems = scraping(page, coditm)

	/*

		db, err := sql.Open("sqlite3", "scraping.db")
		if err != nil {
			log.Fatal("err")
			log.Fatal(err)
		}
		defer db.Close()d

		for _, item := range listItems {
			res, err := db.Exec("insert into WebItems(Title,Precio) values(?,?)", item.Title, item.Precio)

			if err != nil {
				log.Fatal(err)
			}
			id, err := res.LastInsertId()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(id)

		}
	*/
	// // Guardar los datos en un archivo JSON
	fmt.Println("fin")
	return listItems
	// w.Header().Set("Content-Type", "application/json")

	// fmt.Fprintf(w, "Buscando %s", &jsonData)
	// json.NewEncoder(w).Encode(listItems)

}

func scraping(page *rod.Page, search string) []Items {
	page.MustWaitLoad()
	listItems := []Items{}

	containe := page.MustElement(".ui-search-layout")
	element := containe.MustElements(".ui-search-layout__item")

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

		// fmt.Println(title)
		// filtro de titulo

		listSeach := strings.Split(search, " ")
		listTiele := strings.Split(title, " ")

		count := 0
		for _, itemSearch := range listSeach {

			for _, lTitel := range listTiele {

				if strings.Contains(strings.ToLower(itemSearch), strings.ToLower(lTitel)) {

					count++
				}
			}

		}

		// fmt.Println("coun ", count)
		// fmt.Println("coun2 ", len(listSeach))
		if count >= len(listSeach)/2 {
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
			marca := elme.MustElement("span")
			if marca.MustText() == "" {
				marca = elme.MustElement("span.ui-search-item__brand-discoverability")
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

			// if strings.Contains(strings.ToUpper(title), "ZAPATILLA CARINA 2.0 LUX ADP") {
			// 	item.Title = title
			// 	fmt.Println("nop")
			// 	// Aquí puedes hacer lo que necesites con el producto encontrado, como extraer más información
			// } else {
			// 	item.Title = ""

			// }

			price := elme.MustElement("span.andes-money-amount__fraction").MustText()
			// item.Title = title
			item.Precio = price
			item.Marca = marca.MustText()
			item.Imagenes = linksImg
			item.Url = *link
			item.Vendedor = saller
			if item.Title != "" {

				listItems = append(listItems, item)
			}
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
		containerFilters := page.MustElement(".ui-search-filter-groups")
		filters := containerFilters.MustElements(".ui-search-filter-dl")

		for _, item := range filters {
			if len(item.MustElements("h3.ui-search-filter-dt-title")) > 0 {
				filtro := item.MustElements("h3.ui-search-filter-dt-title")[0].MustText()
				// fmt.Println(filtro)
				// fmt.Println(filtro)
				if filtro == key {
					listLink := item.MustElements("a.ui-search-link")
					time.Sleep(1 * time.Second)

					for i, linkgenero := range listLink {
						if i+1 >= len(listLink) {
							if applyFilterWithModal(page, item, filter) {
								return page
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
		fmt.Println(" no modal -> tageta  " + targ + " filtro  " + filter)
		linkgenero.MustClick().Page()
		fmt.Println("click")
		return true
	}
	return false
}

func applyFilterWithModal(page *rod.Page, item *rod.Element, filter string) bool {
	modalsLink, err := item.Element("a.ui-search-modal__link")
	if err != nil {
		fmt.Println("sin modal")
		return false
	}
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
				fmt.Println("tageta  " + targName)
				linkgenero.MustClick().Page()
				fmt.Println("click")
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
