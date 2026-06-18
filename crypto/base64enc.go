package crypto

import (
	"encoding/base64"
	"fmt"
)

func Base64Cifrar(text string) string {
	return base64.StdEncoding.EncodeToString([]byte(text))
}

func Base64Decifrar(text string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return "", fmt.Errorf("texto Base64 inválido: %w", err)
	}
	return string(decoded), nil
}
