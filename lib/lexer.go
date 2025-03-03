package lib

import (
	"regexp"
	"unicode/utf8"
)

type Lexer struct {
	input    string
	filename string

	start       int
	pos         int
	width       int
	lineno      int
	emitNewline bool

	savePoint       int
	savePointLineno int

	tokens  chan token
	chanEof bool
}

type token struct {
	typ TokenType
	val string
}

type TokenType int

const (
	TokenEOF TokenType = iota
	TokenError

	TokenLBracket
	TokenRBracket
)

type transitionType int

const (
	ttPop transitionType = iota
	ttPush
	ttSwitch
	ttEof
	ttNone
	ttKeep
	ttErr
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
		lineno:   1,
	}
}

const RuneEOF = -1

func (l *Lexer) nextRune() rune {
	if l.pos >= len(l.input) {
		return RuneEOF
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += w
	l.emitNewline = (r == '\n')
	if l.emitNewline {
		l.lineno++
	}
	return r
}

func initialState(l *Lexer) nextState {

	for {
		r := l.nextRune()
		if r == RuneEOF {
			return nextState{t: ttEof}
		}
		if isLetter(r) {
			l.backup()
			return nextState{t: ttPush, f: lexIdentifier}
		}
		if isDigit(r) {
			l.backup()
			return lexNumber(l)
		}
		switch r {
		case '\n':
			l.eat()
		case ' ', '\t', '\r':
			l.eat()
		case '-':
			return lexDash(l)
		case ';':
			l.emit(TokenSemicolon)
		case '.':
			l.emit(TokenDot)
		case ',':
			l.emit(TokenComma)
		case ':':
			l.emit(TokenColon)
		case '{':
			l.emit(TokenLBrace)
		case '}':
			l.emit(TokenRBrace)
		case '(':
			l.emit(TokenLParen)
		case ')':
			l.emit(TokenRParen)
		case '[':
			l.emit(TokenLBracket)
		case ']':
			l.emit(TokenRBracket)
		case '@':
			l.emit(TokenAt)
		case '=':
			l.emit(TokenEquals)
		case '/':
			return lexFrontSlash(l)
		case '"':
			return nextState{t: ttPush, f: lexDQuotedString}
		default:
			return nextState{t: ttErr}
		}
	}
}

func lexFrontSlash(l *Lexer) nextState {
	r := l.nextRune()
	switch r {
	case '/':
		return nextState{t: ttPush, f: lexCppStyleComment}
	case '*':
		return lexCStyleComment(l)
	}
	return nextState{t: ttErr}
}

func lexDQuotedString(l *Lexer) nextState {
	start := l.pos
	for {
		r := l.nextRune()
		switch r {
		case RuneEOF:
			return nextState{t: ttEof}
		case '\n':
			return nextState{t: ttErr}
		case '"':
			str := l.input[start : l.pos-1]
			l.tokens <- token{typ: TokenDQoutedString, val: str}
			l.start = l.pos
			return nextState{t: ttPop}
		}
	}
}

func lexCStyleComment(l *Lexer) nextState {

	r := l.nextRune()
	switch r {
	case '*':
		return nextState{t: ttPush, f: func(l *Lexer) nextState { return lexDocString(l, true) }}
	case RuneEOF:
		return nextState{t: ttErr}
	default:
		l.backup()
		return nextState{t: ttPush, f: func(l *Lexer) nextState { return lexDocString(l, false) }}
	}
}

func lexDocString(l *Lexer, emit bool) nextState {
	start := l.pos
	for {
		loopPos := l.pos
		r := l.nextRune()
		if r == RuneEOF {
			return nextState{t: ttEof}
		}
		if r != '*' {
			continue
		}
		r = l.nextRune()
		if r == RuneEOF {
			return nextState{t: ttEof}
		}
		if r == '/' {
			if emit {
				l.tokens <- token{typ: TokenDoc, val: l.input[start:loopPos]}
				l.start = l.pos
			}
			return nextState{t: ttPop}
		}
	}
}
func lexCppStyleComment(l *Lexer) nextState {
	for {
		r := l.nextRune()
		if r == RuneEOF || r == '\n' {
			return nextState{t: ttPop}
		}
	}
}

func lexDash(l *Lexer) nextState {
	l.markSavePoint()
	r := l.nextRune()
	if r == RuneEOF {
		return nextState{t: ttEof}
	}
	switch {
	case r == '>':
		l.emit(TokenArrow)
		return nextState{t: ttKeep}
	case isDigit(r):
		l.restoreSavePoint()
		return lexNumber(l)
	default:
		return nextState{t: ttErr}
	}
}

func (l *Lexer) eat() {
	l.start = l.pos
}

func Lex(
	data []byte,
	filename string,
) *Lexer {
	l := newLexer(data, filename)
	go l.run()
	return l
}

func (l *Lexer) next() token {
	if l.chanEof {
		return token{typ: TokenEOF}
	}
	ret, ok := <-l.tokens
	if !ok {
		l.chanEof = true
		return token{typ: TokenEOF}
	}
	if ret.typ == TokenEOF {
		l.chanEof = true
	}
	return ret
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
		case ttErr:
			l.emit(TokenError)
			state = nil
		case ttEof:
			l.emit(TokenEOF)
			state = nil
		case ttKeep:
			// noop, keep the current state
		}
	}
	close(l.tokens)
}

