package handlers

import "net/http"

func NotificationHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleNotiGetRequest(w, r)
	case http.MethodPost:
		handleNotiPostRequest(w, r)
	case http.MethodDelete:
		handleNotiDeleteRequest(w, r)
	default:
		http.Error(w, "REST method '"+r.Method+"' not supported. "+
			"Currently only '"+http.MethodGet+"' is supported.", http.StatusNotImplemented)
		return
	}
}

func handleNotiGetRequest(w http.ResponseWriter, r *http.Request) {

}

func handleNotiPostRequest(w http.ResponseWriter, r *http.Request) {

}

func handleNotiDeleteRequest(w http.ResponseWriter, r *http.Request) {

}
