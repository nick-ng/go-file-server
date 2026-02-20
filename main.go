package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func fileDownloadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

	if r.Method == http.MethodPost {
		w.Write([]byte("Closing"))
		os.Exit(0)
	}

	filePath := r.URL.Path[len("/"):]
	if len(filePath) == 0 {
		w.Write([]byte(`
		<html>
		<body>
			<h1>File Server</h1>
			<p>Usage is http://localhost:5911/&lt;filename&gt;</p>
		<form method="post">
		<button>Exit</button>
		</body>
		</form>
		</html>
		`))
		return
	}

	if strings.ContainsAny(filePath, `/\?%*:|"<>,;= `) {
		http.Error(w, "File not found.", 404)
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
