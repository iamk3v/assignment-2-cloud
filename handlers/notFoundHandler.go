package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	_, err := fmt.Fprintf(w, "404 - Not Found, only a large amount of dust found\n"+
		strconv.Itoa(http.StatusNotFound))
	if err != nil {
		log.Print("Error occurred when trying to send response: " + err.Error())
		http.Error(w, "An internal error occurred..", http.StatusInternalServerError)
		return
	}
}
