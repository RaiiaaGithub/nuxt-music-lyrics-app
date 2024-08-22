package songbook

import "net/http"

func Routes() {
	http.HandleFunc("/api/lyrics", handleGetLyrics)
}
