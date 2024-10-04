package lib

import (
	"fmt"
	"log"

	"github.com/go-rod/rod"
)

func HandlePanic() {
	if r := recover(); r != nil {
		log.Printf("######################################")
		log.Println("Control de panico")
		log.Printf("Recuperado del pánico: %v", r)
		log.Printf("######################################")
		// fmt.Println("Timeout o contexto cancelado en Puma ")
		stringError := fmt.Sprintf("Timeout o contexto cancelado en Puma  %v", r)
		LoggerWarning(stringError)
	}
}

func HandlePanicScraping(done chan bool, page *rod.Page) {
	if r := recover(); r != nil {
		log.Printf("######################################")
		log.Println("Control de panico")
		// log.Printf("Recuperado del pánico: %v", r)
		stringError := fmt.Sprintf("Recuperado del pánico %v", r)
		log.Printf("######################################")
		// fmt.Println("Timeout o contexto cancelado en Puma ")
		LoggerError(stringError)
		page.Close()
		done <- false
		return
	}
}
