package main

import (
	"fmt"
	"prawn/lexer"
	"prawn/parser"
	"time"
)

func main() {
	lexer := lexer.InitLexer(`
	var myBool = 250;
	write("Hola Mundo XDDDDSODSODSD SD-SOS-OS");`)

	// Crear el parser y pasarle el canal de tokens
	parser := parser.NewParser(lexer.TokenChan)

	// Parsear los tokens que vayan llegando

	/*no sabe el tipo de dato
	 */
	AST := parser.Parse()
	//aqui imprime el AST que retorna el Parser
	// esque como es un channel, para imprimir un chan es asi
	for {
		fmt.Println(<-AST)
		time.Sleep(1 * time.Second)
	}

}
