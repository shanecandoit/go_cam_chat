package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	webview "github.com/webview/webview_go"
)

var html_start_record = ""

func main() {

	// create a localhost server
	http.HandleFunc("/", handler)
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

func today_as_string() string {
	t := time.Now()
	// want yyyymmdd_hhmmss_mmmm
	millisecs := t.Nanosecond() / 1000000
	millsecs_str := fmt.Sprintf("%04d", millisecs)
	return t.Format("20060102_150405") + "_" + millsecs_str
}

// handl POST requests
func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// if the request method is POST
		// read the body of the request
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}

		length := len(body)

		// print the body to the console
		length_str := fmt.Sprintf("%d", length)
		fmt.Println("uploading", string(length_str))

		first_200_bytes := body[:200]
		fmt.Println("first_200_bytes")
		fmt.Println(string(first_200_bytes))
		fmt.Println()

		first_200_bytes_str := string(first_200_bytes)
		first_200_bytes_string_lines := strings.Split(first_200_bytes_str, "\n")

		// Content-Disposition: form-data; name="video"; filename="video_560055.webm"
		// Content-Type: video/webm; codecs="vp9,opus"
		content_disposition := ""
		content_type := ""
		file_name := ""
		for _, line := range first_200_bytes_string_lines {
			if strings.Contains(line, "Content-Disposition") && strings.Contains(line, ": ") {
				content_disposition = string(line)
				content_disposition = strings.Split(line, ": ")[1]

				// get the file name
				file_name = strings.Split(content_disposition, "filename=")[1]
				file_name = strings.Replace(file_name, "\"", "", -1)
				// strip whitespace
				file_name = strings.Trim(file_name, " \t\n\r")
			}
			if strings.Contains(line, "Content-Type") && strings.Contains(line, ": ") {
				content_type = string(line)
				content_type = strings.Split(line, ": ")[1]
			}
		}

		if content_disposition == "" {
			// we won't have the file name, so bail out
			fmt.Println("Content-Disposition not found")
			return
		}

		// we should have the file name and the content type
		fmt.Println("file_name", file_name)
		fmt.Println("content_type", content_type)

		// save the video file from the form data
		video_bytes_index := strings.Index(string(body), "\r\n\r\n") + 4
		// fmt.Println("video_bytes_index", video_bytes_index)
		// video_bytes_index 165
		video_bytes := body[video_bytes_index:]
		// fmt.Println("video_bytes", len(video_bytes))

		// print the first 10 bytes of the video file as hex
		// 1A 45 DF A3
		// fmt.Printf("video_bytes[0] %x \n", video_bytes[0])
		// fmt.Printf("video_bytes[1] %x \n", video_bytes[1])
		// fmt.Printf("video_bytes[2] %x \n", video_bytes[2])
		// fmt.Printf("video_bytes[3] %x \n", video_bytes[3])

		is_webm := false
		if video_bytes[0] == 0x1A && video_bytes[1] == 0x45 && video_bytes[2] == 0xDF && video_bytes[3] == 0xA3 {
			is_webm = true
		}
		fmt.Println("is_webm", is_webm)
		if !is_webm {
			// fmt.Println("Not a webm file")
			return
		}

		// save the video file
		cwd := os.Getenv("PWD")
		fmt.Println("cwd", cwd)
		// if "vids" directory exists, save the video file there
		vids_exists, err := os.Stat("vids")
		if err != nil {
			fmt.Println("Error checking if vids directory exists", err)
			return
		}

		// change name from vids/video_329390.webm
		// to vids/video_329390_some_date.webm
		date_str := today_as_string()
		fmt.Println("file_name before", file_name)
		file_name = strings.Replace(file_name, ".webm", "_"+date_str+".webm", -1)

		fmt.Println("file_name  after", file_name)
		// file_name before video_875991.webm
		// file_name  after video_875991_2024-11-22_17:52:05.3513.webm

		if vids_exists.IsDir() {
			fmt.Println("vids directory exists")
			file_name = "vids/" + file_name
		}
		fmt.Println("file_name", file_name)
		err = os.WriteFile(file_name, video_bytes, 0644)
		if err != nil {
			fmt.Println("Error saving video file", err)
			return
		}

	} else if r.Method == "GET" { // if the request method is GET

		// send favicon
		if r.URL.Path == "/favicon.ico" {
			http.ServeFile(w, r, "favicon.ico")
			return
		}

		// print the request URL to the console
		fmt.Println(r.URL)

		html_start_rec, err := os.ReadFile("record.html")
		if err != nil {
			fmt.Fprintf(w, "Error reading file: %v", err)
			return
		}
		html_start_record = string(html_start_rec)
		// return the html_start_record to the client
		fmt.Fprint(w, html_start_record)
	}
}
