package parser

import (
	"fmt"
	"prawn/lexer/tokenspec"
	"strconv"
)

type Parser struct {
	tokens   []tokenspec.Token
	position int
	errors   []string
}

type Node interface{}

// contiene la declaracion de una variable tipo(Name "miVariable" Value: NumberExpr)
type VarDeclare struct {
	Name  string
	Value Node
}

// contiene una expresion tipo(x + 12)
type BinaryExpr struct {
	Left  Node
	Op    string
	Right Node
}

// contiene numeros
type NumberExpr struct {
	Value int
}

type StringExpr struct {
	Value string
}

// contiene el nombre de variables
type VarExpr struct {
	Name string
}

func NewParser(tokenchan chan tokenspec.Token) *Parser {
	var buffer []tokenspec.Token
	for tok := range tokenchan {
		buffer = append(buffer, tok)
		if tok.Type == tokenspec.EOF {
			break
		}
	}
	return &Parser{
		tokens:   buffer,
		position: 0,
		errors:   []string{},
	}
}

// regresa el token actual sin modificar el parser.position🗣️
func (parser *Parser) currentToken() tokenspec.Token {
	return parser.tokens[parser.position]
}

// version simplificada sin manejo de errores para donde no se necesita del manejo de errores
func Atoi(input string) int {
	value, _ := strconv.Atoi(input)
	return value
}

func (parser *Parser) ParseExpressionType() Node {
	switch parser.currentToken().Type {
	case tokenspec.INT:
		return &NumberExpr{Value: Atoi(parser.currentToken().Literal)}
	case tokenspec.STRING:
		return &StringExpr{Value: parser.currentToken().Literal}
	case tokenspec.IDENT:
		return &VarExpr{Name: parser.currentToken().Literal}
	// otros casos como booleanos, operaciones, etc.
	default:
		return nil
	}
}

func (parser *Parser) ParseVarDeclare() *VarDeclare {
	if parser.currentToken().Type == tokenspec.VAR {
		parser.NextToken()
		varName := &VarExpr{Name: parser.currentToken().Literal}
		parser.NextToken()

		if parser.currentToken().Type != tokenspec.ASSIGN {
			parser.errors = append(parser.errors, fmt.Sprintf("Expected '=' but found '%s'", parser.currentToken().Literal))
			return nil
		}
		parser.NextToken()
		varValue := parser.ParseExpressionType()

		return &VarDeclare{
			Name:  varName.Name,
			Value: varValue,
		}
	}
	return nil
}

func (parser *Parser) NextToken() tokenspec.Token {
	if parser.position < len(parser.tokens) {
		token := parser.tokens[parser.position]
		parser.position++
		return token
	}
	return tokenspec.Token{
		Type: tokenspec.EOF,
	}
}

func (parser *Parser) Parse() {
	for i := 0; i < len(parser.tokens); i++ {
		fmt.Printf("Index: '%d', Literal: '%s'\n", i, parser.tokens[i].Literal)
	}
}
