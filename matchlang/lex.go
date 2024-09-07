package matchlang

import "strings"

type TokenType int

const (
	CodeToken TokenType = iota
	SizeToken
	TextToken
	EqualsToken
	NotEqualsToken
	LiteralToken
)

type LexToken struct {
	Type  TokenType
	Value string
}

func lex(s string) []LexToken {
	result := []LexToken{}
	for _, word := range strings.Split(s, " ") {
		var token LexToken
		switch word {
		case "code":
			token = LexToken{Type: CodeToken}
		case "size":
			token = LexToken{Type: SizeToken}
		case "text":
			token = LexToken{Type: TextToken}
		case "=":
			token = LexToken{Type: EqualsToken}
		case "!=":
			token = LexToken{Type: NotEqualsToken}
		default:
			token = LexToken{Type: LiteralToken, Value: word}
		}
		result = append(result, token)
	}
	return result
}
