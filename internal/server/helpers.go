package server

import (
	"net"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func LimitUploadSize(maxBytes int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
			next.ServeHTTP(w, r)
		})
	}
}

func Chain(h http.Handler, mws ...Middleware) http.Handler {
	for _, mw := range slices.Backward(mws) {
		h = mw(h)
	}
	return h
}

type noListingFS struct {
	fs http.FileSystem
}

func (nfs noListingFS) Open(name string) (http.File, error) {
	f, err := nfs.fs.Open(name)
	if err != nil {
		return nil, err
	}
	stat, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, err
	}
	if stat.IsDir() {
		index := filepath.Join(name, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			f.Close()
			return nil, os.ErrNotExist
		}
	}
	return f, nil
}

func fromFS(mux *http.ServeMux, prefix string, diskPath string) {
	fs := http.FileServer(noListingFS{http.Dir(diskPath)})
	mux.Handle("GET "+prefix, http.StripPrefix(strings.TrimSuffix(prefix, "/"), fs))
}

func clientIP(r *http.Request) string {

	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
