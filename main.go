package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func main() {
	port := os.Getenv("PORT")
	listenAddr := os.Getenv("LISTEN_ADDR")

	if port == "" {
		port = "8000"
	}

	if listenAddr == "" {
		listenAddr = "localhost"
	}

	addr := listenAddr + `:` + port

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Hello World1!")
	})

	http.HandleFunc("/audio", audioHandler)

	http.HandleFunc("/playlist", playlistHandler)

	log.Fatal(http.ListenAndServe(addr, nil))
}

func audioHandler(w http.ResponseWriter, r *http.Request) {
	// Video Id
	v := r.URL.Query().Get("v")

	if v == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "use format /audio?v=...")
		return
	}

	url := fmt.Sprintf("https://www.youtube.com/watch?v=%s", v)

	ytdl := exec.Command("youtube-dl", "--extract-audio", "--audio-format", "mp3", url)
	output, err := ytdl.Output()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.mp3", v))

	w.Header().Set("Content-Type", "audio/mpeg")

	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(output)))

	if _, err := w.Write(output); err != nil {
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

}

func playlistHandler(w http.ResponseWriter, r *http.Request) {
	// Playlist Id
	p := r.URL.Query().Get("p")

	// Download Format
	f := r.URL.Query().Get("f")

	if p == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "use format /playlist?p=...&f=mp3|mp4")
	}

	if f == "" {
		f = "mp3"
	}

	url := fmt.Sprint("https://www.youtube.com/playlist?list=" + p)

	ytdl := exec.Command("youtube-dl", "-cit", "--extract-audio", "--audio-format", "mp3", url)

	output, err := ytdl.Output()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s_playlist.%s\"", p, f))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(output)))

	if _, err := w.Write(output); err != nil {
		fmt.Fprintf(w, "Error: %s", err)
		return
	}
}
