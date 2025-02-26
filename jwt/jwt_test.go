package jwt

import (
	"testing"
	"time"
)

func TestJWT(t *testing.T) {

	// Dados de teste
	claims := map[string]interface{}{
		"email":    "jose@gmail.com",
		"password": "jose",
		"role":     "USER",
		"age":      "12",
	}
	expiration := 0 // 1 segundo
	customKey := "arthur"

	// Teste 1: Criação de token com chave padrão
	tokenDefault, err := CreateToken(claims, customKey, "jose", int64(expiration))
	if err != nil {
		t.Fatalf("Erro ao criar token com chave padrão: %v", err)
	}
	if tokenDefault == "" {
		t.Error("Token com chave padrão não deveria ser vazio")
	}

	// Teste 1: Validação de token válido
	if !IsValidJwt(tokenDefault, customKey) {
		t.Error("Token válido com chave padrão deveria ser validado como true")
	}

	// Teste 2: Extração de valores
	valuesDefault, err := GetValuesFrom(tokenDefault, customKey)
	if err != nil {
		t.Fatalf("Erro ao extrair valores do token com chave padrão: %v", err)
	}
	if valuesDefault["email"] != claims["email"] {
		t.Errorf("Valor 'email' extraído (%v) diferente do esperado (%v)", valuesDefault["email"], claims["email"])
	}

	// Teste 3: Validação de token expirado
	time.Sleep(3 * time.Second) // Espera o token expirar (3s + margem)
	if IsValidJwt(tokenDefault, customKey) {
		t.Error("Um token já expirado não deveria ser válido nem com a chave padrão")
	}
}
