package tokenizer

import (
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `let one = 1;
let two = 2;

let add = fn(a, b) {
	a + b;
};

let three = add(one, two);
!-*/5;
5 < 10 > 5;

if (5 < 10) {
	return true;
} else {
	return false;
}

10 == 10;
10 != 9;
`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{LET, "let"},
		{IDENT, "one"},
		{ASSIGN, "="},
		{INT, "1"},
		{SEMICOLON, ";"},
		{LET, "let"},
		{IDENT, "two"},
		{ASSIGN, "="},
		{INT, "2"},
		{SEMICOLON, ";"},
		{LET, "let"},
		{IDENT, "add"},
		{ASSIGN, "="},
		{FUNCTION, "fn"},
		{LPAREN, "("},
		{IDENT, "a"},
		{COMMA, ","},
		{IDENT, "b"},
		{RPAREN, ")"},
		{LBRACE, "{"},
		{IDENT, "a"},
		{PLUS, "+"},
		{IDENT, "b"},
		{SEMICOLON, ";"},
		{RBRACE, "}"},
		{SEMICOLON, ";"},
		{LET, "let"},
		{IDENT, "three"},
		{ASSIGN, "="},
		{IDENT, "add"},
		{LPAREN, "("},
		{IDENT, "one"},
		{COMMA, ","},
		{IDENT, "two"},
		{RPAREN, ")"},
		{SEMICOLON, ";"},
		{BANG, "!"},
		{MINUS, "-"},
		{ASTER, "*"},
		{SLASH, "/"},
		{INT, "5"},
		{SEMICOLON, ";"},
		{INT, "5"},
		{LT, "<"},
		{INT, "10"},
		{GT, ">"},
		{INT, "5"},
		{SEMICOLON, ";"},
		{IF, "if"},
		{LPAREN, "("},
		{INT, "5"},
		{LT, "<"},
		{INT, "10"},
		{RPAREN, ")"},
		{LBRACE, "{"},
		{RETURN, "return"},
		{TRUE, "true"},
		{SEMICOLON, ";"},
		{RBRACE, "}"},
		{ELSE, "else"},
		{LBRACE, "{"},
		{RETURN, "return"},
		{FALSE, "false"},
		{SEMICOLON, ";"},
		{RBRACE, "}"},
		{INT, "10"},
		{EQ, "=="},
		{INT, "10"},
		{SEMICOLON, ";"},
		{INT, "10"},
		{NOT_EQ, "!="},
		{INT, "9"},
		{SEMICOLON, ";"},
		{EOF, ""},
	}

	tokenizer := New(input)

	for _, tt := range tests {
		token := tokenizer.NextToken()

		if token.Type != tt.expectedType {
			t.Fatalf("token type wrong. got: %q, want: %q", token.Type, tt.expectedType)
		}

		if token.Literal != tt.expectedLiteral {
			t.Fatalf("token literal wrong. got: %q, want: %q", token.Literal, tt.expectedLiteral)
		}
	}
}
