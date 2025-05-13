package main

import (
	"fmt"
	"prawn/lexer/tokenspec"
	"prawn/utils/lexer/review"
	"time"
)

/*
	structura del Lexer 1. Input completo, 2.Posicion actual,

3. Proxima Posicion a Leer,
4. El caracter actual,
*/
type Lexer struct {
	input       string
	position    int
	nextRead    int
	currentChar byte
	currentLine int
}

// constructor
func InitLexer(input string) *Lexer {
	Lexer := &Lexer{input: input}
	//empieza a leer el primer caracter
	Lexer.readChar()
	//retorna la 'struct'
	return Lexer
}

/*
readChar funcion fundamental para avanzar con la lectura del codigo,
lee byte por byte(Letra por letra)
*/
func (lexer *Lexer) readChar() {
	// si llego al final de input resetea el currentChar a '0'
	if lexer.nextRead >= len(lexer.input) {
		lexer.currentChar = 0
	} else {
		// de lo contrario sigue
		lexer.currentChar = byte(lexer.input[lexer.nextRead])
	}
	/*actualizamos los datos
	lexer.position(la posicion actual) lo actualizamos al lexer.nextRead
	que es el siguiente para Leer
	*/
	lexer.position = lexer.nextRead
	lexer.nextRead++
}

func (lexer *Lexer) Tokenizer() tokenspec.Token {
	//crea una variable token para cada token
	var token tokenspec.Token

	if lexer.currentChar == 0 {
		return tokenspec.Token{
			Type:     tokenspec.EOF,
			Literal:  "",
			Position: lexer.position,
			Line:     lexer.currentLine,
			Length:   0,
		}
	}

	//salta los espacios en blanco
	lexer.jumpWhitespaces()
	switch {
	//el manejo de errores ahora esta en default
	case lexer.currentChar == '"':
		token = lexer.readStringToken()
	case review.IsLetter(lexer.currentChar):
		generatedToken := lexer.readIdentifier()
		token = generatedToken
	case review.IsDigit(lexer.currentChar):
		generatedToken := lexer.readNumber()
		token = generatedToken
	case review.IsSymbol(lexer.currentChar):
		generatedToken := lexer.readSymbols(lexer.currentChar)
		token = generatedToken
	default:
		token = tokenspec.Token{Type: tokenspec.ILLEGAL, Literal: string(lexer.currentChar)}
		lexer.readChar()
	}
	return token
}

// salta todos los espacios
func (lexer *Lexer) jumpWhitespaces() {
	for lexer.currentChar == ' ' || lexer.currentChar == '\n' || lexer.currentChar == '\t' {
		lexer.readChar()
	}
}

// parte una porcion del codigo
func (lexer *Lexer) readIdentifier() tokenspec.Token {
	start := lexer.position
	for review.IsLetter(lexer.currentChar) || review.IsDigit(lexer.currentChar) {
		lexer.readChar()
	}
	literal := lexer.input[start:lexer.position]
	return tokenspec.Token{
		Type:     tokenspec.LookupIdent(literal),
		Literal:  literal,
		Length:   lexer.position - start,
		Line:     lexer.currentLine,
		Position: start,
	}
}

func (lexer *Lexer) readNumber() tokenspec.Token {
	var token tokenspec.Token
	start := lexer.position
	for review.IsDigit(lexer.currentChar) {
		lexer.readChar()
	}
	token.Type = tokenspec.INT
	token.Literal = lexer.input[start:lexer.position]
	token.Length = len(lexer.input[start:lexer.position])
	token.Line = lexer.currentLine
	token.Position = start
	return token
}

func (lexer *Lexer) readSymbols(ch byte) tokenspec.Token {
	start := lexer.position
	token := tokenspec.Token{
		Literal:  string(ch),
		Position: start,
		Line:     lexer.currentLine,
		Length:   1,
	}

	switch ch {
	case '-':
		token.Type = tokenspec.MINUS
	case '+':
		token.Type = tokenspec.PLUS
	case '(':
		token.Type = tokenspec.LPAREN
	case ')':
		token.Type = tokenspec.RPAREN
	case '=':
		token.Type = tokenspec.ASSIGN
	case ';':
		token.Type = tokenspec.SEMICOLON
		lexer.currentLine++
	default:
		token.Type = tokenspec.ILLEGAL
	}

	lexer.readChar() // avanzar solo una vez
	return token
}

// esta funcion lee lo que hay dentro de las comillas y lo retorna como string ya que no es un tokentype
func (lexer *Lexer) readStringToken() tokenspec.Token {
	start := lexer.position // guarda la posición de la comilla inicial
	lexer.readChar()        // salta la comilla de apertura

	strStart := lexer.position
	for lexer.currentChar != '"' && lexer.currentChar != 0 {
		lexer.readChar()
	}
	literal := lexer.input[strStart:lexer.position] // sin comillas

	token := tokenspec.Token{
		Type:     tokenspec.STRING,
		Literal:  literal,
		Position: start,
		Length:   lexer.position - strStart,
		Line:     lexer.currentLine,
	}

	lexer.readChar() // salta la comilla final

	return token
}

func main() {
	fmt.Println(`el codigo enviado fue este: 
	var nombre = 235;
	write(nombre)`)
	fmt.Println("----TOKENS CREADOS----")
	lexer := InitLexer(`
	var nombre = 230;
	write("Hola Mundo XDDDDSODSODSD SD-SOS-OS");`)
	for tok := lexer.Tokenizer(); tok.Type != tokenspec.EOF; tok = lexer.Tokenizer() {
		fmt.Printf("Type: '%-7s' Literal: '%s' Position: '%d' Length: '%d' Line: '%d'\n", tok.Type, tok.Literal, tok.Position, tok.Length, tok.Line)
		time.Sleep(1 * time.Second)

	}
	fmt.Scanln()
}
