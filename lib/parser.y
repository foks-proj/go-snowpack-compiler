
%{
package lib

import (
    "fmt"
)

%}

%union {
    file     *FileNode
    uniqueId UniqueID
    uint     uint64
    rawval   string
}

%type <file> top
%type <uniqueId> fileID uniqueID

%token TokenAt TokenSemicolon

%token <rawval> TokenUint64Val

%%

top: 
    fileID
    {
        $$ = &FileNode{ Id: $1 }
    }
    ;

fileID:
    uniqueID TokenSemicolon
    {
        $$ = $1
    }
    ;

uniqueID:
    TokenAt
    TokenUint64Val
    {
        $$ = UniqueID{ Val: $2 }
    }
    ;

%%

type snowpLex struct {
    l *Lexer
}

func (s *snowpLex) Lex(lval *snowpSymType) int {
    s.l.next()
    return 0
}

func (s *snowpLex) Error(es string) {
    fmt.Printf("Lexer error: %s\n", es)
}
