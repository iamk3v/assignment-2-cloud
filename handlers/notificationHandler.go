package handlers

import (
	"assignment-2/config"
	"assignment-2/database"
	"assignment-2/utils"
	"context"
	"encoding/json"
	firebase "firebase.google.com/go"
	"fmt"
	"log"
	"net/http"
	"strings"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

var (
	notiClient      *firestore.Client
	notiFirebaseCtx context.Context
)

func init() {
	notiFirebaseCtx = context.Background()
	sa := option.WithCredentialsFile("../config/service-account.json")
	app, err := firebase.NewApp(notiFirebaseCtx, nil, sa)
	if err != nil {
		log.Fatalf("NotificationHandler init: firebase.NewApp: %v", err)
	}
	notiClient, err = app.Firestore(notiFirebaseCtx)
	if err != nil {
		log.Fatalf("NotificationHandler init: app.Firestore: %v", err)
	}
}

func NotificationHandler(w http.ResponseWriter, r *http.Request) {
	// /dashboard/v1/notifications/{id} or /dashboard/v1/notifications/
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, config.START_URL+"/notifications"), "/")

	if len(pathParts) > 1 && pathParts[1] != "" {
		id := pathParts[1]
		switch r.Method {
		case http.MethodGet:
			handleNotiGetOneRequest(w, r, id)
		case http.MethodDelete:
			handleNotiDeleteRequest(w, r, id)
		default:
			http.Error(w,
				fmt.Sprintf("Method %s not supported on /notifications/{id}", r.Method),
				http.StatusMethodNotAllowed)
		}
	} else {
		// collection
		switch r.Method {
		case http.MethodGet:
			handleNotiGetAllRequest(w, r)
		case http.MethodPost:
			handleNotiPostRequest(w, r)
		default:
			http.Error(w,
				fmt.Sprintf("Method %s not supported on /notifications/", r.Method),
				http.StatusMethodNotAllowed)
		}
	}
}

func handleNotiGetAllRequest(w http.ResponseWriter, r *http.Request) {
	hooks, err := database.GetAllWebhooks(notiFirebaseCtx, notiClient)
	if err != nil {
		log.Println("Error retrieving webhooks:", err)
		http.Error(w, config.ERR_INTERNAL_SERVER_ERROR, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hooks)
}

func handleNotiGetOneRequest(w http.ResponseWriter, r *http.Request, id string) {
	hook, err := database.GetWebhook(notiFirebaseCtx, notiClient, id)
	if err != nil {
		log.Println("Error retrieving webhook:", err)
		http.Error(w, config.ERR_NOT_FOUND, http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hook)
}

func handleNotiPostRequest(w http.ResponseWriter, r *http.Request) {
	var hook utils.Webhook
	if err := json.NewDecoder(r.Body).Decode(&hook); err != nil {
		log.Println("Error decoding webhook body:", err)
		http.Error(w, config.ERR_BAD_REQUEST, http.StatusBadRequest)
		return
	}
	id, err := database.CreateWebhook(notiFirebaseCtx, notiClient, hook)
	if err != nil {
		log.Println("Error creating webhook:", err)
		http.Error(w, config.ERR_INTERNAL_SERVER_ERROR, http.StatusInternalServerError)
		return
	}
	resp := map[string]string{"id": id}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func handleNotiDeleteRequest(w http.ResponseWriter, r *http.Request, id string) {
	err := database.DeleteWebhook(notiFirebaseCtx, notiClient, id)
	if err != nil {
		log.Println("Error deleting webhook:", err)
		http.Error(w, config.ERR_NOT_FOUND, http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
