package review

/*para mientras lo tengo asi para testear luego se pasara esto a la libreria 'tokenspec' usando las value de los tokentypes reales
y separando por simbolos individuales
*/
var symbols = map[string]string{
	"!":  "BANG",
	"=":  "ASSIGN",
	"_":  "UNDERSCORE",
	")":  "RPAREN",
	"(":  "LPAREN",
	"+":  "PLUS",
	"-":  "MINUS",
	";":  "SEMICOLON",
	"!=": "NOT_EQUAL",
	"==": "EQUAL",
	"<":  "LESS_THAN",
	">":  "GREATER_THAN",
	"<=": "LESS_THAN_OR_EQUAL",
	">=": "GREATER_THAN_OR_EQUAL",
}

func IsLetter(currentChar byte) bool {
	return (currentChar >= 'a' && currentChar <= 'z') ||
		(currentChar >= 'A' && currentChar <= 'Z') ||
		currentChar == '_'
}

func IsDigit(currentChar byte) bool {
	return currentChar >= '0' && currentChar <= '9'
}

func IsSymbol(currentChar byte) bool {
	if _, ok := symbols[string(currentChar)]; ok {
		return true
	}
	return false
}

func main() {
	// "!=" en bytes es: '!' = 33, '=' = 61
	// Ejemplo de uso: verificar cada byte

	ok := IsSymbol(33)
	println(ok)

}
