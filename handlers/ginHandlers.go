package handlers

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Nurs005/gql-yt/graph"
	"github.com/gin-gonic/gin"
)

func HandleGQL() gin.HandlerFunc {
	h := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
	return func(ctx *gin.Context) {
		h.ServeHTTP(ctx.Writer, ctx.Request)
	}
}

func HandlePlyground() gin.HandlerFunc {
	h := playground.Handler("GraphQL playground", "/query")
	return func(ctx *gin.Context) {
		h.ServeHTTP(ctx.Writer, ctx.Request)
	}
}
