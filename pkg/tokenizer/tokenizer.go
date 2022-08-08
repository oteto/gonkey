package tokenizer

type Tokenizer struct {
	input        string
	position     int  // 現在の読み込み位置
	readPosition int  // 次の読み込み位置
	char         byte // 読み込んだ文字（unicode にはひとまず対応しない）
}

// 現在の読み込み文字のトークンを取得し、読み込み位置を進める
func (t *Tokenizer) NextToken() Token {
	var token Token

	t.skipWhiteSpace()

	switch t.char {
	case '=':
		if t.peekChar() == '=' {
			token = t.makeTwoCharToken(EQ)
			break
		}
		token = newToken(ASSIGN, t.char)
	case '+':
		token = newToken(PLUS, t.char)
	case '-':
		token = newToken(MINUS, t.char)
	case '*':
		token = newToken(ASTER, t.char)
	case '/':
		token = newToken(SLASH, t.char)
	case '!':
		if t.peekChar() == '=' {
			token = t.makeTwoCharToken(NOT_EQ)
			break
		}
		token = newToken(BANG, t.char)
	case '<':
		token = newToken(LT, t.char)
	case '>':
		token = newToken(GT, t.char)
	case ',':
		token = newToken(COMMA, t.char)
	case ';':
		token = newToken(SEMICOLON, t.char)
	case '(':
		token = newToken(LPAREN, t.char)
	case ')':
		token = newToken(RPAREN, t.char)
	case '{':
		token = newToken(LBRACE, t.char)
	case '}':
		token = newToken(RBRACE, t.char)
	case 0:
		token.Type = EOF
		token.Literal = ""
	default:
		if isLetter(t.char) {
			token.Literal = t.readIdentifer()
			token.Type = LookUpIndent(token.Literal)
			return token
		} else if isDigit(t.char) {
			token.Literal = t.readNumber()
			token.Type = INT
			return token
		} else {
			token = newToken(ILLEGAL, t.char)
		}
	}

	t.readChar()
	return token
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
func (t *Tokenizer) makeTwoCharToken(tokenType TokenType) Token {
	char := t.char
	t.readChar()
	return Token{tokenType, string(char) + string(t.char)}
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
