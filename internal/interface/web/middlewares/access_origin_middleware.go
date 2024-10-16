package middlewares

import (
	"fmt"
	"net/http"
	"strings"
)

// CorsMiddleware é um middleware para permitir CORS
func AccessOriginMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		allowedOrigins := []string{
			"https://notabaiana.com.br",
			"https://www.notabaiana.com.br",
		}
		origin := r.Header.Get("Origin")
		if origin == "" {
			origin = r.Header.Get("Referer")
		}

		fmt.Println(origin)

		allowed := false
		if origin != "" {
			for _, allowedOrigin := range allowedOrigins {
				if strings.HasPrefix(origin, allowedOrigin) {
					allowed = true
					w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
					break
				}
			}
		} else {
			// Se tanto Origin quanto Referer estiverem vazios, você pode optar por permitir a requisição
			// ou implementar uma lógica adicional de verificação
			allowed = true
		}

		if !allowed {
			http.Error(w, "Origem não permitida", http.StatusForbidden)
			return
		}

		//w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		next.ServeHTTP(w, r)
	})
}
