package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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
	// Open the file
	f, err := os.Open("./Archive/5_1_2023_12h_31m.txt")
	if err != nil {
		http.Error(w, "File not found", 404)
		return
	}
	defer f.Close()

	// Set the headers
	w.Header().Set("Content-Disposition", "attachment; filename=5_1_2023_12h_31m.txt")
	w.Header().Set("Content-Type", "text/plain")

	// Copy the file contents to the response writer
	_, err = io.Copy(w, f)
	if err != nil {
		log.Println(err)
	}
}