func (l *Lexer) backup() {
	if l.emitNewline {
		l.lineno--
	}
	l.pos -= l.width
}

func (l *Lexer) markSavePoint() {
	l.savePoint = l.pos
	l.savePointLineno = l.lineno
}

func (l *Lexer) restoreSavePoint() {
	l.pos = l.savePoint
	l.lineno = l.savePointLineno
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

var uint32rxx = regexp.MustCompile(`^0x[0-9a-fA-F]{8}$`)
var uint64rxx = regexp.MustCompile(`^0x[0-9a-fA-F]{16}$`)
var intrxx = regexp.MustCompile(`^-?[0-9]+$`)

func lexNumber(l *Lexer) nextState {
	for {
		r, w := utf8.DecodeRuneInString(l.input[l.pos:])
		if !isDigit(r) && r != '-' && r != 'x' && (r < 'a' || r > 'f') {
			break
		}
		l.pos += w
	}
	var typ TokenType
	switch {
	case uint64rxx.MatchString(l.txt()):
		typ = TokenUint64Val
	case uint32rxx.MatchString(l.txt()):
		typ = TokenUint32Val
	case intrxx.MatchString(l.txt()):
		typ = TokenIntVal
	default:
		return nextState{t: ttErr}
	}
	l.emit(typ)
	return nextState{t: ttKeep}
}

func (l *Lexer) emitIdentifier(start int) {
	txt := l.input[start:l.pos]
	var typ TokenType
	switch txt {
	case "typedef":
		typ = TokenTypedef
	case "List":
		typ = TokenList
	case "Option":
		typ = TokenOption
	case "Blob":
		typ = TokenBlob
	case "struct":
		typ = TokenStruct
	case "Text":
		typ = TokenText
	case "Uint":
		typ = TokenUint
	case "Int":
		typ = TokenInt
	case "Bool":
		typ = TokenBool
	case "enum":
		typ = TokenEnum
	case "variant":
		typ = TokenVariant
	case "case":
		typ = TokenCase
	case "switch":
		typ = TokenSwitch
	case "void":
		typ = TokenVoid
	case "default":
		typ = TokenDefault
	case "protocol":
		typ = TokenProtocol
	case "errors":
		typ = TokenErrors
	case "true":
		typ = TokenTrue
	case "false":
		typ = TokenFalse
	case "argHeader":
		typ = TokenArgHeader
	case "resHeader":
		typ = TokenResHeader
	case "import":
		typ = TokenImport
	case "go:import":
		typ = TokenGoImport
	case "ts:import":
		typ = TokenTypeScriptImport
	case "as":
		typ = TokenAs
	case "Future":
		typ = TokenFuture
	default:
		typ = TokenIdentifier
	}
	l.tokens <- token{typ: typ, val: txt}
	l.start = l.pos
}

func lexIdentifier(l *Lexer) nextState {
	start := l.pos
	for {
		r := l.nextRune()

		// There is one case where it's OK to have a colon in the
		// middle of an indentifier.
		isOkCol := false
		if r == ':' {
			txt := l.input[start:l.pos]
			if txt == "go:" || txt == "ts:" {
				isOkCol = true
			}
		}
		if !isLetter(r) && !isDigit(r) && !isOkCol && r != '_' {
			l.backup()
			break
		}
	}
	l.emitIdentifier(start)
	return nextState{t: ttPop}
}
