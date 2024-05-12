package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Nurs005/gql-yt/handlers"
	"github.com/gin-gonic/gin"
)

const defaultPort = "8081"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	r := gin.Default()
	r.GET("/", handlers.HandlePlyground())
	r.POST("/query", handlers.HandleGQL())

	corsHandler := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			h.ServeHTTP(w, r)
		})
	}

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, corsHandler(r)))
}
