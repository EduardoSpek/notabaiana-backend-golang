package test

import "testing"

type TestCase struct {
	Esperado  any
	Recebido  any
	Descricao string
}

func Resultado(t *testing.T, esperado any, recebido any, descricao string) {
	t.Helper()
	if esperado != recebido {
		t.Errorf("Descricao: %s | Esperado: %s | Recebido: %s", descricao, esperado, recebido)
	}
}