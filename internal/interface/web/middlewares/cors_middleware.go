package middlewares

import (
	"net/http"
)

// CorsMiddleware é um middleware para permitir CORS
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		allowedOrigin := "https://notabaiana.com.br"
		origin := r.Header.Get("Origin")

		if origin != allowedOrigin {
			http.Error(w, "Origem não permitida", http.StatusForbidden)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		//w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		next.ServeHTTP(w, r)
	})
}
