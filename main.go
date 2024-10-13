package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Basic handler
	r.Get("/welcome", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	// Basic handler with params
	r.Get("/welcome/{param}", func(w http.ResponseWriter, r *http.Request) {
		param := chi.URLParam(r, "param")
		w.Write([]byte("welcome " + param))
	})

	// Get path of current directory
	workDir, _ := os.Getwd()
	// Serve static files from the ./static directory
	filesDir := http.Dir(filepath.Join(workDir, "static"))
	FileServer(r, "/files", filesDir)

	http.ListenAndServe(":3000", r)
}

// FileServer could be renamed to fileServer (because its unexported)
func FileServer(r chi.Router, path string, root http.FileSystem) {
	// Check if path contains any URL parameters
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	// Check if path is empty
	if path != "/" && path[len(path)-1] != '/' {
		// Add / to the end of path
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	// Add * to path so it can be trimmed using strings.TrimSuffix()
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())

		// Trim RoutePattern from /files/* to /files
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")

		// Strip pathPrefix from /files/filename to filename
		// /files/other.html => other.html
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))

		// Serve the FileServer
		fs.ServeHTTP(w, r)
	})
}
