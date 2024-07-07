package main

import (
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type spaHandler struct {
	staticPath string
	indexPath  string
}

func guessMimeType(p string) string {
	if strings.HasSuffix(p, ".js") {
		return "application/javascript"
	}
	if strings.HasSuffix(p, ".css") {
		return "text/css"
	}
	if strings.HasSuffix(p, ".html") {
		return "text/html"
	}

	return "text/plain"
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	computedPath := filepath.Join(h.staticPath, strings.Trim(r.URL.Path, "/web"))

	// Check whether a file exists
	fi, err := os.Stat(computedPath)
	if os.IsNotExist(err) || fi.IsDir() {
		slog.Error("invalid file", "path", computedPath, "reqpath", r.URL.Path)
		// Serve index.html anytime a specific existing file isn't the requested object
		w.Header().Set("Content-TYpe", "text/html")
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	}

	if err != nil {
		// this is an unrecoverable error - we had an error checking the filesystem other than "file doesn't exist"
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	mtype := guessMimeType(computedPath)
	w.Header().Set("Content-Type", mtype)

	// If the request is for a specific file, use the static fileserver
	slog.Debug("Serving specific file", "path", computedPath, "reqpath", r.URL.Path, "mtype", mtype)
	http.ServeFile(w, r, computedPath)
}
