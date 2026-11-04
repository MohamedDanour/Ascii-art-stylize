package main

import (
	"net/http"
	"strconv"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		render(w, http.StatusNotFound, PageData{Error: "That page doesn't exist. You can start drawing here instead."})
		return
	}
	render(w, http.StatusOK, PageData{})
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
