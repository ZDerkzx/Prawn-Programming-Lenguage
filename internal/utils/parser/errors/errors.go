package errors

import "fmt"

// errores para caracteres esperados pero enves de eso pusieron otros JSSJSJJSJSJSJ
func CreateErrorExpected(expected string, found string) string {
	return fmt.Sprintf("Expected '%s' but found '%s'.", expected, found)
}

func CreateErrorUnrecognizableToken(tokenLiteral string) string {
	return fmt.Sprintf("Token Unrecognizable '%s'.", tokenLiteral)
}
