package middlewares

import (
	"net/http"

	"github.com/eduardospek/notabaiana-backend-golang/internal/utils"
)

func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" {
			http.Error(w, "Acesso não autorizado", http.StatusForbidden)
			return
		}

		tokenStr = tokenStr[len("Bearer "):]

		_, err := utils.ValidateJWT(tokenStr)
		if err != nil {
			http.Error(w, "Acesso não autorizado: Token inválido!", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
