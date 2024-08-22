package main

import (
	"github.com/RaiiaaGithub/vue-music-lyrics-app/pkg/genius"
	"github.com/RaiiaaGithub/vue-music-lyrics-app/pkg/utils"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	searchQuery := "When I was your man"
	topHitSong, err := genius.GetTopHitSong(searchQuery)
	if err != nil {
		panic(err)
	}

	utils.LogDebug("Top Hit Song: %s by %s", topHitSong.Title, topHitSong.Artist)
	utils.LogDebug("Lyrics: %v", topHitSong.Lyrics)
}
