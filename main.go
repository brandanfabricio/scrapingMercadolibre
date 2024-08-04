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
	Title  string
	Precio string
	Marca  string
}

func main() {

	http.HandleFunc("GET /api/data", GetDataMercadolibre)

	http.ListenAndServe(":3000", AddCORSHeaders(http.DefaultServeMux))

}

func GetDataMercadolibre(w http.ResponseWriter, r *http.Request) {

	search := r.URL.Query().Get("search")
	marca := r.URL.Query().Get("marca")
	categoria := r.URL.Query().Get("categoria")
	genero := r.URL.Query().Get("genero")

	fmt.Println(search + " " + marca + " " + categoria)

	// page := rod.New().MustConnect().MustPage("https://listado.mercadolibre.com.ar/mochilas-hombre#D[A:mochilas%20hombre%20]")

	browser := rod.New().MustConnect()

	defer browser.Close()
	fmt.Println("entrando en mercado libre ")
	page := browser.MustPage("https://www.mercadolibre.com.ar/")

	// Llenar el formulario y hacer clic en el botón de búsqueda
	fmt.Println("Buscando ", search)
	page.MustElement("#cb1-edit").MustInput(search + " " + marca)
	page.MustElement(".nav-search-btn").MustClick()

	// class="ui-search-layout__item"
	// ui-search-filter-groups
	var listItems []Items

	if marca == "" {
		listItems = scraping(page)

	} else {

		fmt.Println("Buscando Marca")
		link := getMarc(page, marca, categoria, genero)
		fmt.Println(link)
		if link != "" {

			fmt.Println("scraping Marca")
			listItems = scraping(browser.MustPage(link))
		} else {
			fmt.Println("scraping")
			listItems = scraping(page)

		}

	}
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

	fmt.Println("guardando")

	w.Header().Set("Content-Type", "application/json")

	// fmt.Fprintf(w, "Buscando %s", &jsonData)
	json.NewEncoder(w).Encode(listItems)

}
func scraping(page *rod.Page) []Items {
	page.MustWaitLoad()
	listItems := []Items{}

	containe := page.MustElement(".ui-search-layout")
	element := containe.MustElements(".ui-search-layout__item")

	for _, elme := range element {
		item := Items{}
		marca := elme.MustElement("span")
		if marca.MustText() == "" {
			marca = elme.MustElement("span.ui-search-item__brand-discoverability")
		}
		// pr := elMarca == marca
		title := elme.MustElement("h2").MustText()
		price := elme.MustElement("span.andes-money-amount__fraction").MustText()
		item.Title = title
		item.Precio = price
		item.Marca = marca.MustText()
		listItems = append(listItems, item)
	}
	return listItems
}
func getMarc(page *rod.Page, marca string, categoria, genero string) string {

	Containerfilters := page.MustElement(".ui-search-filter-groups")

	filters := Containerfilters.MustElements(".ui-search-filter-dl")
	fmt.Println("Buscnado filtros")
	var linkFilter string
	for _, item := range filters {
		// fmt.Println(item.MustElements("h3.ui-search-filter-dt-title")[0].MustText())

		if len(item.MustElements("h3.ui-search-filter-dt-title")) > 0 {

			filtro := item.MustElements("h3.ui-search-filter-dt-title")[0].MustText()
			// fmt.Println(filtro)
			if filtro == "Género" || filtro == "Genero" {

				link, err := item.Element("a.ui-search-modal__link")

				if err != nil {
					fmt.Println("error")
					Listlink := item.MustElements("a.ui-search-link")

					for _, linkgenero := range Listlink {

						targ := removeAccents(strings.ToLower(linkgenero.MustElement("span.ui-search-filter-name").MustText()))
						genero = removeAccents(strings.ToLower(genero))

						if strings.Contains(targ, genero) || strings.Contains(targ+"s", genero) {

							fmt.Println("antes del ewrr")
							fmt.Println(linkgenero)

							item = linkgenero.MustClick()

							// linkgenero.MustClick()
							// linkFilter = *link

							if linkgenero == nil {
								fmt.Println("linkgenero es nil, no se puede hacer click")
							} else {
								link := linkgenero.MustAttribute("href") // Descomentado
								fmt.Println(*link)
								linkFilter = *link

							}
						}

					}

					// fmt.Println(Listlink)
					// link.MustClick()
					// time.Sleep(1 * time.Second)

				} else {
					fmt.Println(link)

				}

			}
			if filtro == "Categorías" {
				link, err := item.Element("a.ui-search-modal__link")
				if err != nil {
					fmt.Println("error")
					Listlink := item.MustElements("a.ui-search-link")

					for _, categoriLink := range Listlink {

						targ := removeAccents(strings.ToLower(categoriLink.MustElement("span.ui-search-filter-name").MustText()))
						// categoria := removeAccents(strings.ToLower(categoria))

						if strings.Contains(targ, categoria) {
							categoriLink.MustClick()

						}

					}

					// fmt.Println(Listlink)
					// link.MustClick()
					// time.Sleep(1 * time.Second)

				} else {
					fmt.Println(link)

				}
			}

			if filtro == "Marca" {

				link, err := item.Element("a.ui-search-modal__link")
				if err != nil {

					Listlink := item.MustElements("a.ui-search-link")

					for _, modalLink := range Listlink {
						fmt.Println("cuscando en el listado")

						targ := removeAccents(strings.ToLower(modalLink.MustText()))
						marcas := removeAccents(strings.ToLower(marca))
						marca := strings.Split(marcas, " ")
						for _, marc := range marca {
							if strings.Contains(targ, marc) {
								link := modalLink.MustAttribute("href")

								linkFilter = *link
								break

							}
						}

					}
				} else {
					link.MustClick()
					time.Sleep(1 * time.Second)
					fmt.Println("abriendo modal")
					modal := page.MustElement("#modal")
					// ui-search-search-modal-list
					// .ui-search-search-modal-grid-columns
					modalItem := modal.MustElements("a.ui-search-link")
					for _, modalLink := range modalItem {
						targ := removeAccents(strings.ToLower(modalLink.MustText()))
						marca := removeAccents(strings.ToLower(marca))
						if strings.Contains(targ, marca) {
							link := modalLink.MustAttribute("href")
							linkFilter = *link
							break

						}
					}
				}

				break

				// linkFilter = link.MustAttribute("href")
			}
		}

	}
	fmt.Println("aki2")

	return linkFilter

	// page = browser.MustPage(*linkFilter)
	// page.MustWaitLoad()
	// modal := page.MustElements("a")
	// for _, item := range modal {
	// 	targ := strings.ToLower(item.MustText())

	// 	marca := strings.ToLower(marca)

	// 	if strings.Contains(targ, marca[0:2]) {
	// 		link = item.MustAttribute("href")
	// 	}
	// }

}

func removeAccents(s string) string {
	t := norm.NFD.String(s)
	t = strings.Map(func(r rune) rune {
		if unicode.Is(unicode.Mn, r) { // Mn: Nonspacing_Mark
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
