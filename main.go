package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func fileDownloadHandler(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path[len("/"):] // Assumes file path is after /download/
	if len(filePath) == 0 {
		http.Error(w, "no file provided. usage is http://localhost:5911/<filename>", 400)
		return
	}
	_, fileName := filepath.Split(filePath)
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "File not found.", 404)
		return
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		http.Error(w, "Internal server error.", 500)
		return
	}

	// Set headers
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	w.Header().Set("Content-Length", strconv.FormatInt(fileStat.Size(), 10))

	// Stream the file
	http.ServeContent(w, r, fileName, fileStat.ModTime(), file)
}

func main() {
	http.HandleFunc("/", fileDownloadHandler)
	fmt.Println("Starting server on :5911")
	if err := http.ListenAndServe(":5911", nil); err != nil {
		log.Fatal(err)
	}
}
