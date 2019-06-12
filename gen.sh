#!/usr/bin/env bash

go run ./gen/gen.go # generate entities

mv resolver.go resolver.old
go run github.com/99designs/gqlgen # generate graphQL

cp resolver.go resolver.go.new
git merge-file resolver.go resolver.go resolver.old