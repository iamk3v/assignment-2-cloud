package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

/*
NotFoundHandler Handles all requests to endpoints that does not yet exist.
*/
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("html/index.html")
	if err != nil {
		log.Print("Error occurred when trying to open file: " + err.Error())
		_, err := fmt.Fprintf(w, "%d - %s, only a large amount of dust found, double check your URL!",
			http.StatusNotFound, http.StatusText(http.StatusNotFound))
		if err != nil {
			log.Print("Error occurred when trying to send response: " + err.Error())
			http.Error(w, "An internal error occurred..", http.StatusInternalServerError)
			return
		}
	}

	defer file.Close()

	w.WriteHeader(http.StatusNotFound)
	_, err = io.Copy(w, file)
	if err != nil {
		log.Print("Error occurred when trying to copy file: " + err.Error())
		_, err := fmt.Fprintf(w, "%d - %s, only a large amount of dust found, double check your URL!",
			http.StatusNotFound, http.StatusText(http.StatusNotFound))
		if err != nil {
			log.Print("Error occurred when trying to send response: " + err.Error())
			http.Error(w, "An internal error occurred..", http.StatusInternalServerError)
			return
		}
	}
}
