package utils

import (
	"net"
	"net/http"
)

func GetIP(r *http.Request) string {	

	// Verifica se o IP está no cabeçalho X-Forwarded-For
	ip := r.Header.Get("X-Forwarded-For")
	
	if ip == "" {
		// Se não estiver, verifica o cabeçalho X-Real-Ip
		ip = r.Header.Get("X-Real-Ip")
	}

	if ip == "" {
		// Caso contrário, usa o IP remoto
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	}

	return ip
}