package songbook

import "net/http"

func Routes(mux *http.ServeMux) {
	mux.HandleFunc("/api/lyrics", getLyricsHandler)
}
