package main

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type PageData struct {
	Text   string
	Result string
	Banner string
	Error  string
}

// render fills the template and sends it with the given status code.
// Every page (normal or error) goes through here, so they all look the same.
func render(w http.ResponseWriter, status int, data PageData) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "404 template not found", http.StatusNotFound)
		return
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		log.Println("error failed in the execution", err)
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	buf.WriteTo(w)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		render(w, http.StatusNotFound, PageData{Error: "That page doesn't exist. You can start drawing here instead."})
		return
	}
	render(w, http.StatusOK, PageData{})
}

// buildArt reads the form, checks it, and draws the art.
// Both /ascii-art and /download use it.
func buildArt(r *http.Request) (PageData, int) {
	text := r.FormValue("theInput")
	text = strings.ReplaceAll(text, "\r\n", "\n") // Enter in a textarea sends \r\n
	text = strings.ReplaceAll(text, "\\n", "\n")  // still allow typing \n by hand
	banner := r.FormValue("banner")
	data := PageData{Text: text, Banner: banner}

	for _, char := range text {
		if char != '\n' && (char < 32 || char > 126) {
			data.Error = "\"" + string(char) + "\" can't be drawn. Use letters, numbers, spaces and keyboard symbols only."
			return data, http.StatusBadRequest
		}
	}

	if banner != "standard.txt" && banner != "shadow.txt" && banner != "thinkertoy.txt" {
		data.Banner = ""
		data.Error = "Please choose one of the three fonts."
		return data, http.StatusBadRequest
	}

	ContentFromTheBanner, err := os.ReadFile(banner)
	if err != nil {
		data.Error = "The " + strings.TrimSuffix(banner, ".txt") + " font file is missing on the server."
		return data, http.StatusNotFound
	}

	cleaned := strings.ReplaceAll(string(ContentFromTheBanner), "\r\n", "\n")
	AdjustedContent := strings.Split(cleaned, "\n")

	data.Result = drawing(text, AdjustedContent)
	return data, http.StatusOK
}

func asciiHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		render(w, http.StatusMethodNotAllowed, PageData{Error: "Use the form below to draw your art."})
		return
	}

	data, status := buildArt(r)
	render(w, status, data)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		render(w, http.StatusMethodNotAllowed, PageData{Error: "Use the Download button to save your art."})
		return
	}

	data, status := buildArt(r)
	if status != http.StatusOK {
		render(w, status, data)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Content-Disposition", `attachment; filename="ascii-art.txt"`)
	w.Header().Set("Content-Length", strconv.Itoa(len(data.Result)))
	w.Write([]byte(data.Result))
}

func startLine(TheRune rune) int {
	return (int(TheRune)-32)*9 + 1
}

func drawing(text string, lines []string) string {
	art := ""
	words := strings.Split(text, "\n")
	if strings.ReplaceAll(text, "\n", "") == "" {
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

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ascii-art", asciiHandler)
	http.HandleFunc("/download", downloadHandler)

	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
