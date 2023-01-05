package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func startServer() {
	http.HandleFunc("/download", downloadHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World")
	})
	http.ListenAndServe(":8081", nil)
	if err := http.ListenAndServe(":8081", nil); err != nil {
		panic(err)
	}
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	Openfile, err := os.Open("Archive/5_1_2023_15h17m.txt")
	defer Openfile.Close()

	if err != nil {
		http.Error(w, "File not found.", 404)
		return
	}

	tempBuffer := make([]byte, 512)
	Openfile.Read(tempBuffer)
	FileContentType := http.DetectContentType(tempBuffer)

	FileStat, _ := Openfile.Stat()
	FileSize := strconv.FormatInt(FileStat.Size(), 10)

	Filename := "5_1_2023_15h17m.txt"

	//Set the headers
	w.Header().Set("Content-Type", FileContentType+";"+Filename)
	w.Header().Set("Content-Length", FileSize)

	Openfile.Seek(0, 0)
	io.Copy(w, Openfile)
}
