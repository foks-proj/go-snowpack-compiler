package lib

import "fmt"

//go:generate go tool golang.org/x/tools/cmd/goyacc -o parser.go -p "snowp" parser.y

type snowpLex struct {
	l *Lexer
}

func (s *snowpLex) Lex(yylval *snowpSymType) int {
	tok := s.l.next()
	yylval.rawval = tok.val
	return int(tok.typ)
}

var lexErr error
var parseErr error
var top *Root

type LexerError struct {
	Line     int
	Filename string
	Msg      string
}

func (e LexerError) Error() string {
	return fmt.Sprintf("%s:%d: %s", e.Filename, e.Line, e.Msg)
}

func (s *snowpLex) Error(es string) {
	lexErr = LexerError{Msg: es, Filename: s.l.filename, Line: s.l.lineno}
}

func Parse(
	indat []byte,
	nm string,
) (
	*Root,
	error,
) {
	lexer := Lex(indat, nm)
	l := &snowpLex{l: lexer}
	snowpParse(l)
	if parseErr != nil {
		return nil, parseErr
	}
	if lexErr != nil {
		return nil, lexErr
	}
	return top, nil
}

func init() {
	snowpErrorVerbose = true
}
