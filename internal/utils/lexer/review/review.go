package review

import "prawn/lexer/tokenspec"

func IsLetter(currentChar byte) bool {
	return (currentChar >= 'a' && currentChar <= 'z') ||
		(currentChar >= 'A' && currentChar <= 'Z') ||
		currentChar == '_'
}

func IsDigit(currentChar byte) bool {
	return currentChar >= '0' && currentChar <= '9'
}

func IsSymbol(currentChar byte) bool {
	if _, ok := tokenspec.SymbolTokens[string(currentChar)]; ok {
		return true
	}
	return false
}
