package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		urls := r.URL.Query()["url"]
		if len(urls) == 0 {
			sendError(w, "missing url parameter", http.StatusBadRequest)
			return
		}

		zipped, err := downloadAndZip(urls[0])
		if err != nil {
			sendError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, path.Base(urls[0])+".zip"))
		_, err = w.Write(zipped.Bytes())

		if err != nil {
			sendError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("started...")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func sendError(w http.ResponseWriter, msg string, code int) {
	_, _ = w.Write([]byte(msg))
	w.WriteHeader(code)
}

func downloadAndZip(exeURL string) (*bytes.Buffer, error) {
	fmt.Printf("Getting %s...\n", exeURL)

	resp, err := http.Get(exeURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var buff bytes.Buffer

	_, err = io.Copy(&buff, resp.Body)
	if err != nil {
		return nil, err
	}

	return zipBuffer(&buff, path.Base(exeURL))
}

func zipBuffer(in *bytes.Buffer, name string) (*bytes.Buffer, error) {

	fmt.Printf("Zipping %s...\n", name)

	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)

	// Create a new zip archive.
	w := zip.NewWriter(buf)

	f, err := w.Create(name)
	if err != nil {
		log.Fatal(err)
	}
	_, err = f.Write(in.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	// Make sure to check the error on Close.
	err = w.Close()
	if err != nil {
		log.Fatal(err)
	}

	return buf, nil
}
