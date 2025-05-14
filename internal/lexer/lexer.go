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

// da un vistazo del proximo caracter sin mover el puntero del Lexer
func (lexer *Lexer) previewNextChar() rune {
	// retorna tipo de dato 'rune' que representa el proximo caracter sin mover el actual
	return rune(lexer.input[lexer.nextRead])
}

func (lexer *Lexer) Tokenizer() tokenspec.Token {
	//crea una variable token para cada token
	var token tokenspec.Token

	lexer.jumpWhitespaces()

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

	switch {
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
		token = tokenspec.NewToken(tokenspec.ILLEGAL, string(lexer.currentChar), lexer.position, lexer.currentLine, 1, 1)
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
	tokenType := tokenspec.LookupIdent(literal)
	return tokenspec.NewToken(tokenType, literal, start, lexer.currentLine, lexer.position-start, 0)
}

func (lexer *Lexer) readNumber() tokenspec.Token {
	start := lexer.position
	for review.IsDigit(lexer.currentChar) {
		lexer.readChar()
	}
	literal := lexer.input[start:lexer.position]
	return tokenspec.NewToken(tokenspec.INT, literal, start, lexer.currentLine, len(literal), 0)
}

func (lexer *Lexer) readSymbols(ch byte) tokenspec.Token {
	var token tokenspec.Token
	switch ch {
	case '-':
		token = tokenspec.NewToken(tokenspec.MINUS, string(lexer.currentChar), lexer.position, lexer.currentLine, 1, 0)
	case '+':
		token = tokenspec.NewToken(tokenspec.PLUS, string(lexer.currentChar), lexer.position, lexer.currentLine, 1, 0)
	case '(':
		token = tokenspec.NewToken(tokenspec.LPAREN, string(lexer.currentChar), lexer.position, lexer.currentLine, 1, 0)
	case ')':
		token = tokenspec.NewToken(tokenspec.RPAREN, string(lexer.currentChar), lexer.position, lexer.currentLine, 1, 0)
	case '=':
		if char := lexer.previewNextChar(); char == rune(lexer.currentChar) {
			token = tokenspec.NewToken(tokenspec.EQUAL, "==", lexer.position, lexer.currentLine, 2, 0)
			lexer.readChar()
		} else {
			token = tokenspec.NewToken(tokenspec.ASSIGN, string(lexer.currentChar), lexer.position, lexer.currentLine, 1, 0)
		}
	case '!':
		if char := lexer.previewNextChar(); char == '=' {
			token = tokenspec.NewToken(tokenspec.NOT_EQUAL, "!=", lexer.position, lexer.currentLine, 2, 0)
			lexer.readChar()
		} else {
			token = tokenspec.NewToken(tokenspec.BANG, string(lexer.currentChar), lexer.position, lexer.currentLine, 1, 0)
		}
	case '<':
		if char := lexer.previewNextChar(); char == '=' {
			token = tokenspec.NewToken(tokenspec.LESS_OR_EQUAL, "<=", lexer.position, lexer.currentLine, 2, 0)
			lexer.readChar()
		} else {
			token = tokenspec.NewToken(tokenspec.LESS, string(lexer.currentChar), lexer.position, lexer.currentLine, 1, 0)
		}
	case '>':
		if char := lexer.previewNextChar(); char == '=' {
			token = tokenspec.NewToken(tokenspec.GREATER_OR_EQUAL, ">=", lexer.position, lexer.currentLine, 2, 0)
			lexer.readChar()
		} else {
			token = tokenspec.NewToken(tokenspec.GREATER, string(lexer.currentChar), lexer.position, lexer.currentLine, 1, 0)
		}
	case ';':
		token = tokenspec.NewToken(tokenspec.SEMICOLON, string(lexer.currentChar), lexer.position, lexer.currentLine, 1, 0)
		lexer.currentLine++
	default:
		token = tokenspec.NewToken(tokenspec.ILLEGAL, string(lexer.currentChar), lexer.position, lexer.currentLine, 1, 1)
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
	literal := lexer.input[strStart:lexer.position]

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
	lexer := InitLexer(`\\#var nombre = 250 == 250;
	var nombre = 250 != 250;
	var nombre = 250 >= 250;
	var nombre = 250 <= 250;
	write("Hola Mundo XDDDDSODSODSD SD-SOS-OS");`)
	for tok := lexer.Tokenizer(); tok.Type != tokenspec.EOF; tok = lexer.Tokenizer() {
		fmt.Printf("Type: '%-7s' Literal: '%s' Position: '%d' Length: '%d' Line: '%d' Code: '%d'\n", tok.Type, tok.Literal, tok.Position, tok.Length, tok.Line, tok.Code)
		time.Sleep(1 * time.Second)
	}
	fmt.Scanln()
}
