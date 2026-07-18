package server

import (
	"errors"
	"net/http"
	"os"
	"strings"
)

var cachedCSS []byte

func loadViewerCSS() error {
	viewerCSS, err := os.ReadFile("/opt/gochive/lib/pdfjs/web/viewer.css")
	if err != nil {
		return errors.New("cannot read viewer.css")
	}

	userCSS, err := os.ReadFile("static/styles/userContent.css")
	if err != nil {
		return errors.New("cannot read userContent.css")
	}

	buf := make([]byte, 0, len(viewerCSS)+len(userCSS)+32)
	buf = append(buf, viewerCSS...)
	buf = append(buf, "\n\n/* userContent.css */\n"...)
	buf = append(buf, userCSS...)

	cachedCSS = buf
	return nil
}

func injectContentCSS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !strings.HasSuffix(r.URL.Path, "viewer.css") {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Type", "text/css")
		w.WriteHeader(http.StatusOK)
		w.Write(cachedCSS)

	})
}
