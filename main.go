package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/gocolly/colly/v2"
)

type Items struct {
	Title  string
	Precio string
}

func main() {

	url := flag.String("url", "", "url mercado libre ej: https://listado.mercadolibre.com.ar/pantalon#D[A:pantalon]")
	name := flag.String("name", "data", "nombre del archivo json")
	flag.Parse()
	namejson := fmt.Sprintf("%s.json", *name)

	c := colly.NewCollector()
	file, err := os.Create(namejson)

	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	listItems := []Items{}

	// Find and visit all links
	c.OnHTML("div.ui-search-result__content", func(container *colly.HTMLElement) {

		// fmt.Println(container.DOM.Html())
		item := Items{}

		container.ForEach("h2", func(i int, titles *colly.HTMLElement) {
			item.Title = titles.Text
		})
		container.ForEach("span.andes-money-amount__fraction", func(i int, prices *colly.HTMLElement) {

			item.Precio = prices.Text
			// fmt.Println(h.Text)

		})
		listItems = append(listItems, item)

	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit(*url)

	json, err := json.Marshal(listItems)
	if err != nil {
		fmt.Println(err)
	}
	file.WriteString(string(json))

}
