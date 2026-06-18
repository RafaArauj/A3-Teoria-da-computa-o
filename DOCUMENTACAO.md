# Compilador de Criptografia — CipherScript

Trabalho acadêmico inspirado na Máquina Enigma utilizada na Segunda Guerra Mundial.  
O projeto implementa um compilador para uma linguagem própria chamada **CipherScript**, cujos programas descrevem operações de criptografia e descriptografia de mensagens.

---

## Sumário

1. [Visão Geral](#visão-geral)
2. [Como Executar](#como-executar)
3. [Arquitetura do Compilador](#arquitetura-do-compilador)
4. [A Linguagem CipherScript](#a-linguagem-cipherscript)
5. [Pacotes](#pacotes)
   - [token](#token)
   - [lexer](#lexer)
   - [ast](#ast)
   - [parser](#parser)
   - [interpreter](#interpreter)
   - [crypto](#crypto)
   - [interactive](#interactive)
6. [Algoritmos de Criptografia](#algoritmos-de-criptografia)
7. [Modo Interativo](#modo-interativo)
8. [Formato do Pacote](#formato-do-pacote)

---

## Visão Geral

O compilador funciona em duas formas de uso:

- **Modo interativo** — terminal para cifrar e decifrar mensagens em tempo real, gerando um código compacto para compartilhar.
- **Modo arquivo** — lê um arquivo `.cipher` escrito em CipherScript, analisa o código nas três fases clássicas de um compilador e executa as instruções.

```
Código fonte (.cipher)
        │
        ▼
  ┌───────────┐
  │   Lexer   │  Análise Léxica   → sequência de tokens
  └───────────┘
        │
        ▼
  ┌───────────┐
  │  Parser   │  Análise Sintática → Árvore Sintática Abstrata (AST)
  └───────────┘
        │
        ▼
  ┌─────────────┐
  │ Interpreter │  Execução       → operações de criptografia
  └─────────────┘
```

---

## Como Executar

**Pré-requisito:** Go 1.21 ou superior.

```bash
# Modo interativo (sem argumentos)
go run .

# Modo arquivo
go run . exemplos/exemplo.cipher

# Modo arquivo com debug léxico
go run . -tokens exemplos/exemplo.cipher

# Modo arquivo com debug da AST
go run . -ast exemplos/exemplo.cipher

# Compilar o binário
go build -o compilador .
./compilador
```

---

## Arquitetura do Compilador

### Fase 1 — Análise Léxica (Lexer)

O **lexer** lê o código fonte caractere por caractere e agrupa os caracteres em unidades chamadas **tokens**. Cada token tem um tipo e um valor literal.

Exemplo de entrada:
```
CIFRAR CESAR DESLOCAMENTO 7
```

Tokens gerados:
```
CIFRAR        "CIFRAR"        linha 1 col 1
CESAR         "CESAR"         linha 1 col 8
DESLOCAMENTO  "DESLOCAMENTO"  linha 1 col 14
NUMBER        "7"             linha 1 col 27
```

### Fase 2 — Análise Sintática (Parser)

O **parser** consome os tokens e verifica se eles seguem a gramática da linguagem. Se sim, constrói a **AST** (Árvore Sintática Abstrata), uma estrutura de dados que representa o programa de forma hierárquica.

Cada instrução vira um nó da árvore. O programa inteiro é uma lista de nós (`Program.Statements`).

### Fase 3 — Interpretação

O **interpreter** percorre a AST e executa cada instrução em ordem, mantendo um estado interno com a mensagem e a chave atuais.

---

## A Linguagem CipherScript

CipherScript é case-sensitive. Todos os comandos são escritos em maiúsculas.

### Comandos

| Comando | Descrição |
|---|---|
| `MENSAGEM "texto"` | Define a mensagem que será processada |
| `CHAVE "chave"` | Define a chave de criptografia (usada pelo XOR) |
| `CIFRAR CESAR DESLOCAMENTO n` | Cifra com a Cifra de César usando deslocamento n (padrão: 3) |
| `CIFRAR ROT13` | Cifra com ROT13 |
| `CIFRAR BASE64` | Codifica em Base64 |
| `CIFRAR XOR` | Cifra com XOR byte a byte (requer CHAVE) |
| `DECIFRAR CESAR DESLOCAMENTO n` | Reverte a Cifra de César |
| `DECIFRAR ROT13` | Reverte o ROT13 |
| `DECIFRAR BASE64` | Decodifica Base64 |
| `DECIFRAR XOR` | Reverte o XOR (requer CHAVE) |
| `EXIBIR` | Exibe a mensagem atual no terminal |
| `// texto` | Comentário de linha (ignorado pelo compilador) |

### Gramática

```
Programa     → Instrucao*
Instrucao    → MensagemCmd | ChaveCmd | CifrarCmd | DecifrarCmd | ExibirCmd
MensagemCmd  → MENSAGEM STRING
ChaveCmd     → CHAVE STRING
CifrarCmd    → CIFRAR Algoritmo [DESLOCAMENTO NUMBER]
DecifrarCmd  → DECIFRAR Algoritmo [DESLOCAMENTO NUMBER]
ExibirCmd    → EXIBIR
Algoritmo    → CESAR | ROT13 | BASE64 | XOR
```

### Exemplo de programa

```cipher
// Cifra de César
MENSAGEM "Atacar ao amanhecer"
CIFRAR CESAR DESLOCAMENTO 7
EXIBIR
DECIFRAR CESAR DESLOCAMENTO 7
EXIBIR

// XOR com chave
MENSAGEM "Operacao secreta"
CHAVE "SENHA"
CIFRAR XOR
EXIBIR
DECIFRAR XOR
EXIBIR
```

---

## Pacotes

### token

**Arquivo:** `token/token.go`

Define os tipos de tokens reconhecidos pelo lexer e a tabela de palavras-chave.

```go
type TokenType string

type Token struct {
    Type    TokenType
    Literal string
    Line    int
    Column  int
}
```

**Tipos de token:**

| Constante | Valor | Descrição |
|---|---|---|
| `STRING` | `"STRING"` | Texto entre aspas duplas |
| `NUMBER` | `"NUMBER"` | Número inteiro |
| `ILLEGAL` | `"ILLEGAL"` | Caractere não reconhecido |
| `EOF` | `"EOF"` | Fim do arquivo |
| `MENSAGEM` | `"MENSAGEM"` | Palavra-chave |
| `CHAVE` | `"CHAVE"` | Palavra-chave |
| `CIFRAR` | `"CIFRAR"` | Palavra-chave |
| `DECIFRAR` | `"DECIFRAR"` | Palavra-chave |
| `EXIBIR` | `"EXIBIR"` | Palavra-chave |
| `CESAR` | `"CESAR"` | Nome de algoritmo |
| `ROT13` | `"ROT13"` | Nome de algoritmo |
| `BASE64` | `"BASE64"` | Nome de algoritmo |
| `XOR` | `"XOR"` | Nome de algoritmo |
| `DESLOCAMENTO` | `"DESLOCAMENTO"` | Parâmetro do César |

**Função principal:**

```go
func LookupIdent(ident string) TokenType
```

Recebe uma string lida pelo lexer e retorna o `TokenType` correspondente. Se a string for uma palavra-chave conhecida, retorna o token dela; caso contrário, retorna `ILLEGAL`.

---

### lexer

**Arquivo:** `lexer/lexer.go`

Responsável pela **análise léxica**: percorre o código fonte e produz tokens um a um.

```go
type Lexer struct {
    input        string  // código fonte completo
    position     int     // posição do caractere atual
    readPosition int     // próxima posição a ler
    ch           byte    // caractere atual
    line         int     // linha atual (para erros)
    column       int     // coluna atual (para erros)
}
```

**Funções:**

| Função | Descrição |
|---|---|
| `New(input string) *Lexer` | Cria um novo lexer e já avança para o primeiro caractere |
| `NextToken() token.Token` | Lê e retorna o próximo token |
| `readChar()` | Avança um caractere, atualizando posição e linha |
| `peekChar() byte` | Espia o próximo caractere sem avançar |
| `skipWhitespace()` | Ignora espaços, tabs e quebras de linha |
| `skipLineComment()` | Ignora tudo após `//` até o fim da linha |
| `readString() string` | Lê o conteúdo entre aspas duplas |
| `readIdentifier() string` | Lê uma sequência de letras e dígitos |
| `readNumber() string` | Lê uma sequência de dígitos |

**Fluxo do `NextToken`:**

1. Ignora espaços em branco.
2. Se encontra `//`, ignora até o fim da linha e chama a si mesmo recursivamente.
3. Se encontra `"`, lê uma string.
4. Se encontra `0` (byte nulo), retorna `EOF`.
5. Se é letra, lê um identificador e consulta a tabela de palavras-chave.
6. Se é dígito, lê um número.
7. Qualquer outro caractere vira `ILLEGAL`.

---

### ast

**Arquivo:** `ast/ast.go`

Define os nós da **Árvore Sintática Abstrata**. Cada tipo de instrução da linguagem tem um nó correspondente.

```go
type Program struct {
    Statements []Statement
}

type Statement interface {
    statementNode()
    String() string
}
```

**Nós disponíveis:**

| Tipo | Campos | Representa |
|---|---|---|
| `MensagemStatement` | `Value string` | `MENSAGEM "texto"` |
| `ChaveStatement` | `Value string` | `CHAVE "chave"` |
| `CifrarStatement` | `Algorithm string`, `Deslocamento int`, `HasDeslocamento bool` | `CIFRAR <algo> [DESLOCAMENTO n]` |
| `DecifrarStatement` | `Algorithm string`, `Deslocamento int`, `HasDeslocamento bool` | `DECIFRAR <algo> [DESLOCAMENTO n]` |
| `ExibirStatement` | — | `EXIBIR` |

Todos os nós implementam `String()`, que reconstrói a instrução em texto para fins de debug.

---

### parser

**Arquivo:** `parser/parser.go`

Responsável pela **análise sintática**: consome tokens do lexer e constrói a AST, validando a gramática.

```go
type Parser struct {
    l         *lexer.Lexer
    curToken  token.Token   // token atual
    peekToken token.Token   // próximo token (lookahead)
    errors    []string      // erros encontrados
}
```

O parser usa **lookahead de 1 token** (`peekToken`) para tomar decisões sem precisar recuar.

**Funções principais:**

| Função | Descrição |
|---|---|
| `New(l *lexer.Lexer) *Parser` | Cria o parser e já carrega os dois primeiros tokens |
| `ParseProgram() *ast.Program` | Ponto de entrada: processa todas as instruções |
| `Errors() []string` | Retorna a lista de erros sintáticos encontrados |
| `parseStatement()` | Despacha para o método correto conforme o token atual |
| `parseMensagem()` | Processa `MENSAGEM STRING` |
| `parseChave()` | Processa `CHAVE STRING` |
| `parseCifrar()` | Processa `CIFRAR Algoritmo [DESLOCAMENTO NUMBER]` |
| `parseDecifrar()` | Processa `DECIFRAR Algoritmo [DESLOCAMENTO NUMBER]` |
| `parseAlgorithm()` | Avança o cursor e retorna o nome canônico do algoritmo |
| `expectPeek(t TokenType)` | Avança se o próximo token for do tipo esperado, senão registra erro |
| `addError(format, args)` | Adiciona erro com número de linha e coluna |

**Tratamento de erros:** quando uma instrução inválida é encontrada, o parser registra o erro, descarta o nó (`nil`) e continua tentando analisar as próximas instruções. Ao final, se `Errors()` não estiver vazio, a execução é abortada.

---

### interpreter

**Arquivo:** `interpreter/interpreter.go`

Percorre a AST e **executa** cada instrução, mantendo estado interno.

```go
type Interpreter struct {
    mensagem string  // mensagem sendo processada no momento
    chave    string  // chave de criptografia atual
}
```

**Funções:**

| Função | Descrição |
|---|---|
| `New() *Interpreter` | Cria um interpreter com estado vazio |
| `Execute(program *ast.Program)` | Percorre e executa todos os nós do programa |
| `executeStatement(stmt)` | Despacha para o handler correto por tipo de nó |
| `executeCifrar(s)` | Aplica o algoritmo de cifra e atualiza `mensagem` |
| `executeDecifrar(s)` | Aplica o algoritmo de decifra e atualiza `mensagem` |

**Estado:** o interpreter mantém `mensagem` e `chave` entre instruções. `MENSAGEM` sobrescreve o valor atual; cada operação `CIFRAR`/`DECIFRAR` modifica `mensagem` no lugar. `EXIBIR` apenas imprime sem alterar o estado.

---

### crypto

Pacote com as implementações puras dos algoritmos. Sem estado; todas as funções são puras (entrada → saída).

#### caesar.go

```go
func CesarCifrar(text string, shift int) string
func CesarDecifrar(text string, shift int) string
```

Desloca cada letra `shift` posições no alfabeto. Letras maiúsculas e minúsculas são tratadas separadamente; outros caracteres (espaços, pontuação, números) passam sem alteração.

O deslocamento é normalizado com `((shift % 26) + 26) % 26` para suportar qualquer valor inteiro, incluindo negativos.

`CesarDecifrar` simplesmente chama `CesarCifrar` com `26 - shift`.

#### rot13.go

```go
func ROT13Cifrar(text string) string
func ROT13Decifrar(text string) string
```

Caso especial da Cifra de César com deslocamento fixo de 13. Como 13 + 13 = 26 (volta ao início), cifrar e decifrar são a mesma operação. Ambas as funções chamam `CesarCifrar(text, 13)`.

#### base64enc.go

```go
func Base64Cifrar(text string) string
func Base64Decifrar(text string) (string, error)
```

Codificação Base64 padrão (`encoding/base64` da stdlib). Converte bytes arbitrários em texto ASCII. `Base64Decifrar` retorna erro se a entrada não for Base64 válido.

#### xor.go

```go
func XORCifrar(text, key string) string
func XORDecifrar(hexText, key string) (string, error)
```

Aplica XOR byte a byte entre a mensagem e a chave (repetida ciclicamente). O resultado é codificado em **hexadecimal** para que o texto cifrado seja imprimível.

`XORDecifrar` decodifica o hexadecimal e aplica XOR novamente (operação simétrica).

**Propriedade importante:** `x XOR x = 0`. Se a chave for idêntica à mensagem, o resultado será uma sequência de zeros. Decifrar zeros com qualquer chave `k` produz `k` repetida. Por isso, a chave deve sempre ser diferente da mensagem.

---

### interactive

**Arquivo:** `interactive/terminal.go`

Implementa o **modo interativo** do compilador: um terminal de linha de comando para cifrar e decifrar mensagens sem precisar escrever um arquivo `.cipher`.

```go
type Terminal struct {
    reader *bufio.Reader
}
```

**Funções:**

| Função | Descrição |
|---|---|
| `New() *Terminal` | Cria o terminal ligado ao stdin |
| `Run()` | Loop principal: exibe menu e despacha opções |
| `flowCifrar()` | Guia o usuário para cifrar uma mensagem |
| `flowDecifrar()` | Guia o usuário para decifrar um pacote recebido |
| `chooseAlgorithm()` | Exibe o submenu de algoritmos e retorna a escolha |
| `buildPacket(algo, param, ciphertext)` | Monta o código compacto para compartilhar |
| `parsePacket(packet)` | Extrai algoritmo, parâmetro e texto cifrado de um pacote |
| `readLine()` | Lê uma linha do stdin, removendo espaços e `\n` |
| `paramToShift(param)` | Converte a string do parâmetro em inteiro (padrão 3) |

---

## Algoritmos de Criptografia

### Cifra de César

Desloca cada letra `n` posições no alfabeto. Com deslocamento 3: A → D, B → E, etc.

```
Original:  A B C D E F G ... X Y Z
Cifrado:   D E F G H I J ... A B C
```

Para decifrar, aplica o deslocamento inverso: `26 - n`.

### ROT13

Caso específico da Cifra de César com deslocamento 13. Por ser exatamente a metade de 26, a operação é seu próprio inverso: ROT13(ROT13(x)) = x.

```
Original: A B C ... M N O ... Z
Cifrado:  N O P ... Z A B ... M
```

### Base64

Representa dados binários usando apenas 64 caracteres ASCII (`A-Z`, `a-z`, `0-9`, `+`, `/`, `=` para padding). Não é uma cifra de segurança; é uma codificação para transmitir dados em meios que aceitam apenas texto.

```
"Ola" → "T2xh"
```

### XOR

Aplica a operação lógica XOR (OU exclusivo) bit a bit entre cada byte da mensagem e o byte correspondente da chave (repetida ciclicamente). O resultado é exibido em hexadecimal.

```
Mensagem: O  l  a
ASCII:    79 6C 61
Chave:    K  E  Y  (repetida)
ASCII:    4B 45 59
XOR:      32 29 38  → "32293b" em hex
```

Para decifrar: aplica XOR novamente com a mesma chave (operação simétrica).

---

## Modo Interativo

Executar `go run .` sem argumentos inicia o terminal interativo.

```
Compilador de Criptografia
Inspirado na Maquina Enigma

1 - Cifrar mensagem
2 - Decifrar mensagem
0 - Sair
>
```

### Cifrar

1. Digitar a mensagem.
2. Escolher o algoritmo (1 = César, 2 = ROT13, 3 = Base64, 4 = XOR).
3. Para César: informar o deslocamento (padrão: 3).
4. Para XOR: informar a chave.
5. O terminal exibe o **pacote** gerado, pronto para copiar e enviar.

### Decifrar

1. Colar o pacote recebido.
2. Se o algoritmo precisar de chave (XOR), informá-la.
3. O terminal exibe a mensagem original.

---

## Formato do Pacote

O pacote é uma string auto-descritiva gerada pelo modo interativo:

```
[ALGORITMO:PARAMETRO]texto_cifrado
```

O parâmetro só existe para o César (o deslocamento). Para os demais algoritmos:

```
[ALGORITMO]texto_cifrado
```

**Exemplos:**

| Algoritmo | Pacote |
|---|---|
| César (desl. 7) | `[CESAR:7]Hwwhjhy hw hhwkljly` |
| ROT13 | `[ROT13]Zrafntrz frpergn` |
| Base64 | `[BASE64]T2xhIE11bmRv` |
| XOR | `[XOR]0200074441371d05` |

O receptor cola o pacote na opção "Decifrar". O algoritmo é detectado automaticamente a partir do cabeçalho. Para XOR, a chave deve ser combinada previamente por outro canal.
