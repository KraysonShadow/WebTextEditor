package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

const filePath = "notepad_content.txt"

func main() {
	// Обработчик для корневого пути (HTML-страница)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		http.ServeFile(w, r, "editor.html")
	})

	// Обработчик для сохранения содержимого
	http.HandleFunc("/save", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var content map[string]string
			if err := json.NewDecoder(r.Body).Decode(&content); err != nil {
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}
			if err := os.WriteFile(filePath, []byte(content["content"]), 0644); err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			w.Write([]byte("Содержимое сохранено"))
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	// Обработчик для загрузки содержимого
	http.HandleFunc("/load", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			content, err := os.ReadFile(filePath)
			if err != nil {
				if os.IsNotExist(err) {
					content = []byte("{}") // Возвращаем пустой JSON, если файл не существует
				} else {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(content)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	// Запуск сервера на порту 8080
	log.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
