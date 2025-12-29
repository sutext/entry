package main

import (
	"sutext.github.io/entry/model/sqlite"
	"sutext.github.io/entry/server"
)

func main() {
	s := server.New(
		server.WithDriver(sqlite.Named("tester.db")),
		server.WithIssuerURL("http://localhost:8080/"),
	)
	s.Serve()
}
