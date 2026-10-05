package helpers

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"path"
	"strings"
)

//go:embed web_dist/*
var Files embed.FS

type spaHandler struct {
	fs http.FileSystem
}

func (h *spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := path.Clean(r.URL.Path)
	fsPath := strings.TrimPrefix(p, "/")
	if fsPath == "" || fsPath == "." {
		fsPath = "index.html"
	}

	f, err := h.fs.Open(fsPath)
	if err != nil {
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
		http.ServeContent(w, r, "index.html", stat.ModTime(), indexFile)
		return
	}
	defer f.Close()

	stat, err := f.Stat()
	if err == nil && stat.IsDir() {
		indexPath := path.Join(fsPath, "index.html")
		if idxF, err := h.fs.Open(indexPath); err == nil {
			defer idxF.Close()
			idxStat, _ := idxF.Stat()
			http.ServeContent(w, r, "index.html", idxStat.ModTime(), idxF)
			return
		}

		indexFile, err := h.fs.Open("index.html")
		if err == nil {
			defer indexFile.Close()
			idxStat, _ := indexFile.Stat()
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			http.ServeContent(w, r, "index.html", idxStat.ModTime(), indexFile)
			return
		}
	}

	http.FileServer(h.fs).ServeHTTP(w, r)
}

func Handler() http.Handler {
	subFS, err := fs.Sub(Files, "web_dist")
	if err != nil {
		panic("embedded web UI files are unavailable: " + err.Error())
	}

	return &spaHandler{
		fs: http.FS(subFS),
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
