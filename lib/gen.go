package lib

import "fmt"

//go:generate go tool golang.org/x/tools/cmd/goyacc -o parser.go -p "snowp" parser.y

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
