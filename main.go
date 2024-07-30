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

	fmt.Println(search)

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
		link := getMarc(page, marca)
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
func getMarc(page *rod.Page, marca string) string {

	Containerfilters := page.MustElement(".ui-search-filter-groups")

	filters := Containerfilters.MustElements(".ui-search-filter-dl")

	var linkFilter string
	for _, item := range filters {

		if len(item.MustElements("h3.ui-search-filter-dt-title")) > 0 && item.MustElements("h3.ui-search-filter-dt-title")[0].MustText() == "Marca" {

			fmt.Println(item.MustText())
			link, err := item.Element("a.ui-search-modal__link")
			if err != nil {

				Listlink := item.MustElements("a.ui-search-link")

				for _, modalLink := range Listlink {

					targ := removeAccents(strings.ToLower(modalLink.MustText()))
					fmt.Println(targ)
					marca := removeAccents(strings.ToLower(marca))
					if strings.Contains(targ, marca) {
						link := modalLink.MustAttribute("href")

						linkFilter = *link
						break

					}
				}
			} else {
				link.MustClick()
				time.Sleep(1 * time.Second)
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
