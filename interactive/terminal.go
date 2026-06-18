package interactive

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"Compilador/crypto"
)

type Terminal struct {
	reader *bufio.Reader
}

func New() *Terminal {
	return &Terminal{reader: bufio.NewReader(os.Stdin)}
}

func (t *Terminal) Run() {
	fmt.Println("Compilador de Criptografia")
	fmt.Println("Inspirado na Maquina Enigma")
	fmt.Println()

	for {
		fmt.Println("1 - Cifrar mensagem")
		fmt.Println("2 - Decifrar mensagem")
		fmt.Println("0 - Sair")
		fmt.Print("> ")

		switch t.readLine() {
		case "1":
			fmt.Println()
			t.flowCifrar()
		case "2":
			fmt.Println()
			t.flowDecifrar()
		case "0", "q", "sair":
			return
		default:
			fmt.Println("opcao invalida")
		}
		fmt.Println()
	}
}

func (t *Terminal) flowCifrar() {
	fmt.Print("mensagem: ")
	message := t.readLine()
	if message == "" {
		fmt.Println("erro: mensagem vazia")
		return
	}

	algo, param := t.chooseAlgorithm()
	if algo == "" {
		return
	}

	var ciphertext, keyNote string

	switch algo {
	case "CESAR":
		ciphertext = crypto.CesarCifrar(message, paramToShift(param))

	case "ROT13":
		ciphertext = crypto.ROT13Cifrar(message)

	case "BASE64":
		ciphertext = crypto.Base64Cifrar(message)

	case "XOR":
		fmt.Print("chave: ")
		key := t.readLine()
		if key == "" {
			fmt.Println("erro: chave vazia")
			return
		}
		ciphertext = crypto.XORCifrar(message, key)
		keyNote = "lembre de enviar a chave para o destinatario por fora"
	}

	packet := buildPacket(algo, param, ciphertext)

	fmt.Println()
	fmt.Println("-- codigo gerado --")
	fmt.Println(packet)
	fmt.Println("-------------------")
	if keyNote != "" {
		fmt.Println("obs:", keyNote)
	}
}

func (t *Terminal) flowDecifrar() {
	fmt.Print("codigo recebido: ")
	packet := t.readLine()
	if packet == "" {
		fmt.Println("erro: nenhum codigo fornecido")
		return
	}

	algo, param, ciphertext, err := parsePacket(packet)
	if err != nil {
		fmt.Println("erro:", err)
		return
	}

	var message string
	var decErr error

	switch algo {
	case "CESAR":
		message = crypto.CesarDecifrar(ciphertext, paramToShift(param))

	case "ROT13":
		message = crypto.ROT13Decifrar(ciphertext)

	case "BASE64":
		message, decErr = crypto.Base64Decifrar(ciphertext)
		if decErr != nil {
			fmt.Println("erro:", decErr)
			return
		}

	case "XOR":
		fmt.Print("chave: ")
		key := t.readLine()
		if key == "" {
			fmt.Println("erro: chave vazia")
			return
		}
		message, decErr = crypto.XORDecifrar(ciphertext, key)
		if decErr != nil {
			fmt.Println("erro:", decErr)
			return
		}
	}

	fmt.Println()
	fmt.Println("-- mensagem decifrada --")
	fmt.Println(message)
	fmt.Println("------------------------")
}

func (t *Terminal) chooseAlgorithm() (algo, param string) {
	fmt.Println()
	fmt.Println("algoritmo:")
	fmt.Println("1 - Cesar")
	fmt.Println("2 - ROT13")
	fmt.Println("3 - Base64")
	fmt.Println("4 - XOR")
	fmt.Print("> ")

	switch t.readLine() {
	case "1":
		fmt.Print("deslocamento (padrao 3): ")
		d := t.readLine()
		if d == "" || !isNumeric(d) {
			d = "3"
		}
		return "CESAR", d
	case "2":
		return "ROT13", ""
	case "3":
		return "BASE64", ""
	case "4":
		return "XOR", ""
	default:
		fmt.Println("opcao invalida")
		return "", ""
	}
}

func (t *Terminal) readLine() string {
	line, _ := t.reader.ReadString('\n')
	return strings.TrimSpace(line)
}

func buildPacket(algo, param, ciphertext string) string {
	if param != "" {
		return fmt.Sprintf("[%s:%s]%s", algo, param, ciphertext)
	}
	return fmt.Sprintf("[%s]%s", algo, ciphertext)
}

func parsePacket(packet string) (algo, param, ciphertext string, err error) {
	if !strings.HasPrefix(packet, "[") {
		return "", "", "", fmt.Errorf("formato invalido: deve comecar com '['")
	}
	end := strings.Index(packet, "]")
	if end == -1 {
		return "", "", "", fmt.Errorf("formato invalido: ']' nao encontrado")
	}

	parts := strings.SplitN(packet[1:end], ":", 2)
	algo = strings.ToUpper(parts[0])
	if len(parts) == 2 {
		param = parts[1]
	}
	ciphertext = packet[end+1:]

	switch algo {
	case "CESAR", "ROT13", "BASE64", "XOR":
	default:
		return "", "", "", fmt.Errorf("algoritmo desconhecido: %s", algo)
	}
	return algo, param, ciphertext, nil
}

func paramToShift(param string) int {
	n, err := strconv.Atoi(param)
	if err != nil || n == 0 {
		return 3
	}
	return n
}

func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return len(s) > 0
}
