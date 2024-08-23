package main

import (
	"fmt"
	"net/http"

	"github.com/RaiiaaGithub/vue-music-lyrics-app/pkg/songbook"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

	// Create a new CORS handle
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"}, // Add your frontend URL
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		Debug:          true,
	})

	songbook.Routes(mux)

	handler := c.Handler(mux)

	fmt.Println("Server listening on http://localhost:8080/api/lyrics")
	http.ListenAndServe(":8080", handler)
}
