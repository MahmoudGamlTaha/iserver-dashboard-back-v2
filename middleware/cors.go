package middleware

import (
	"net/http"

	"github.com/rs/cors"
)

// CorsMiddleware returns a new CORS middleware handler.
func CorsMiddleware() func(http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization", "X-Requested-With", "Origin", "Accept", "Application/json", "User-Agent"},
	})
	return c.Handler
}
