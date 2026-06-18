package ast

import "fmt"

type Program struct {
	Statements []Statement
}

type Statement interface {
	statementNode()
	String() string
}

type MensagemStatement struct {
	Value string
}

func (m *MensagemStatement) statementNode() {}
func (m *MensagemStatement) String() string { return fmt.Sprintf(`MENSAGEM "%s"`, m.Value) }

type ChaveStatement struct {
	Value string
}

func (c *ChaveStatement) statementNode() {}
func (c *ChaveStatement) String() string { return fmt.Sprintf(`CHAVE "%s"`, c.Value) }

type CifrarStatement struct {
	Algorithm       string
	Deslocamento    int
	HasDeslocamento bool
}

func (c *CifrarStatement) statementNode() {}
func (c *CifrarStatement) String() string {
	if c.HasDeslocamento {
		return fmt.Sprintf("CIFRAR %s DESLOCAMENTO %d", c.Algorithm, c.Deslocamento)
	}
	return fmt.Sprintf("CIFRAR %s", c.Algorithm)
}

type DecifrarStatement struct {
	Algorithm       string
	Deslocamento    int
	HasDeslocamento bool
}

func (d *DecifrarStatement) statementNode() {}
func (d *DecifrarStatement) String() string {
	if d.HasDeslocamento {
		return fmt.Sprintf("DECIFRAR %s DESLOCAMENTO %d", d.Algorithm, d.Deslocamento)
	}
	return fmt.Sprintf("DECIFRAR %s", d.Algorithm)
}

type ExibirStatement struct{}

func (e *ExibirStatement) statementNode() {}
func (e *ExibirStatement) String() string { return "EXIBIR" }
