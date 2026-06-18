package crypto

// ROT13 é simétrico: cifrar e decifrar são a mesma operação.
func ROT13Cifrar(text string) string   { return CesarCifrar(text, 13) }
func ROT13Decifrar(text string) string { return CesarCifrar(text, 13) }
