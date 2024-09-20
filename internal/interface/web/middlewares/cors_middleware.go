package middlewares

import (
	"fmt"
	"net/http"
)

// CorsMiddleware é um middleware para permitir CORS
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		allowedOrigins := []string{
			"https://notabaiana.com.br",
			"https://www.notabaiana.com.br",
		}

		fmt.Println("ORIGIN: ", r.Header.Get("Origin"))

		origin := r.Header.Get("Origin")

		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				allowed = true
				w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
				break
			}
		}

		if !allowed {
			fmt.Println(origin)
			http.Error(w, "Origem não permitida", http.StatusForbidden)
			return
		}

		//w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		next.ServeHTTP(w, r)
	})
}
