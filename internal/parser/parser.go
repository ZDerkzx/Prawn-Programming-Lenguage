package parser

import (
	"fmt"
	"prawn/lexer/tokenspec"
	"prawn/utils/lexer/review"
	"prawn/utils/parser/errors"
	"strconv"
)

type Parser struct {
	tokens   []tokenspec.Token
	position int
	errors   []string
}

// Node contiene todo tipo de datos
type Node interface{}

// contiene la declaracion de una variable tipo(Name "miVariable" Value: NumberExpr)
type VarDeclare[T any] struct {
	Name  string
	Value T
}

func (varDecl *VarDeclare[T]) Payload() map[string]map[string]interface{} {
	return map[string]map[string]interface{}{
		"Payload": {
			"Ident": varDecl.Name,
			"Value": varDecl.Value,
		},
	}
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

type WriteDecl struct {
	Value Node
}

func (wrt *WriteDecl) Payload() map[string]map[string]interface{} {
	return map[string]map[string]interface{}{
		"Payload": {
			"Value": wrt.Value,
		},
	}
}

// el constructor de la 'struct'
/*
El lexer le pasa al parser la informacion por un channel ?Como Funciona un Channel{
	al leer un indice se borra
	no puedes leer el proximo indice sin borrarlo
}
	asi que por eso decidi tener una variable para almacenarlo,primero lee
	todo el contenido del channel para poder usarlo como pegue la gana
*/
func NewParser(tokenchan chan tokenspec.Token) *Parser {
	var buffer []tokenspec.Token
	//va agregando token por token al buffer
	for tok := range tokenchan {
		//Si es EOF termina(osea el fin del codigo)
		if tok.Type == tokenspec.EOF {
			break
		}
		buffer = append(buffer, tok)
	}
	//retorna la 'struct'
	return &Parser{
		tokens:   buffer,
		position: 0,
		errors:   []string{},
	}
}

// verifica el proximo token sin saltarlo
func (parser *Parser) previewNextToken() tokenspec.Token {
	return parser.tokens[parser.position+1]
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

/*
Esta funcion Lee el tipo de Dato INT,STRING,IDENT y lo retorna en Un 'struct' Node
*/
//error esto regresa un Nil en un binaryexpr, no detecta que sea un int
func (parser *Parser) ParseExpressionType() Node {
	//crea un switch para validar el currentToken.Type
	switch parser.currentToken().Type {
	//si es INT crea el Nodo y lo pasa a Int por que esta en String(Necesita manejo de errores(ahorita nadamas por test))
	case tokenspec.INT:
		valueInt := Atoi(parser.currentToken().Literal)
		return NumberExpr{Value: valueInt}
	case tokenspec.STRING:
		valueStr := parser.currentToken().Literal
		return StringExpr{Value: valueStr}
	case tokenspec.IDENT:
		valueIDENT := parser.currentToken().Literal
		return VarExpr{Name: valueIDENT}
	// otros casos como booleanos, operaciones, etc.
	default:
		return nil
	}

}

// tipo constructor que crea 'BinaryExpr',left node,Operator string,right node
func (parser *Parser) CreateBinaryExpression(leftValue Node) *BinaryExpr {
	//Guarda el primer valor(lo guarda como literal)
	//avanza al siguiente token que seria el operador tipo (*,-,+)
	parser.NextToken()
	OpValue := parser.currentToken()
	parser.NextToken()
	rightValue := parser.ParseExpressionType()
	parser.NextToken()
	return &BinaryExpr{
		Left:  leftValue,
		Op:    OpValue.Literal,
		Right: rightValue,
	}
}

// Parsea write("Hola Mundo")
func (parser *Parser) ParseWriteDecl() *WriteDecl {
	/*current token 'write' asi que asemos un parser.NextToken
	y pasamos al siguiente token
	*/
	parser.NextToken()
	//el current token '(' y verificamos si existe si no agregamos el error y retornamos nil
	if parser.currentToken().Type != tokenspec.LPAREN {
		parser.errors = append(parser.errors, "Expected '(' but found '%s'.", parser.currentToken().Literal)
		return nil
	}
	//pasa al contenido
	parser.NextToken()
	//aqui guarda el contenido en un tipo de dato 'Node'
	writeContentValue := parser.ParseExpressionType()
	//pasa al siguiente token que deberia de ser ')'
	parser.NextToken()
	if parser.currentToken().Type != tokenspec.RPAREN {
		parser.errors = append(parser.errors, "Expected ')' but found '%s'.", parser.currentToken().Literal)
		return nil
	}
	parser.NextToken()
	if parser.currentToken().Type != tokenspec.SEMICOLON {
		parser.errors = append(parser.errors, "Expected ';' but found '%s'.", parser.currentToken().Literal)
		return nil
	}
	return &WriteDecl{
		Value: writeContentValue,
	}
}

func (parser *Parser) ParseVarDeclare() *VarDeclare[any] {
	//pasa al siguiente token 'IDENT' (nombre de la variable)
	parser.NextToken()
	if parser.currentToken().Type != tokenspec.IDENT {
		parser.errors = append(parser.errors, errors.CreateErrorExpected(parser.position, 2, tokenspec.IDENT, string(parser.currentToken().Type)))
		return nil
	}
	//aguarda el nombre de la variable
	varName := &VarExpr{Name: parser.currentToken().Literal}
	//pasa al siguiente token 'ASSIGN'(=)
	parser.NextToken()
	// si no encuentra el token tira error y lo almacena en un slice de errores
	if parser.currentToken().Type != tokenspec.ASSIGN {
		//hay que mejorar este mensaje de error
		parser.errors = append(parser.errors, errors.CreateErrorExpected(parser.position, 2, tokenspec.ASSIGN, string(parser.currentToken().Type)))
		fmt.Println(parser.errors)
		//no retorna nada
		return nil
	}
	/*si no se encontro ningun error sigue y lo proximo deberia ser
	el contenido
	*/
	parser.NextToken()
	if review.IsArithmeticSymbol(parser.previewNextToken().Literal) {
		varValueLeft := parser.ParseExpressionType()
		varValue := parser.CreateBinaryExpression(varValueLeft)
		if parser.currentToken().Type != tokenspec.SEMICOLON {
			fmt.Println("Error expected ';' xd")
			parser.errors = append(parser.errors, fmt.Sprintf("Expected ';' but found '%s'", parser.currentToken().Literal))
			return nil
		}
		/*avanza al siguiente token para no dejar al currentToken con el mismo si no se pone esto
		podria causar error*/
		parser.NextToken()
		return &VarDeclare[any]{
			Name:  varName.Name,
			Value: *varValue,
		}
	}
	varValue := parser.ParseExpressionType()
	parser.NextToken()
	if parser.currentToken().Type != tokenspec.SEMICOLON {
		parser.errors = append(parser.errors, fmt.Sprintf("Expected ';' but found '%s'", parser.currentToken().Literal))
		return nil
	}
	parser.NextToken()
	return &VarDeclare[any]{
		Name:  varName.Name,
		Value: varValue,
	}
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

// parsea Nodo por nodo
func (parser *Parser) parseNode() map[string]interface{} {
	//leer que tipo de token es
	//aqui vamos a leer los tipos de tokens
	currentToken := parser.currentToken()
	switch currentToken.Type {
	case tokenspec.VAR:
		node := parser.ParseVarDeclare()
		return map[string]interface{}{
			"VarDeclare": node.Payload(),
		}
	case tokenspec.WRITE:
		node := parser.ParseWriteDecl()
		return map[string]interface{}{
			"Write": node.Payload(),
		}
	default:
		parser.errors = append(parser.errors, fmt.Sprintf("Token '%s' no reconocido position '%d'", parser.currentToken().Literal, parser.position))
		return nil
	}
}

// Crea el AST(Abstract Sintaxys Tree)
func (parser *Parser) Parse() (map[string][]map[string]interface{}, []string) {
	var program []map[string]interface{}

	for parser.position < len(parser.tokens) {
		node := parser.parseNode()
		if node != nil {
			program = append(program, node)
		} else {
			// avanzar para no quedar en bucle infinito si hay error
			parser.NextToken()
		}
	}
	return map[string][]map[string]interface{}{
		"Program": program,
	}, parser.errors
}
