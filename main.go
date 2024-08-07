package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"unicode"

	"github.com/go-rod/rod"
	"golang.org/x/text/unicode/norm"
)

type Items struct {
	Title    string
	Precio   string
	Marca    string
	Url      string
	Imagenes []string
}

func main() {

	http.HandleFunc("GET /api/data", GetDataMercadolibre)

	http.ListenAndServe(":3000", AddCORSHeaders(http.DefaultServeMux))

}

func GetDataMercadolibre(w http.ResponseWriter, r *http.Request) {

	coditm := r.URL.Query().Get("search")
	marca := r.URL.Query().Get("marca")
	categoria := r.URL.Query().Get("categoria")
	genero := r.URL.Query().Get("genero")
	talle := r.URL.Query().Get("Talle")
	material := r.URL.Query().Get("material")

	search := fmt.Sprintf("%s %s %s %s %s %s", coditm, marca, categoria, material, genero, talle)

	fmt.Println(search)

	// page := rod.New().MustConnect().MustPage("https://listado.mercadolibre.com.ar/mochilas-hombre#D[A:mochilas%20hombre%20]")

	browser := rod.New().MustConnect()
	defer browser.Close()

	fmt.Println("entrando en mercado libre ")
	page := browser.MustPage("https://www.mercadolibre.com.ar/")

	// Llenar el formulario y hacer clic en el botón de búsqueda
	page.MustElement("#cb1-edit").MustInput(search)
	page.MustElement(".nav-search-btn").MustClick()

	// class="ui-search-layout__item"
	// ui-search-filter-groups

	var listItems []Items

	fils := []string{"Marca:" + marca, "Género:" + genero, "Categorías:" + categoria, "Talle:" + talle}

	getMarc(page, fils)
	fmt.Println("Iniciado scraping")
	listItems = scraping(page)

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
	w.Header().Set("Content-Type", "application/json")

	// fmt.Fprintf(w, "Buscando %s", &jsonData)
	json.NewEncoder(w).Encode(listItems)

}

// func wr(n string, d string) {
// 	file, err := os.Create(n + ".html")
// 	if err != nil {
// 		fmt.Println("ewr")
// 		fmt.Println(err)
// 	}
// 	defer file.Close()
// 	file.WriteString(d)

// }
func scraping(page *rod.Page) []Items {
	page.MustWaitLoad()
	listItems := []Items{}

	containe := page.MustElement(".ui-search-layout")
	element := containe.MustElements(".ui-search-layout__item")

	for _, elme := range element {

		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recover en scraping")
			}
		}()

		item := Items{}
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

		title := elme.MustElement("h2").MustText()
		price := elme.MustElement("span.andes-money-amount__fraction").MustText()
		item.Title = title
		item.Precio = price
		item.Marca = marca.MustText()
		item.Imagenes = linksImg
		item.Url = *link

		listItems = append(listItems, item)

	}
	return listItems
}

