package main

import (
	"fmt"
	"net/http"
	"os"

	webview "github.com/webview/webview_go"
)

var html_start_record = ""

func main() {

	// create a localhost server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// read html_start_record from file record.html
		// get it fresh every time, so we can change it without restarting the server
		html_start_rec, err := os.ReadFile("record.html")
		if err != nil {
			fmt.Fprintf(w, "Error reading file: %v", err)
			return
		}
		html_start_record = string(html_start_rec)

		fmt.Fprint(w, html_start_record)
	})
	go http.ListenAndServe(":8080", nil)

	// create a webview
	debug := true
	w := webview.New(debug)
	defer w.Destroy()
	w.SetTitle("Basic Webview")
	w.SetSize(600, 400, webview.HintNone)
	// w.SetHtml(html_start_record)
	w.Navigate("http://localhost:8080")
	w.Run()
}

// upload and create separate audio file?
