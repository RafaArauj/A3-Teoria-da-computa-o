package interpreter

import (
	"fmt"

	"Compilador/ast"
	"Compilador/crypto"
)

type Interpreter struct {
	mensagem string
	chave    string
}

func New() *Interpreter { return &Interpreter{} }

func (i *Interpreter) Execute(program *ast.Program) {
	for _, stmt := range program.Statements {
		i.executeStatement(stmt)
	}
}

func (i *Interpreter) executeStatement(stmt ast.Statement) {
	switch s := stmt.(type) {
	case *ast.MensagemStatement:
		i.mensagem = s.Value
		fmt.Printf("  [MENSAGEM] \"%s\"\n", i.mensagem)

	case *ast.ChaveStatement:
		i.chave = s.Value
		fmt.Printf("  [CHAVE]    \"%s\"\n", i.chave)

	case *ast.CifrarStatement:
		i.executeCifrar(s)

	case *ast.DecifrarStatement:
		i.executeDecifrar(s)

	case *ast.ExibirStatement:
		fmt.Printf("\n>>> %s <<<\n\n", i.mensagem)
	}
}

func (i *Interpreter) executeCifrar(s *ast.CifrarStatement) {
	original := i.mensagem
	switch s.Algorithm {
	case "CESAR":
		desl := 3
		if s.HasDeslocamento {
			desl = s.Deslocamento
		}
		i.mensagem = crypto.CesarCifrar(i.mensagem, desl)
		fmt.Printf("  [CIFRAR]   César  (desl=%d)  \"%s\" → \"%s\"\n", desl, original, i.mensagem)

	case "ROT13":
		i.mensagem = crypto.ROT13Cifrar(i.mensagem)
		fmt.Printf("  [CIFRAR]   ROT13   \"%s\" → \"%s\"\n", original, i.mensagem)

	case "BASE64":
		i.mensagem = crypto.Base64Cifrar(i.mensagem)
		fmt.Printf("  [CIFRAR]   Base64  \"%s\" → \"%s\"\n", original, i.mensagem)

	case "XOR":
		if i.chave == "" {
			fmt.Println("  [ERRO] XOR requer CHAVE definida antes.")
			return
		}
		i.mensagem = crypto.XORCifrar(i.mensagem, i.chave)
		fmt.Printf("  [CIFRAR]   XOR     (chave=%q) \"%s\" → \"%s\"\n", i.chave, original, i.mensagem)
	}
}

func (i *Interpreter) executeDecifrar(s *ast.DecifrarStatement) {
	original := i.mensagem
	switch s.Algorithm {
	case "CESAR":
		desl := 3
		if s.HasDeslocamento {
			desl = s.Deslocamento
		}
		i.mensagem = crypto.CesarDecifrar(i.mensagem, desl)
		fmt.Printf("  [DECIFRAR] César  (desl=%d)  \"%s\" → \"%s\"\n", desl, original, i.mensagem)

	case "ROT13":
		i.mensagem = crypto.ROT13Decifrar(i.mensagem)
		fmt.Printf("  [DECIFRAR] ROT13   \"%s\" → \"%s\"\n", original, i.mensagem)

	case "BASE64":
		decoded, err := crypto.Base64Decifrar(i.mensagem)
		if err != nil {
			fmt.Printf("  [ERRO] %v\n", err)
			return
		}
		i.mensagem = decoded
		fmt.Printf("  [DECIFRAR] Base64  \"%s\" → \"%s\"\n", original, i.mensagem)

	case "XOR":
		if i.chave == "" {
			fmt.Println("  [ERRO] XOR requer CHAVE definida antes.")
			return
		}
		decoded, err := crypto.XORDecifrar(i.mensagem, i.chave)
		if err != nil {
			fmt.Printf("  [ERRO] %v\n", err)
			return
		}
		i.mensagem = decoded
		fmt.Printf("  [DECIFRAR] XOR     (chave=%q) \"%s\" → \"%s\"\n", i.chave, original, i.mensagem)
	}
}
