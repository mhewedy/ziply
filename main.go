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
		<style>
		.text {
			padding: .5rem 1rem;
			font-size: 1.25rem;
			line-height: 1.5;
			border-radius: .3rem;
			display: inline-block;
			width: 85%;
			color: #495057;
			background-color: #fff;
			background-clip: padding-box;
			border: 1px solid #ced4da;
		}
		.button {
			cursor: pointer;
			color: #fff;
			background-color: #007bff;
			border-color: #007bff;
		  	display: inline-block;
			font-weight: 400;
			text-align: center;
			white-space: nowrap;
			vertical-align: middle;
			-webkit-user-select: none;
			-moz-user-select: none;
			-ms-user-select: none;
			user-select: none;
			padding: .7rem 1.5rem;
			font-size: 1rem;
			line-height: 1.5;
			border-radius: .25rem;
    		margin-top: -8px;
		}
		</style>
		<div style="margin: 15% 1% 25% 1%;">
			<h2 style="font-family: sans-serif">Enter HTTP URL to zip and download</h2>
			<form action="/download">
				<div>
					<span><input type="text" name="url" placeholder="Enter HTTP URL" class="text" /></span>
					<span><input type="submit" value="Zip and Download" class="button" /></span>
				</div>
			</form>
		</div>
		<div style="text-align: center; font-family: sans-serif; font-size: small">
		<a href="https://github.com/mhewedy/ziply" target="_blank">
			<img alt="github" src="https://github.githubassets.com/images/modules/logos_page/GitHub-Mark.png" width="30px">
		</a>
		</div>
		`

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintln(w, homeHtml)
	})

	http.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {

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
