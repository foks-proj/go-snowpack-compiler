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

var lexError error
var top *FileNode

type LexerError struct {
	Line     int
	Filename string
	Msg      string
}

func (e LexerError) Error() string {
	return fmt.Sprintf("%s:%d: %s", e.Filename, e.Line, e.Msg)
}

func (s *snowpLex) Error(es string) {
	lexError = LexerError{Msg: es, Filename: s.l.filename, Line: s.l.lineno}
}

func Parse(
	indat []byte,
	nm string,
) (
	*FileNode,
	error,
) {
	lexer := Lex(indat, nm)
	l := &snowpLex{l: lexer}
	ret := snowpParse(l)
	if ret != 0 {
		return nil, fmt.Errorf("parse error in file %s", nm)
	}
	if lexError != nil {
		return nil, lexError
	}
	return top, nil
}
