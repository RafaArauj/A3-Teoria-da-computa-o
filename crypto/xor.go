package crypto

import (
	"encoding/hex"
	"fmt"
)

// XORCifrar aplica XOR byte a byte e codifica o resultado em hexadecimal.
func XORCifrar(text, key string) string {
	if key == "" {
		return text
	}
	result := make([]byte, len(text))
	for i := 0; i < len(text); i++ {
		result[i] = text[i] ^ key[i%len(key)]
	}
	return hex.EncodeToString(result)
}

// XORDecifrar decodifica o hexadecimal e reverte o XOR.
func XORDecifrar(hexText, key string) (string, error) {
	if key == "" {
		return hexText, nil
	}
	data, err := hex.DecodeString(hexText)
	if err != nil {
		return "", fmt.Errorf("texto XOR inválido (esperado hexadecimal): %w", err)
	}
	result := make([]byte, len(data))
	for i, b := range data {
		result[i] = b ^ key[i%len(key)]
	}
	return string(result), nil
}
