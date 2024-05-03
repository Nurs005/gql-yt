package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Nurs005/gql-yt/graph"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	corsHandler := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Разрешаем доступ со всех доменов
			w.Header().Set("Access-Control-Allow-Origin", "*")
			// Разрешаем принимать куки с запросами CORS
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			// Разрешаем методы запросов
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			// Разрешаем использование определенных заголовков
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			// Продолжаем выполнение цепочки обработчиков
			h.ServeHTTP(w, r)
		})
	}

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, corsHandler(http.DefaultServeMux)))
}
