package main

import (
	"fmt"
	"prawn/lexer/tokenspec"
	"prawn/utils/lexer/review"
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
	//salta los espacios en blanco
	lexer.jumpWhitespaces()
	startPos := lexer.position
	if review.IsLetter(lexer.currentChar) {
		//guarda en una variable el literal que contiene la sintax(no tokentype en si)
		token.Literal, token.Length = lexer.readIdentifier()
		// verifica si el literal existe y si existe lo guarda en token.Type
		token.Type = tokenspec.LookupIdent(token.Literal)
		token.Position = startPos
		token.Line = 0
		return token
	}

	if review.IsDigit(lexer.currentChar) {
		token.Type = tokenspec.INT
		token.Literal, token.Length = lexer.readNumber()
		token.Position = startPos
		token.Line = 0
		return token
	}
	switch lexer.currentChar {
	//el manejo de errores ahora esta en default
	case '=':
		token = tokenspec.NewToken(tokenspec.ASSIGN, lexer.currentChar, lexer.position, lexer.currentLine, 1)
	case '+':
		token = tokenspec.NewToken(tokenspec.PLUS, lexer.currentChar, lexer.position, lexer.currentLine, 1)
	case '-':
		token = tokenspec.NewToken(tokenspec.MINUS, lexer.currentChar, lexer.position, lexer.currentLine, 1)
	case '(':
		token = tokenspec.NewToken(tokenspec.LPAREN, lexer.currentChar, lexer.position, lexer.currentLine, 1)
	case ')':
		token = tokenspec.NewToken(tokenspec.RPAREN, lexer.currentChar, lexer.position, lexer.currentLine, 1)
	case '"':
		startPos := lexer.position
		token.Literal, token.Length = lexer.readString()
		token.Type = tokenspec.STRING
		token.Position = startPos
		token.Line = lexer.currentLine
	case ';':
		token = tokenspec.NewToken(tokenspec.SEMICOLON, lexer.currentChar, lexer.position, lexer.currentLine, 1)
		lexer.currentLine += 1
	case 0:
		token.Type = tokenspec.EOF
		token.Literal = ""
		token.Position = lexer.position
		token.Line = lexer.currentLine
	default:
		token = tokenspec.Token{Type: tokenspec.ILLEGAL, Literal: string(lexer.currentChar)}
		lexer.readChar()
	}
	lexer.readChar()
	return token
}

// salta todos los espacios
func (lexer *Lexer) jumpWhitespaces() {
	for lexer.currentChar == ' ' || lexer.currentChar == '\n' {
		lexer.readChar()
	}
}

// parte una porcion del codigo
func (lexer *Lexer) readIdentifier() (string, int) {
	/*guarda la posicion actual por ejemplo
	un "print" empieza, guarda en la variable 'start' el inicio osea 0 que
	contiene la letra 'p' y recorre todo dependiendo si es letra y al final
	como ya no da True retorna lexer.input cortado del punto de inicio al final de donde se quedo
	*/
	start := lexer.position
	for review.IsLetter(lexer.currentChar) || review.IsDigit(lexer.currentChar) {
		lexer.readChar()
	}
	//corta el pedazo del input
	return lexer.input[start:lexer.position], len(lexer.input[start:lexer.position])
}

func (lexer *Lexer) readNumber() (string, int) {
	start := lexer.position
	for review.IsDigit(lexer.currentChar) {
		lexer.readChar()
	}
	return lexer.input[start:lexer.position], len(lexer.input[start:lexer.position])
}

// esta funcion lee lo que hay dentro de las comillas y lo retorna como string ya que no es un tokentype
func (lexer *Lexer) readString() (string, int) {
	position := lexer.position + 1
	for {
		lexer.readChar()
		if lexer.currentChar == '"' {
			break
		}
	}
	/*corta de la posicion actual a la posicion final y lo retorna
	como string y tambien la longitud de la cadena de texto
	*/
	return lexer.input[position:lexer.position], len(lexer.input[position:lexer.position])
}

func main() {
	lexer := InitLexer(`var xx = "Hola Mundo";write("Hola Mundo");write(500 + 500;var nombre = "Pedro";)`)
	for tok := lexer.Tokenizer(); tok.Type != tokenspec.EOF; tok = lexer.Tokenizer() {
		fmt.Printf("Type: '%-7s' Literal: '%s' Position: '%d' Length: '%d' Line: '%d'\n", tok.Type, tok.Literal, tok.Position, tok.Length, tok.Line)
	}
}
