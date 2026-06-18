package parser

import (
	"fmt"
	"strconv"

	"Compilador/ast"
	"Compilador/lexer"
	"Compilador/token"
)

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
	errors    []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) Errors() []string { return p.errors }

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.MENSAGEM:
		return p.parseMensagem()
	case token.CHAVE:
		return p.parseChave()
	case token.CIFRAR:
		return p.parseCifrar()
	case token.DECIFRAR:
		return p.parseDecifrar()
	case token.EXIBIR:
		return &ast.ExibirStatement{}
	default:
		p.addError("comando desconhecido '%s'", p.curToken.Literal)
		return nil
	}
}

func (p *Parser) parseMensagem() *ast.MensagemStatement {
	if !p.expectPeek(token.STRING) {
		return nil
	}
	return &ast.MensagemStatement{Value: p.curToken.Literal}
}

func (p *Parser) parseChave() *ast.ChaveStatement {
	if !p.expectPeek(token.STRING) {
		return nil
	}
	return &ast.ChaveStatement{Value: p.curToken.Literal}
}

func (p *Parser) parseCifrar() *ast.CifrarStatement {
	algo := p.parseAlgorithm()
	if algo == "" {
		return nil
	}
	stmt := &ast.CifrarStatement{Algorithm: algo}
	if p.peekToken.Type == token.DESLOCAMENTO {
		p.nextToken()
		if !p.expectPeek(token.NUMBER) {
			return nil
		}
		n, err := strconv.Atoi(p.curToken.Literal)
		if err != nil {
			p.addError("deslocamento inválido '%s'", p.curToken.Literal)
			return nil
		}
		stmt.Deslocamento = n
		stmt.HasDeslocamento = true
	}
	return stmt
}

func (p *Parser) parseDecifrar() *ast.DecifrarStatement {
	algo := p.parseAlgorithm()
	if algo == "" {
		return nil
	}
	stmt := &ast.DecifrarStatement{Algorithm: algo}
	if p.peekToken.Type == token.DESLOCAMENTO {
		p.nextToken()
		if !p.expectPeek(token.NUMBER) {
			return nil
		}
		n, err := strconv.Atoi(p.curToken.Literal)
		if err != nil {
			p.addError("deslocamento inválido '%s'", p.curToken.Literal)
			return nil
		}
		stmt.Deslocamento = n
		stmt.HasDeslocamento = true
	}
	return stmt
}

func (p *Parser) parseAlgorithm() string {
	p.nextToken()
	switch p.curToken.Type {
	case token.CESAR:
		return "CESAR"
	case token.ROT13:
		return "ROT13"
	case token.BASE64:
		return "BASE64"
	case token.XOR:
		return "XOR"
	default:
		p.addError("algoritmo esperado (CESAR, ROT13, BASE64, XOR), encontrado '%s'", p.curToken.Literal)
		return ""
	}
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	}
	p.addError("esperado %s, encontrado '%s'", t, p.peekToken.Literal)
	return false
}

func (p *Parser) addError(format string, args ...any) {
	msg := fmt.Sprintf("linha %d, coluna %d: ", p.curToken.Line, p.curToken.Column)
	msg += fmt.Sprintf(format, args...)
	p.errors = append(p.errors, msg)
}
