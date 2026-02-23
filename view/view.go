package view

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed dist/*
var files embed.FS

func FileServer() (http.Handler, error) {
	dfs, err := fs.Sub(files, "dist")
	if err != nil {
		return nil, err
	}
	return http.FileServerFS(dfs), nil
}
