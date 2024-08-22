package main

import (
	"fmt"
	"net/http"

	"github.com/RaiiaaGithub/vue-music-lyrics-app/pkg/songbook"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	songbook.Routes()

	fmt.Println("Server listening on http://localhost:8080/api/lyrics")
	http.ListenAndServe(":8080", nil)
}
