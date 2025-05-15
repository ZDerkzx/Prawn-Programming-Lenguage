package main

import (
	"prawn/lexer"
	"prawn/parser"
)

func main() {
	lexer := lexer.InitLexer(`
	var myBool = 250 != 900
	if(myBool){
		write("Es true")
	};
	write("Hola Mundo XDDDDSODSODSD SD-SOS-OS");`)

	// Crear el parser y pasarle el canal de tokens
	parser := parser.NewParser(lexer.TokenChan)

	// Parsear los tokens que vayan llegando
	parser.Parse()
}
