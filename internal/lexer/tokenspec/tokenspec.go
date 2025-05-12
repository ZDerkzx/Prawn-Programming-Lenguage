package tokenspec

type TokenType string

type Token struct {
	Type     TokenType
	Literal  string
	Position int
	Line     int
	Length   int
}

const (
	ILLEGAL   = "ILLEGAL"
	EOF       = "EOF"
	IDENT     = "IDENT"
	INT       = "INT"
	STRING    = "STRING"
	LPAREN    = "LPAREN"
	RPAREN    = "RPAREN"
	ASSIGN    = "ASSIGN"
	PLUS      = "PLUS"
	MINUS     = "MINUS"
	SEMICOLON = "SEMICOLON"

	//keywords
	VAR   = "VAR"
	WRITE = "WRITE"
)

var Keywords = map[string]TokenType{
	"write": WRITE,
	"var":   VAR,
}

// funcion para crear el token que usa la 'struct' TokenType y 'c' que es el currentChar
func NewToken(tokenType TokenType, ch byte, position int, line int, length int) Token {
	//retorna el token con el tipo y el literal
	return Token{Type: tokenType, Literal: string(ch), Position: position, Line: line, Length: length}
}

func LookupIdent(ident string) TokenType {
	// verifica si el IDENT que le pasamos existe en keywords
	/*Que ase?
	1. Crea 2 variables una del token y otro de la confirmacion
	2. entra a el mapa 'keywords' a la clave ident
	3. guarda en token la key ident como TokenType
	4. Retorna el tokentype
	*/
	if token, exists := Keywords[ident]; exists {
		return token
	}
	// Si no existe retorna ident
	return IDENT
}
