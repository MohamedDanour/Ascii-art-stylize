package main

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
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
