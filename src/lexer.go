package main

func lexer(code string) []string {
	var lexed []string
	var var_name []rune

	flush := func() {
		if len(var_name) != 0 {
			lexed = append(lexed, string(var_name))
			var_name = nil
		}
	}

	for _, char := range code {
		switch char {
		case '(', ')', '[', ']', ',', '\'', ':', '=':
			flush()
			lexed = append(lexed, string(char))
		case '\r', ' ', '\t', '\n':
			flush()
			//ignore \r to stop potential problems on windows
		default:
			var_name = append(var_name, char)
		}
	}

	flush()

	return lexed
}
