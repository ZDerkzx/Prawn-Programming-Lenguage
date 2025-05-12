package main

import (
	"fmt"
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
}

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	IDENT   = "IDENT"
	INT     = "INT"
	STRING  = "STRING"
	LPAREN  = "LPAREN"
	RPAREN  = "RPAREN"
	ASSIGN  = "ASSIGN"
	PLUS    = "PLUS"
	MINUS   = "MINUS"

	//keywords
	VAR   = "VAR"
	WRITE = "WRITE"
)

var keywords = map[string]TokenType{
	"write": WRITE,
	"var":   VAR,
}

func LookupIdent(ident string) TokenType {
	// verifica si el IDENT que le pasamos existe en keywords
	/*Que ase?
	1. Crea 2 variables una del token y otro de la confirmacion
	2. entra a el mapa 'keywords' a la clave ident
	3. guarda en token la key ident como TokenType
	4. Retorna el tokentype
	*/
	if token, exists := keywords[ident]; exists {
		return token
	}
	// Si no existe retorna ident
	return IDENT
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

func (lexer *Lexer) Tokenizer() Token {
	//crea una variable token para cada token
	var token Token
	//salta los espacios en blanco
	lexer.jumpWhitespaces()

	if isLetter(lexer.currentChar) {
		//guarda en una variable el literal que contiene la sintax(no tokentype en si)
		literal := lexer.readIdentifier()
		// verifica si el literal existe y si existe lo guarda en token.Type
		token.Type = LookupIdent(literal)
		token.Literal = literal
		return token
	}

	if isDigit(lexer.currentChar) {
		token.Type = INT
		token.Literal = lexer.readNumber()
		return token
	}
	switch lexer.currentChar {
	//el manejo de errores ahora esta en default
	case '=':
		token = newToken(ASSIGN, lexer.currentChar)
	case '+':
		token = newToken(PLUS, lexer.currentChar)
	case '-':
		token = newToken(MINUS, lexer.currentChar)
	case '(':
		token = newToken(LPAREN, lexer.currentChar)
	case ')':
		token = newToken(RPAREN, lexer.currentChar)
	case '"':
		literal := lexer.readString()
		token.Type = STRING
		token.Literal = literal
	case 0:
		token.Type = EOF
		token.Literal = ""
	default:
		token = Token{Type: ILLEGAL, Literal: string(lexer.currentChar)}
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

// funcion para crear el token que usa la 'struct' TokenType y 'c' que es el currentChar
func newToken(tokenType TokenType, ch byte) Token {
	//retorna el token con el tipo y el literal
	return Token{Type: tokenType, Literal: string(ch)}
}

func isLetter(currentChar byte) bool {
	return (currentChar >= 'a' && currentChar <= 'z') ||
		(currentChar >= 'A' && currentChar <= 'Z') ||
		currentChar == '_'
}

func isDigit(currentChar byte) bool {
	return currentChar >= '0' && currentChar <= '9'
}

func (lexer *Lexer) readIdentifier() string {
	/*guarda la posicion actual por ejemplo
	un "print" empieza, guarda en la variable 'start' el inicio osea 0 que
	contiene la letra 'p' y recorre todo dependiendo si es letra y al final
	como ya no da True retorna lexer.input cortado del punto de inicio al final de donde se quedo
	*/
	start := lexer.position
	for isLetter(lexer.currentChar) || isDigit(lexer.currentChar) {
		lexer.readChar()
	}
	//corta el pedazo del input
	return lexer.input[start:lexer.position]
}

func (lexer *Lexer) readNumber() string {
	start := lexer.position
	for isDigit(lexer.currentChar) {
		lexer.readChar()
	}
	return lexer.input[start:lexer.position]
}

// esta funcion lee lo que hay dentro de las comillas y lo retorna como string ya que no es un tokentype
func (lexer *Lexer) readString() string {
	position := lexer.position + 1
	for {
		lexer.readChar()
		if lexer.currentChar == '"' {
			break
		}
	}
	//corta de la posicion actual a la posicion final y lo retorna como string
	return lexer.input[position:lexer.position]
}

func main() {
	lexer := InitLexer(`write("Hello World")`)
	for tok := lexer.Tokenizer(); tok.Type != EOF; tok = lexer.Tokenizer() {
		fmt.Printf("Type: %-7s Literal: %s\n", tok.Type, tok.Literal)
	}
}
