package handlers

import "net/http"

func RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleRegGetRequest(w, r)
	case http.MethodPost:
		handleRegPostRequest(w, r)
	case http.MethodDelete:
		handleRegDeleteRequest(w, r)
	case http.MethodPut:
		handleRegPutRequest(w, r)
	default:
		http.Error(w, "REST method '"+r.Method+"' not supported. "+
			"Currently only '"+http.MethodGet+"' is supported.", http.StatusNotImplemented)
		return
	}
}

func handleRegGetRequest(w http.ResponseWriter, r *http.Request) {

}

func handleRegPostRequest(w http.ResponseWriter, r *http.Request) {

}

func handleRegDeleteRequest(w http.ResponseWriter, r *http.Request) {

}

func handleRegPutRequest(w http.ResponseWriter, r *http.Request) {

}
