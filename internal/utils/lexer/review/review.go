package review

func IsLetter(currentChar byte) bool {
	return (currentChar >= 'a' && currentChar <= 'z') ||
		(currentChar >= 'A' && currentChar <= 'Z') ||
		currentChar == '_'
}

func IsDigit(currentChar byte) bool {
	return currentChar >= '0' && currentChar <= '9'
}

func IsSymbol(currentChar byte) bool {
	return currentChar == '=' || currentChar == '_' || currentChar == ')' || currentChar == '(' || currentChar == '+' || currentChar == '-' || currentChar == ';'
}
