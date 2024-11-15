package main

import (
	"fmt"
	"net/http"

	webview "github.com/webview/webview_go"
)

var html_start_record = `
<!DOCTYPE html>
<html>
  <head>
    <title>Webcam Recorder</title>
    <style>
      #recordButton {
        background-color: red;
        border: none;
        color: white;
        padding: 15px 32px;
        text-align: center;
        text-decoration: none;
        display: inline-block;
        font-size: 16px;
        margin: 4px 2px;
        cursor: pointer;
      }
    </style>
  </head>
  <body>
    <h1>Webcam Recorder</h1>
    
    <video id="video" width="640" height="480" autoplay></video>
    <br>
    <button id="recordButton">Record</button>
    <br>

    <script>
      const video = document.getElementById("video");
      const recordButton = document.getElementById("recordButton");
      let recorder;
      let isRecording = false;

      navigator.mediaDevices
          .getUserMedia({ video: true, audio: true })
          .then(function (stream) {
              video.srcObject = stream;
              video.muted = true; // Mute the video element to avoid feedback noise.
              recorder = new MediaRecorder(stream);

              recorder.ondataavailable 
= (event) => {
                  // Handle recorded data (e.g., save to file or send to server)
                  console.log(event.data);

                  const blob = new Blob([event.data], { type: 'video/webm; codecs="vp9,opus"' });
                  const url = URL.createObjectURL(blob);

                  var random_number = Math.floor(Math.random() * 1000000);
                  var random_string = 'video_' + random_number.toString();
                  random_string = random_string.replace('.', '') + '.webm';

                  console.log(random_string);
                  console.log(url);

                  const a = document.createElement('a');
                  a.href = url;
                  a.download = random_string;
                  // a.click();
                  // URL.revokeObjectURL(url); // free up storage--no longer needed.

                  // save the video to disk
                  const link = document.createElement('a');
                  link.href = url;
                  link.setAttribute('download', random_string);
                  link.setAttribute('target', '_blank'); // open in new tab
                  link.innerText = 'Download - ' + random_string; // or display the name of the file
                  document.body.appendChild(link);
                  link.click();
                  // document.body.removeChild(link);
                  document.appendChild(document.createElement('br'));

              };

              recorder.onerror = (error) => {
                  console.error("Error recording:", error);
              };
          })
          .catch(function (error) {
              console.error("Error accessing media devices.", error);
          });

      recordButton.addEventListener("click", () => {
          if (!isRecording) {
              recorder.start();
              recordButton.textContent = "Stop Recording";
              isRecording = true;
          } else {
              recorder.stop();
              recordButton.textContent = "Start Recording";
              isRecording = false;
          }
      });
  </script>
  </body>
</html>

`

func main() {

	// create a localhost server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, html_start_record)
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

// we have to serve the page over https
