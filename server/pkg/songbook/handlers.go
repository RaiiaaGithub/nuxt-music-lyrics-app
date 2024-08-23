package songbook

import (
	"encoding/json"
	"net/http"

	"github.com/RaiiaaGithub/vue-music-lyrics-app/pkg/utils"
)

func getLyricsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	song, err := GetTopHitSong(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	utils.LogDebug("GET LYRICS - %s by %s", song.Title, song.Artist)
	json.NewEncoder(w).Encode(song)
}
