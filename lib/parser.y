
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
    stmts    []Statement
    stmt     Statement
    imprt    Importer
}

%type <file> top
%type <uniqueId> fileID uniqueID
%type <stmts> statements
%type <stmt> statement 
%type <imprt> import genericImport tsImport goImport

%token TokenAt TokenSemicolon TokenAs
%token TokenImport TokenTypeScriptImport TokenGoImport

%token <rawval> TokenUint64Val
%token <rawval> TokenDQoutedString TokenIdentifier

%%

top: 
    fileID statements
    {
        $$ = &FileNode{ Id: $1, Stmts : $2 }
        top = $$ // Set the global top variable
    }
    ;

statements:
    /* empty */ { $$ = []Statement{} }
    | statements statement
    {
        $$ = append($1, $2)
    }
    ;

genericImport: 
    TokenImport TokenDQoutedString TokenAs TokenIdentifier TokenSemicolon
    { 
        $$ = &GenericImport{ BaseImport : BaseImport { Path: $2, Name : $4 }  }
    }
    ;

tsImport: 
    TokenTypeScriptImport TokenDQoutedString TokenAs TokenIdentifier TokenSemicolon
    { 
        $$ = &TypeScriptImport{ BaseImport : BaseImport { Path: $2, Name : $4 } }
    } 
    ;

goImport: 
    TokenGoImport
    { $$ = nil } 
    ;

import: 
    genericImport { $$ = $1 }
    | tsImport    { $$ = $1 }
    | goImport    { $$ = $1 }
    ;

statement:
    import { $$ = $1 }
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

const parserEOF = 0

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

func (s *snowpLex) Error(es string) {
    fmt.Printf("Lexer error: %s\n", es)
    lexError = fmt.Errorf("Lexer error: %s", es)
}
