package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
)

const EXTENSION_NAME = "GAAFL" // the name of the extension
const EXTENSION_ID = "hkghffhfggadmlknehbpfmpocbngafpe"

type GUpdate struct {
	XMLName  xml.Name `xml:"gupdate"`
	Xmlns    string   `xml:"xmlns,attr"`
	Protocol string   `xml:"protocol,attr"`
	App      App      `xml:"app"`
}

type App struct {
	ID          string      `xml:"appid,attr"`
	UpdateCheck UpdateCheck `xml:"updatecheck"`
}

type UpdateCheck struct {
	Codebase string `xml:"codebase,attr"` // url to the crx file
	Version  string `xml:"version,attr"`
}

func updatesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml")
	// Write XML declaration
	if _, err := w.Write([]byte(xml.Header)); err != nil {
		http.Error(w, "Failed to write XML header", http.StatusInternalServerError)
		return
	}
	// Prepare the response structure
	resp := GUpdate{
		Xmlns:    "http://www.google.com/update2/response",
		Protocol: "2.0",
		App: App{
			ID: EXTENSION_ID,
			UpdateCheck: UpdateCheck{
				Codebase: fmt.Sprintf("http://%s/%s/dist.crx", r.Host, EXTENSION_NAME),
				Version:  "2.0",
			},
		},
	}
	// Encode the XML body
	if err := xml.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode XML", http.StatusInternalServerError)
	}
}

func distHandler(w http.ResponseWriter, r *http.Request) {
	// Return a simple plain‑text response
	w.Header().Set("Content-Type", "application/x-chrome-extension")
	// fmt.Fprint(w, "Hello, this is the dist handler returning text.")
	// Serve the extension .crx file
	// Adjust the path as necessary for your project layout
	filePath := "./dist.crx"
	http.ServeFile(w, r, filePath)
}

func main() {
	http.HandleFunc("/updates.xml", updatesHandler)
	http.HandleFunc(fmt.Sprintf("/%s/dist.crx", EXTENSION_NAME), distHandler)
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
