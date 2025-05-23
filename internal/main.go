package main

import (
	"fmt"
	"prawn/lexer"
	"prawn/parser"
)

func main() {
	lexer := lexer.InitLexer(`var test = 200;
	var myxd = 400;
	write("Hola Mundo XDDDDSODSODSD-A2-2");`)

	// Crear el parser y pasarle el canal de tokens
	parser := parser.NewParser(lexer.TokenChan)

	// Parsear los tokens que vayan llegando

	/*no sabe el tipo de dato
	 */
	AST, errors := parser.Parse()
	//aqui imprime el AST que retorna el Parser
	// esque como es un channel, para imprimir un chan es asis
	for node := range AST {
		fmt.Println(node)
	}
	fmt.Println("Errors: ", errors)
}
