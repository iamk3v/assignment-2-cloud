package assignment_2

import (
	"log"
	"net/http"
	"os"
)

func main() {

	// Create a new router
	router := http.NewServeMux()

	// Routes
	router.HandleFunc()

	// Define port
	PORT := "8080"
	if os.Getenv("PORT") != "" {
		PORT = os.Getenv("PORT")
	}

	// Listen on the designated port for traffic
	log.Println("Starting server on port " + PORT)
	err := http.ListenAndServe(":"+PORT, router)
	if err != nil {
		log.Fatal(err.Error())
	}
}
