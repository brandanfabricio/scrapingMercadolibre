package lib

import (
	"fmt"
	"log/slog"
	"os"
)

func LoggerInfo(msj string) {

	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("error")
	}
	defer file.Close()
	logger := slog.New(slog.NewJSONHandler(file, nil))

	logger.Info("Inicio de busqueda",
		"Buscar: ", msj)

}
func LoggerWarning(msj string) {

	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("error")
	}
	defer file.Close()
	logger := slog.New(slog.NewJSONHandler(file, nil))

	logger.Warn("Advertencia",
		"Warning: ", msj)

}
func LoggerError(msj string) {

	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("error")
	}
	defer file.Close()
	logger := slog.New(slog.NewJSONHandler(file, nil))

	logger.Error("Error en la busqueda",
		"ERROR: ", msj)

}
