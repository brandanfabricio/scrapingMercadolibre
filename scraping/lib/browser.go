package lib

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

// NewBrowserManager crea una nueva instancia de BrowserManager.
type BrowserManager struct {
	browser   *rod.Browser
	mu        sync.Mutex
	once      sync.Once
	idleTimer *time.Timer
}

// NewBrowserManager crea una nueva instancia de BrowserManager.
func NewBrowserManager() *BrowserManager {
	return &BrowserManager{}
}

// initializeBrowser inicializa la instancia del navegador.
func (brm *BrowserManager) initializeBrowser() {
	path, _ := launcher.LookPath()
	fmt.Println(path)
	laucher, err := launcher.New().
		Bin(path).
		Headless(true). // Modo headless para reducir consumo de recursos
		NoSandbox(true).
		Leakless(false).
		Devtools(false).
		Set("disable-web-security"). // Desactivar seguridad web (CORS)
		Set("disable-extensions").   // Desactivar extensiones
		Set("disable-blink-features", "AutomationControlled").
		Launch()
	if err != nil {
		fmt.Println("Error en creacion de laucher ", err)
	}
	brm.browser = rod.New().
		ControlURL(laucher).MustConnect().Timeout(70 * time.Second).CancelTimeout()
	// Habilitar caché del navegador
	// brm.browser.MustPage().SetCacheEnabled

	// Configurar el User-Agent y otros headers
	brm.browser.MustPage().MustSetUserAgent(&proto.NetworkSetUserAgentOverride{
		UserAgent:      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, como Gecko) Chrome/91.0.4472.124 Safari/537.36",
		AcceptLanguage: "en-US,en;q=0.9",
		Platform:       "Windows",
	})

	// Iniciar el temporizador de inactividad
	brm.startIdleTimer()
}

// GetPage devuelve una nueva página utilizando la instancia compartida del navegador.
func (brm *BrowserManager) GetPage(ctx context.Context, url string) (*rod.Page, error) {
	brm.once.Do(brm.initializeBrowser) // Inicializa el navegador solo una vez
	brm.mu.Lock()
	defer brm.mu.Unlock() // Bloqueo para acceso concurrente

	// Verifica si el navegador fue cerrado y vuelve a inicializar si es necesario
	// fmt.Println(brm.browser)
	_, err := brm.browser.Incognito()
	if err != nil {
		brm.initializeBrowser()
	}
	if brm.browser == nil {
		fmt.Println("Navegador cerrado. Reinicializando...")
		brm.once = sync.Once{} // Permitir que el navegador se vuelva a inicializar
		brm.initializeBrowser()
	}

	// Reinicia el temporizador de inactividad cada vez que se llama a GetPage
	brm.startIdleTimer()
	page := brm.browser.Timeout(70 * time.Second).CancelTimeout().MustPage(url) // Timeout extendido para manejar internet lento

	page.Mouse.MustMoveTo(100, 45) // Simular movimiento del mouse para evitar detección de bots

	// Controla la cancelación por contexto
	select {
	case <-ctx.Done():
		fmt.Println("Context cancelado Navegador", ctx.Err())
		defer page.Close() // Cierra la página si el contexto es cancelado
		return nil, ctx.Err()
	default:
		return page, nil
	}
}

// Close cierra la instancia del navegador.
func (brm *BrowserManager) Close() {
	if brm.browser != nil {
		brm.browser.MustClose()
	}
}
func (brm *BrowserManager) KillChromeProcesses() {
	// cmd := exec.Command("pkill", "-f", "chrome")
	cmd := exec.Command("killall", "chrome")
	err := cmd.Run()
	if err != nil {
		log.Printf("Error matando procesos de Chrome: %v", err)
	} else {
		log.Println("Procesos de Chrome terminados exitosamente.")
	}
}

// startIdleTimer inicia un temporizador para cerrar el navegador si está inactivo.
func (brm *BrowserManager) startIdleTimer() {
	// Si el temporizador ya existe, lo reseteamos
	if brm.idleTimer != nil {
		brm.idleTimer.Stop()
	}
	// Iniciar un nuevo temporizador de 60 segundos
	brm.idleTimer = time.AfterFunc(120*time.Second, func() {
		fmt.Println("Inactividad detectada. Cerrando navegador.")
		brm.Close() // Cierra el navegador
		// brm.KillChromeProcesses() // Mata los procesos de Chrome
	})
}
