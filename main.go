package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"pawmap/database"
	"pawmap/server"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Initialize the database
	database.InitDB()

	// Init log folder
	if _, err := os.Stat("./logs"); os.IsNotExist(err) {
		os.Mkdir("./logs", 0755)
	}

	// Init Logger
	f, err := os.OpenFile("./logs/api.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(wrt)

	// Initialize the router
	r := mux.NewRouter()
	v1 := r.PathPrefix("/api/v1").Subrouter()

	// Pins
	v1.HandleFunc("/pin/submitPin", server.SubmitPinHandler).Methods("POST")
	v1.HandleFunc("/pin/getAllPins", server.GetAllPinsHandler).Methods("GET")
	v1.HandleFunc("/pin/getPinsByLocation", server.GetPinsByLocationHandler).Methods("GET")

	// Users
	v1.HandleFunc("/user/createUser", server.CreateUserHandler).Methods("POST")
	v1.HandleFunc("/user/getUserByUsername", server.GetUserByUsernameHandler).Methods("GET")
	v1.HandleFunc("/user/getUserByEmail", server.GetUserByMailHandler).Methods("GET")

	// Sessions
	v1.HandleFunc("/session/createSession", server.CreateSessionHandler).Methods("POST")
	v1.HandleFunc("/session/getSession", server.GetSessionHandler).Methods("GET")
	v1.HandleFunc("/session/checkSession", server.CheckSessionHandler).Methods("GET")
	v1.HandleFunc("/session/deleteSession", server.DeleteSessionHandler).Methods("DELETE")

	//allowedOrigin := os.Getenv("ALLOWED_ORIGIN")

	// CORS middleware
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := corsMiddleware.Handler(r)

	// Start the server
	fmt.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe("localhost:8080", handler))
}
