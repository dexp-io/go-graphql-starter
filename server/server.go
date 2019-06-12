package main

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/handler"
	"github.com/drupalexp/dexp"
	"github.com/go-chi/chi"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
)

func main() {

	var post = os.Getenv("PORT")

	dexp.Mysql()
	dexp.Test()

	router := chi.NewRouter()
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:8080"},
		AllowCredentials: true,
		Debug:            false,
		AllowedHeaders:   []string{"Authorization", "Content-Type", "Origin"},
	}).Handler)

	router.Use(dexp.Middleware())

	c := dexp.Config{Resolvers: &dexp.Resolver{}}
	c.Directives.HasRole = func(ctx context.Context, obj interface{}, next graphql.Resolver, role []*dexp.UserRole) (res interface{}, err error) {

		if len(role) > 0 {
			// process handle role here
		}

		return next(ctx)

	}
	router.Handle("/", handler.Playground("GraphQL playground", "/query"))
	router.Handle("/query", handler.GraphQL(dexp.NewExecutableSchema(c)))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", post)
	err := http.ListenAndServe(":"+post, router)
	if err != nil {
		panic(err)
	}

}
