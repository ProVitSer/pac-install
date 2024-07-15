package main

import (
	"fmt"
	"os"
	"pac-install/internal/web"
	"runtime"
)

func main() {
	// Определение ОС
	fmt.Println("Определение операционной системы...")

	switch client_os := runtime.GOOS; client_os {
	case "linux":
		fmt.Println("Это Linux.")
	case "windows":
		fmt.Println("Это Windows.")
	default:
		fmt.Printf("Неизвестная операционная система: %s\n", client_os)
		os.Exit(1)

	}

	// Запуск веб-сервера для ввода лицензионного ключа
	go web.StartServer()

	// Открытие браузера
	fmt.Println("Открытие браузера для ввода лицензионного ключа...")
	web.OpenBrowser("http://localhost:8080/license")

	// Ожидание завершения веб-сервера
	select {}
}
