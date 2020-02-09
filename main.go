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
	"strconv"
	"strings"
)

var homeHtml = `
<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css">
<h1 style="text-align: center; margin-top: 5%;">Ziply: Zip&Download URL content</h1>
<div style="margin: 10% 1% 20% 1%;">
	<form action="/dl">
		<div>
			<input type="text" name="url" placeholder="Enter HTTP URL to zip and download" class="form-control form-control-lg" required />
			<div style="width: 20%; margin: 20px 0 0 40%;">
				<input type="submit" value="Zip and Download" class="btn btn-primary btn-lg btn-block"/>
			</div>
		</div>
	</form>
</div>
<div style="text-align: center;">
<a href="https://github.com/mhewedy/ziply" target="_blank">
	<img alt="github" src="https://github.githubassets.com/images/modules/logos_page/GitHub-Mark.png" width="30px">
</a>
</div>`

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintln(w, homeHtml)
	})

	http.HandleFunc("/dl", func(w http.ResponseWriter, r *http.Request) {

		urls := r.URL.Query()["url"]
		if len(urls) == 0 || strings.TrimSpace(urls[0]) == "" {
			sendError(w, "missing url parameter", http.StatusBadRequest)
			return
		}

		zipped, err := downloadAndZip(urls[0])
		if err != nil {
			sendError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fileName := path.Base(urls[0]) + ".zip"
		fmt.Printf("downloading %s\n", fileName)
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
		w.Header().Set("Content-Length", strconv.Itoa(zipped.Len()))
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
	w.WriteHeader(code)
	w.Write([]byte(msg))
}

func downloadAndZip(exeURL string) (*bytes.Buffer, error) {

	fmt.Printf("getting %s...\n", exeURL)

	resp, err := http.Get(exeURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, resp.Body)
	if err != nil {
		return nil, err
	}

	return doZip(&buf, path.Base(exeURL))
}

func doZip(in *bytes.Buffer, name string) (*bytes.Buffer, error) {

	fmt.Printf("zipping %s...\n", name)

	var buf bytes.Buffer
	w := zip.NewWriter(&buf)

	f, err := w.Create(name)
	if err != nil {
		return nil, err
	}

	_, err = f.Write(in.Bytes())
	if err != nil {
		return nil, err
	}

	err = w.Close()
	if err != nil {
		return nil, err
	}

	return &buf, nil
}
