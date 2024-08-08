package scraping

import (
	"fmt"
	"os"
)

type Items struct {
	Title    string
	Precio   string
	Marca    string
	Url      string
	Vendedor string

	Imagenes []string
}

func Wr(n string, d string) {
	file, err := os.Create(n + ".html")
	if err != nil {
		fmt.Println("ewr")
		fmt.Println(err)
	}
	defer file.Close()
	file.WriteString(d)

}
