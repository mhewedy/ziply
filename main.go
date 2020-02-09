package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"path"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		url := r.URL.Query()["url"][0]

		zipped, _ := downloadAndZip(url)
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, path.Base(url)+".zip"))
		_, err := w.Write(zipped.Bytes())

		if err != nil {
			_, err = w.Write([]byte(err.Error()))
			if err != nil {
				log.Println(err)
			}
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))

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
