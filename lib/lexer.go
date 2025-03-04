package lib

import "unicode/utf8"

type Lexer struct {
	input    string
	filename string

	start int
	pos   int
	width int

	tokens chan token
}

type token struct {
	typ TokenType
	val string
}

type TokenType int

const (
	TokenEOF TokenType = iota
	TokenError

	TokenAt
	TokenUint64Val
	TokenUint32Val

	TokenTypedef
	TokenList
	TokenOption
	TokenBlob
	TokenStruct
	TokenText
	TokenUint
	TokenInt
	TokenBool

	TokenEnum
	TokenVariant
	TokenCase
	TokenSwitch
	TokenVoid
	TokenDefault

	TokenFalse
	TokenTrue

	TokenProtocol
	TokenErrors
	TokenArgHeader
	TokenResHeader

	TokenImport
	TokenGoImport
	TokenTypeScriptImpot

	TokenAs
	TokenFuture

	TokenArrow
	TokenSemicolon
	TokenDot
	TokenComma
	TokenColon
	TokenLBrace
	TokenRBrace
	TokenLParen
	TokenRParen
	TokenLBracket
	TokenRBracket

	TokenEquals

	TokenIdentifier
	TokenIntVal

	TokenDQoutedString
	TokenDocFrag
)

type transitionType int

const (
	ttPop transitionType = iota
	ttPush
	ttSwitch
	ttEof
	ttNone
)

type nextState struct {
	t transitionType
	f stateFn
}

type stateFn func(*Lexer) nextState

func newLexer(
	data []byte,
	filename string,
) *Lexer {
	return &Lexer{
		input:    string(data),
		filename: filename,
		tokens:   make(chan token),
	}
}

func initialState(l *Lexer) nextState {

	for {
		if l.pos >= len(l.input) {
			l.emit(TokenEOF)
			return nextState{t: ttEof}
		}
		r, w := utf8.DecodeRuneInString(l.input[l.pos:])
		l.width = w
		l.pos += w

		if isLetter(r) {
			l.backup()
			return lexIdentifier(l)
		}
		if isDigit(r) || r == '-' {
			l.backup()
			return lexNumber(l)
		}
	}
	return nextState{t: ttEof}
}

func Lex(
	data []byte,
	filename string,
) *Lexer {
	l := newLexer(data, filename)
	go l.run()
	return l
}

func (l *Lexer) run() {
	var stack []stateFn
	for state := initialState; state != nil; {
		ns := state(l)
		switch ns.t {
		case ttPop:
			if len(stack) == 0 {
				state = nil
			} else {
				state = stack[len(stack)-1]
				stack = stack[:len(stack)-1]
			}
		case ttPush:
			stack = append(stack, state)
			state = ns.f
		case ttSwitch:
			state = ns.f
		case ttEof:
			state = nil
		}
	}
	close(l.tokens)
}

func (l *Lexer) backup() {
	l.pos -= l.width
}

func (l *Lexer) emit(t TokenType) {
	l.tokens <- token{typ: t, val: l.txt()}
	l.start = l.pos
}

func isLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func (l *Lexer) txt() string {
	return l.input[l.start:l.pos]
}

func lexNumber(l *Lexer) nextState {
	for {
		r, w := utf8.DecodeRuneInString(l.input[l.pos:])
		if !isDigit(r) && r != '-' {
			break
		}
		l.pos += w
	}
	l.emit(TokenIntVal)
	return nextState{t: ttPop}
}

func (l *Lexer) emitIdentifier() {
	switch l.txt() {
	case "typedef":
		l.emit(TokenTypedef)
	case "List":
		l.emit(TokenList)
	case "Option":
		l.emit(TokenOption)
	case "Blob":
		l.emit(TokenBlob)
	case "struct":
		l.emit(TokenStruct)
	case "Text":
		l.emit(TokenText)
	case "Uint":
		l.emit(TokenUint)
	case "Int":
		l.emit(TokenInt)
	case "Bool":
		l.emit(TokenBool)
	case "enum":
		l.emit(TokenEnum)
	case "variant":
		l.emit(TokenVariant)
	case "case":
		l.emit(TokenCase)
	case "switch":
		l.emit(TokenSwitch)
	case "void":
		l.emit(TokenVoid)
	case "default":
		l.emit(TokenDefault)
	case "protocol":
		l.emit(TokenProtocol)
	case "errors":
		l.emit(TokenErrors)
	case "true":
		l.emit(TokenTrue)
	case "false":
		l.emit(TokenFalse)
	case "argHeader":
		l.emit(TokenArgHeader)
	case "resHeader":
		l.emit(TokenResHeader)
	case "import":
		l.emit(TokenImport)
	case "go:import":
		l.emit(TokenGoImport)
	case "ts:import":
		l.emit(TokenTypeScriptImpot)
	case "as":
		l.emit(TokenAs)
	case "Future":
		l.emit(TokenFuture)
	default:
		l.emit(TokenIdentifier)
	}
}

func lexIdentifier(l *Lexer) nextState {
	for {
		r, w := utf8.DecodeRuneInString(l.input[l.pos:])

		// There is one case where it's OK to have a colon in the
		// middle of an indentifier.
		isOkCol := false
		if r == ':' && (l.txt() == "go" || l.txt() == "ts") {
			isOkCol = true
		}
		if !isLetter(r) && !isDigit(r) && isOkCol {
			break
		}
		l.pos += w
	}
	l.emitIdentifier()
	return nextState{t: ttPop}
}
