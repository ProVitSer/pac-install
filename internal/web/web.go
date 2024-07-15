package web

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"os/exec"
	"runtime"
	"time"
)

func StartServer() {
	http.HandleFunc("/license", licenseHandler)

	http.HandleFunc("/ws", wsHandler)

	http.ListenAndServe(":8080", nil)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func licenseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		_ = r.FormValue("license_key")
		//valid := license.ValidateLicenseKey(licenseKey)
		//if !valid {
		//	http.Error(w, "Invalid license key", http.StatusUnauthorized)
		//	return
		//}

		fmt.Fprintln(w, "Активация прошла успешно! Установка продолжается...")
		return
	}

	http.ServeFile(w, r, "web/index.html")
}

func wsHandler(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Не удалось настроить обновление ws:", err)
		return
	}

	defer conn.Close()

	if err := processInstallation(conn); err != nil {
		sendProgress(conn, "", false, err.Error())
	}
}

func processInstallation(conn *websocket.Conn) error {
	stages := []string{"download", "installComponents", "finishInstallation"}

	for _, stage := range stages {
		if err := sendProgress(conn, stage, false, "Ошибка установки"); err != nil {
			return err
		}
		time.Sleep(2 * time.Second)
	}
	return nil
}

func sendProgress(conn *websocket.Conn, stage string, success bool, errorMsg string) error {

	message := map[string]interface{}{
		"stage":   stage,
		"success": success,
		"error":   errorMsg,
	}

	err := conn.WriteJSON(message)
	if err != nil {
		fmt.Println("Ошибка обмена данными:", err)
	}

	return err
}

func OpenBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	default:
		err = fmt.Errorf("Не поддерживаемая платформа")
	}

	if err != nil {
		fmt.Println("Не удалось открыть браузер: ", err)
	}
}
