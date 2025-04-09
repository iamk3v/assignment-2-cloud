package handlers

import (
	"fmt"
	"log"
	"net/http"
)

/*
NotFoundHandler Handles all requests to endpoints that does not yet exist.
*/
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	_, err := fmt.Fprintf(w, "%d - %s, only a large amount of dust found, double check your URL!",
		http.StatusNotFound, http.StatusText(http.StatusNotFound))

	if err != nil {
		log.Print("Error occurred when trying to send response: " + err.Error())
		http.Error(w, "An internal error occurred..", http.StatusInternalServerError)
		return
	}
}
