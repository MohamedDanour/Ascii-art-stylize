package main

import (
	"net/http"
	"os"
	"strings"
)

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
