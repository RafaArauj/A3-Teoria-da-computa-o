package main

import (
	"flag"
	"fmt"
	"os"

	"Compilador/interactive"
	"Compilador/interpreter"
	"Compilador/lexer"
	"Compilador/parser"
	"Compilador/token"
)

func main() {
	showTokens := flag.Bool("tokens", false, "exibir lista de tokens (fase léxica)")
	showAST := flag.Bool("ast", false, "exibir AST (fase sintática)")
	flag.Parse()

	if flag.NArg() == 0 {
		interactive.New().Run()
		return
	}

	runFile(flag.Arg(0), *showTokens, *showAST)
}

func runFile(path string, showTokens, showAST bool) {
	source, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "erro ao ler arquivo: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Compilador de Criptografia")
	fmt.Printf("arquivo: %s\n", path)

	// Análise Léxica
	if showTokens {
		fmt.Println("\n── Fase 1: Tokens ───────────────────────")
		tl := lexer.New(string(source))
		for {
			tok := tl.NextToken()
			fmt.Printf("  %-14s  %-28q  linha %-3d col %d\n",
				tok.Type, tok.Literal, tok.Line, tok.Column)
			if tok.Type == token.EOF {
				break
			}
		}
	}

	// Análise Sintática
	l := lexer.New(string(source))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		fmt.Println("\nErros sintáticos:")
		for _, e := range p.Errors() {
			fmt.Println("  ✗", e)
		}
		os.Exit(1)
	}

	if showAST {
		fmt.Println("\n── Fase 2: AST ──────────────────────────")
		for idx, stmt := range program.Statements {
			fmt.Printf("  [%d] %s\n", idx+1, stmt.String())
		}
	}

	// Interpretação
	fmt.Printf("\n── Fase 3: Execução (%d instruções) ─────\n", len(program.Statements))
	interp := interpreter.New()
	interp.Execute(program)
	fmt.Println("── Fim ──────────────────────────────────")
}
