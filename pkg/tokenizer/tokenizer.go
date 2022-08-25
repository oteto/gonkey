package tokenizer

import "github.com/oteto/gonkey/pkg/token"

type Tokenizer struct {
	input        string
	position     int  // 現在の読み込み位置
	readPosition int  // 次の読み込み位置
	char         byte // 読み込んだ文字（unicode にはひとまず対応しない）
}

// 現在の読み込み文字のトークンを取得し、読み込み位置を進める
func (t *Tokenizer) NextToken() token.Token {
	var tkn token.Token

	t.skipWhiteSpace()

	switch t.char {
	case '=':
		if t.peekChar() == '=' {
			tkn = t.makeTwoCharToken(token.EQ)
			break
		}
		tkn = token.NewToken(token.ASSIGN, t.char)
	case '+':
		tkn = token.NewToken(token.PLUS, t.char)
	case '-':
		tkn = token.NewToken(token.MINUS, t.char)
	case '*':
		tkn = token.NewToken(token.ASTER, t.char)
	case '/':
		tkn = token.NewToken(token.SLASH, t.char)
	case '!':
		if t.peekChar() == '=' {
			tkn = t.makeTwoCharToken(token.NOT_EQ)
			break
		}
		tkn = token.NewToken(token.BANG, t.char)
	case '<':
		tkn = token.NewToken(token.LT, t.char)
	case '>':
		tkn = token.NewToken(token.GT, t.char)
	case ',':
		tkn = token.NewToken(token.COMMA, t.char)
	case ';':
		tkn = token.NewToken(token.SEMICOLON, t.char)
	case '(':
		tkn = token.NewToken(token.LPAREN, t.char)
	case ')':
		tkn = token.NewToken(token.RPAREN, t.char)
	case '{':
		tkn = token.NewToken(token.LBRACE, t.char)
	case '}':
		tkn = token.NewToken(token.RBRACE, t.char)
	case 0:
		tkn.Type = token.EOF
		tkn.Literal = ""
	case '"':
		tkn.Type = token.STRING
		tkn.Literal = t.readString()
	case '[':
		tkn = token.NewToken(token.LBRACKET, t.char)
	case ']':
		tkn = token.NewToken(token.RBRACKET, t.char)
	case ':':
		tkn = token.NewToken(token.COLON, t.char)
	default:
		if isLetter(t.char) {
			tkn.Literal = t.readIdentifer()
			tkn.Type = token.LookUpIndent(tkn.Literal)
			return tkn
		} else if isDigit(t.char) {
			tkn.Literal = t.readNumber()
			tkn.Type = token.INT
			return tkn
		} else {
			tkn = token.NewToken(token.ILLEGAL, t.char)
		}
	}

	t.readChar()
	return tkn
}

// 読み込み位置を１つ進める
// 終端の場合は char に 0(NUL) がセットされる
func (t *Tokenizer) readChar() {
	if t.readPosition >= len(t.input) {
		// NUL 文字（まだ読み込んでいない or ファイルの終端）
		t.char = 0
	} else {
		t.char = t.input[t.readPosition]
	}
	t.position = t.readPosition
	t.readPosition += 1
}

// 識別子の終端まで読み進め、識別子のリテラルを返す
func (t *Tokenizer) readIdentifer() string {
	start := t.position
	for isLetter(t.char) {
		t.readChar()
	}
	return t.input[start:t.position]
}

// 整数値を終端まで読み進め、リテラルを返す
func (t *Tokenizer) readNumber() string {
	start := t.position
	for isDigit(t.char) {
		t.readChar()
	}
	return t.input[start:t.position]
}

func (t *Tokenizer) readString() string {
	start := t.position + 1 // " の次 "foo" なら f
	for {
		t.readChar()
		if t.char == '"' || t.char == 0 {
			break
		}
	}
	return t.input[start:t.position]
}

// 空白、改行は無視して読み進める
func (t *Tokenizer) skipWhiteSpace() {
	for t.char == ' ' || t.char == '\t' || t.char == '\n' || t.char == '\r' {
		t.readChar()
	}
}

// 次の読み込み位置の文字を取得する。
// 読み込み位置は進めない
func (t *Tokenizer) peekChar() byte {
	if t.readPosition >= len(t.input) {
		return 0
	}
	return t.input[t.readPosition]
}

// ２文字トークンを作成する
// １文字目を読んだ状態でコールする
func (t *Tokenizer) makeTwoCharToken(tokenType token.TokenType) token.Token {
	char := t.char
	t.readChar()
	return token.Token{
		Type:    tokenType,
		Literal: string(char) + string(t.char),
	}
}

// トークナイザを作成する
func New(input string) *Tokenizer {
	tokenizer := &Tokenizer{input: input}
	tokenizer.readChar()
	return tokenizer
}

// 識別子に使用できる文字かどうかをチェック
// 英字、アンスコ(_)を許可
func isLetter(char byte) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_'
}

// 数値チェック
func isDigit(char byte) bool {
	return '0' <= char && char <= '9'
}
