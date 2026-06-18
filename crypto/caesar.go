package crypto

// CesarCifrar aplica a cifra de César com o deslocamento dado.
func CesarCifrar(text string, shift int) string {
	shift = ((shift % 26) + 26) % 26
	result := make([]byte, len(text))
	for i := 0; i < len(text); i++ {
		ch := text[i]
		switch {
		case ch >= 'A' && ch <= 'Z':
			result[i] = 'A' + (ch-'A'+byte(shift))%26
		case ch >= 'a' && ch <= 'z':
			result[i] = 'a' + (ch-'a'+byte(shift))%26
		default:
			result[i] = ch
		}
	}
	return string(result)
}

// CesarDecifrar reverte a cifra de César.
func CesarDecifrar(text string, shift int) string {
	return CesarCifrar(text, 26-shift)
}