func getMarc(page *rod.Page, fils []string) {

	fmt.Println("Buscnado filtros")

	for _, fil := range fils {
		fil := strings.Split(fil, ":")
		key, search := fil[0], fil[1]
		fmt.Printf("%s = %s \n", key, search)
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

				if filtro == key {
					listLink := item.MustElements("a.ui-search-link")
					// fmt.Println(listLink)
					time.Sleep(1 * time.Second)
					for i, linkgenero := range listLink {
						if applyFilterToElement(linkgenero, filter) {
							return page
						}

						if i+1 == len(listLink) {
							if applyFilterWithModal(page, item, filter) {
								return page
							}
						}
					}
				}
			}
		}
		time.Sleep(2 * time.Second) // Esperar antes de intentar nuevamente
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
	modalsLink := item.MustElement("a.ui-search-modal__link")
	modalsLink.MustClick()
	fmt.Println("En el modal")
	time.Sleep(1 * time.Second)
	fmt.Println("abriendo modal")
	modal := page.MustElement("#modal")
	modalItem := modal.MustElements("a.ui-search-link")

	for _, linkgenero := range modalItem {
		targName, err := linkgenero.Text()

		if err == nil {
			targName = normalizeFilter(targName)
			filter = normalizeFilter(filter)

			if filter == "damas" {
				filter = "mujer"
			}
			if strings.Contains(targName, filter) || strings.Contains(targName+"s", filter) {
				fmt.Println("tageta  " + targName)
				page = linkgenero.MustClick().Page()
				fmt.Println("click")
				return true
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

func AddCORSHeaders(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		handler.ServeHTTP(w, r)
	})
}

/*

func appliFiltro(page *rod.Page, key string, filter string) *rod.Page {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recoverd en el appliFilter: ", r)
		}
	}()

	page.MustWaitLoad()

	Containerfilters := page.MustElement(".ui-search-filter-groups")
	filters := Containerfilters.MustElements(".ui-search-filter-dl")

	for _, item := range filters {

		if len(item.MustElements("h3.ui-search-filter-dt-title")) > 0 {
			filtro := item.MustElements("h3.ui-search-filter-dt-title")[0].MustText()

			if filtro == key {

				Listlink := item.MustElements("a.ui-search-link")

				for i, linkgenero := range Listlink {
					// linkgenero.MustElement("span.ui-search-filter-name").MustText()
					// fmt.Println(linkgenero.MustElement("span.ui-search-filter-name").HTML())

					modalsLink, err := item.Element("a.ui-search-modal__link")

					if err == nil {

						targ := removeAccents(strings.ToLower(linkgenero.MustElement("span.ui-search-filter-name").MustText()))

						Splifiltro := strings.Split(strings.ToLower(filter), " ")
						if len(Splifiltro) > 0 {

							filter = removeAccents(Splifiltro[0])
						} else {
							filter = removeAccents(strings.ToLower(filter))

						}

						if filter == "damas" {
							filter = "mujer"
						}

						if strings.Contains(targ, filter) || strings.Contains(targ+"s", filter) {
							fmt.Println(" no modal -> tageta  "+targ, " filtro  "+filter)

							linkgenero.MustClick()
							page = linkgenero.Page()

							fmt.Println("click")

							return page

						}

					} else {

						if i+1 == len(Listlink) {

							fmt.Println("En el modal")

							modalsLink.MustClick()

							fmt.Println(modalsLink.HTML())

							time.Sleep(1 * time.Second)
							fmt.Println("abriendo modal")
							modal := page.MustElement("#modal")
							// ui-search-search-modal-list
							// .ui-search-search-modal-grid-columns
							modalItem := modal.MustElements("a.ui-search-link")

							for _, linkgenero := range modalItem {
								targName, err := linkgenero.Text()

								if err != nil {
									fmt.Println("errr")
									fmt.Println(err)
								} else {
									targName = removeAccents(strings.ToLower(targName))
									Splifiltro := strings.Split(strings.ToLower(filter), " ")
									if len(Splifiltro) > 0 {

										filter = removeAccents(Splifiltro[0])
									} else {
										filter = removeAccents(strings.ToLower(filter))

									}
									// filter = removeAccents(strings.ToLower(filter))
									if filter == "damas" {
										filter = "mujer"
									}
									if strings.Contains(targName, filter) || strings.Contains(targName+"s", filter) {
										fmt.Println("tageta  " + targName)

										page = linkgenero.MustClick().Page()
										// wr("ht/"+targ, page.MustHTML())
										fmt.Println("click")

										return page

									}

								}

							}
						} else {

							targ := removeAccents(strings.ToLower(linkgenero.MustElement("span.ui-search-filter-name").MustText()))

							Splifiltro := strings.Split(strings.ToLower(filter), " ")
							if len(Splifiltro) > 0 {

								filter = removeAccents(Splifiltro[0])
							} else {
								filter = removeAccents(strings.ToLower(filter))

							}

							if filter == "damas" {
								filter = "mujer"
							}

							if strings.Contains(targ, filter) || strings.Contains(targ+"s", filter) {
								fmt.Println(" no modal2  -> tageta  "+targ, " filtro  "+filter)

								linkgenero.MustClick()
								page = linkgenero.Page()

								fmt.Println("click")

								return page

							}

						}

					}

				}
			}
		}

	}
	return page

}

*/

/*
	for i, linkgenero := range Listlink {
					// linkgenero.MustElement("span.ui-search-filter-name").MustText()
					// fmt.Println(linkgenero.MustElement("span.ui-search-filter-name").HTML())




					if i+1 == len(Listlink) {

						modalsLink, err := item.Element("a.ui-search-modal__link")

						if err == nil {
							fmt.Println("En el modal")

							modalsLink.MustClick()

							time.Sleep(1 * time.Second)
							fmt.Println("abriendo modal")
							modal := page.MustElement("#modal")
							// ui-search-search-modal-list
							// .ui-search-search-modal-grid-columns
							modalItem := modal.MustElements("a.ui-search-link")

							for _, linkgenero := range modalItem {
								targName, err := linkgenero.Text()

								if err != nil {
									fmt.Println("errr")
									fmt.Println(err)
								} else {
									targName = removeAccents(strings.ToLower(targName))
									Splifiltro := strings.Split(strings.ToLower(filter), " ")
									if len(Splifiltro) > 0 {

										filter = removeAccents(Splifiltro[0])
									} else {
										filter = removeAccents(strings.ToLower(filter))

									}
									// filter = removeAccents(strings.ToLower(filter))
									if filter == "damas" {
										filter = "mujer"
									}
									if strings.Contains(targName, filter) || strings.Contains(targName+"s", filter) {
										fmt.Println("tageta  " + targName)

										page = linkgenero.MustClick().Page()
										// wr("ht/"+targ, page.MustHTML())
										fmt.Println("click")

										return page

									}

								}

							}
						} else {

							targ := removeAccents(strings.ToLower(linkgenero.MustElement("span.ui-search-filter-name").MustText()))

							Splifiltro := strings.Split(strings.ToLower(filter), " ")
							if len(Splifiltro) > 0 {

								filter = removeAccents(Splifiltro[0])
							} else {
								filter = removeAccents(strings.ToLower(filter))

							}

							if filter == "damas" {
								filter = "mujer"
							}
							fmt.Println(" no modal2  -> tageta  "+targ, " filtro  "+filter)

							if strings.Contains(targ, filter) || strings.Contains(targ+"s", filter) {

								linkgenero.MustClick()
								page = linkgenero.Page()

								fmt.Println("click")

								return page

							}
						}

					} else {

						targ := removeAccents(strings.ToLower(linkgenero.MustElement("span.ui-search-filter-name").MustText()))

						Splifiltro := strings.Split(strings.ToLower(filter), " ")
						if len(Splifiltro) > 0 {

							filter = removeAccents(Splifiltro[0])
						} else {
							filter = removeAccents(strings.ToLower(filter))

						}

						if filter == "damas" {
							filter = "mujer"
						}

						fmt.Println(" no modal -> tageta  "+targ, " filtro  "+filter)
						if strings.Contains(targ, filter) || strings.Contains(targ+"s", filter) {

							linkgenero.MustClick()
							page = linkgenero.Page()

							fmt.Println("click")

							return page

						}
					}

				}

*/
