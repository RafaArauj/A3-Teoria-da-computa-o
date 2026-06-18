package token

type TokenType string

const (
	STRING  TokenType = "STRING"
	NUMBER  TokenType = "NUMBER"
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	MENSAGEM TokenType = "MENSAGEM"
	CHAVE    TokenType = "CHAVE"
	CIFRAR   TokenType = "CIFRAR"
	DECIFRAR TokenType = "DECIFRAR"
	EXIBIR   TokenType = "EXIBIR"

	CESAR        TokenType = "CESAR"
	ROT13        TokenType = "ROT13"
	BASE64       TokenType = "BASE64"
	XOR          TokenType = "XOR"
	DESLOCAMENTO TokenType = "DESLOCAMENTO"
)

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

var keywords = map[string]TokenType{
	"MENSAGEM":     MENSAGEM,
	"CHAVE":        CHAVE,
	"CIFRAR":       CIFRAR,
	"DECIFRAR":     DECIFRAR,
	"EXIBIR":       EXIBIR,
	"CESAR":        CESAR,
	"ROT13":        ROT13,
	"BASE64":       BASE64,
	"XOR":          XOR,
	"DESLOCAMENTO": DESLOCAMENTO,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return ILLEGAL
}
