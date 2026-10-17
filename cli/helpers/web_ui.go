package helpers

import (
	"embed"
	"io"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"path"
	"path/filepath"
	"strings"
)

//go:embed web_dist/*
var Files embed.FS

type spaHandler struct {
	fs fs.FS
}

func (h *spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cleanPath := path.Clean(r.URL.Path)
	fsPath := strings.TrimPrefix(cleanPath, "/")
	if fsPath == "" || fsPath == "." {
		fsPath = "index.html"
	}

	f, err := h.fs.Open(fsPath)
	if err == nil {
		defer f.Close()

		stat, err := f.Stat()
		if err == nil && !stat.IsDir() {
			ext := filepath.Ext(fsPath)
			if mimeType := mime.TypeByExtension(ext); mimeType != "" {
				w.Header().Set("Content-Type", mimeType)
			}

			http.ServeContent(w, r, stat.Name(), stat.ModTime(), f.(io.ReadSeeker))
			return
		}
	}

	if path.Ext(fsPath) != "" {
		http.NotFound(w, r)
		return
	}

	indexFile, err := h.fs.Open("index.html")
	if err != nil {
		http.Error(w, "index.html not found", http.StatusInternalServerError)
		return
	}
	defer indexFile.Close()

	stat, _ := indexFile.Stat()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.ServeContent(w, r, "index.html", stat.ModTime(), indexFile.(io.ReadSeeker))
}

func Handler() http.Handler {
	subFS, err := fs.Sub(Files, "web_dist")
	if err != nil {
		panic("embedded web UI files are unavailable: " + err.Error())
	}

	return &spaHandler{
		fs: subFS,
	}
}

func StartServer(addr string) error {
	server := &http.Server{
		Addr:    addr,
		Handler: Handler(),
	}

	log.Printf("Nagare Web UI listening on http://%s", addr)
	return server.ListenAndServe()
}
