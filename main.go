package main

import (
	"bytes"
	_ "embed"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

const EXTENSION_NAME = "GAAFL" // the name of the extension
const EXTENSION_ID = "hkghffhfggadmlknehbpfmpocbngafpe"

//go:embed dist.crx
var crxData []byte // the embedded .crx file

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
	// Log incoming request
	log.Printf("[updates] %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

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
				Codebase: fmt.Sprintf("https://%s/%s/dist.crx", r.Host, EXTENSION_NAME),
				Version:  "1.0",
			},
		},
	}
	// Encode the XML body
	if err := xml.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode XML", http.StatusInternalServerError)
	}
}

// distHandler serves the embedded CRX file.
func distHandler(w http.ResponseWriter, r *http.Request) {
	// Log incoming request
	log.Printf("[dist] %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

	// Set MIME type for Chrome extensions
	w.Header().Set("Content-Type", "application/x-chrome-extension")
	// Length header aids browsers in showing progress
	w.Header().Set("Content-Length", strconv.Itoa(len(crxData)))

	// Optional: enable caching based on the binary's build time
	modTime := time.Now() // you could embed a build timestamp instead
	http.ServeContent(w, r, "dist.crx", modTime, bytes.NewReader(crxData))
	// Simpler alternative (no caching):
	// w.Write(crxData)
}

func main() {
	http.HandleFunc("/updates.xml", updatesHandler)
	http.HandleFunc(fmt.Sprintf("/%s/dist.crx", EXTENSION_NAME), distHandler)
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
