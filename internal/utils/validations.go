package utils

import "regexp"

func IsValidEmail(email string) bool {
	// Define a expressão regular para validar o email
	// Esta expressão cobre a maioria dos casos comuns de email, mas pode não cobrir todos os casos possíveis
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// Compila a expressão regular
	re := regexp.MustCompile(emailRegex)

	// Retorna verdadeiro se o email corresponder à expressão regular, falso caso contrário
	return re.MatchString(email)
}