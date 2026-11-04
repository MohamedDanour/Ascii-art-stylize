package main

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Error status not found", http.StatusNotFound)
		return
	}

	err = tmpl.Execute(w, PageData{})
	if err != nil {
		log.Println("error failed in the execution", err)
	}

}

func asciiHandler(w http.ResponseWriter, r *http.Request) {
	text := r.FormValue("theInput")

	for _, char := range text {
		if !(char >= 32 && char <= 126) {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
	}

	banner := r.FormValue("banner")
	if banner != "standard.txt" && banner != "shadow.txt" && banner != "thinkertoy.txt" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	ContentFromTheBanner, err := os.ReadFile(banner)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	AdjustedContent := strings.Split(string(ContentFromTheBanner), "\n")

	art := drawing(text, AdjustedContent)

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, PageData{Result: art, Banner: banner})
	if err != nil {
		log.Println("error failed in the execution", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	buf.WriteTo(w)
}

func startLine(TheRune rune) int {
	return (int(TheRune)-32)*9 + 1
}

func drawing(text string, lines []string) string {
	art := ""
	words := strings.Split(text, "\\n")
	if strings.ReplaceAll(text, "\\n", "") == "" {
		return strings.Repeat("\n", len(words)-1)
	}
	for _, word := range words {
		if word == "" {
			art += "\n"
			continue
		}
		for row := 0; row < 8; row++ {
			for _, letter := range word {
				art += lines[startLine(letter)+row]
			}
			art += "\n"
		}
	}
	return art
}

type PageData struct {
	Result string
	Banner string
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ascii-art", asciiHandler)
	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
