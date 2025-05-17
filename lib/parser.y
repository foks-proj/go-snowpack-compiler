
%{
package lib

%}

%union {
    file     *FileNode
    uniqueId UniqueID
    uint     uint64
    rawval   string
    stmts    []Statement
    stmt     Statement
    imprt    Importer
    dec      *Decorators
    doc      Docstring
    docRaw   string
}

%type <file> top
%type <uniqueId> fileID uniqueID
%type <stmts> statements
%type <stmt> statement typedef
%type <imprt> import genericImport tsImport goImport
%type <dec> decorators
%type <doc> doc 
%type <docRaw> docRaw

%token TokenAt TokenSemicolon TokenAs
%token TokenImport TokenTypeScriptImport TokenGoImport

%token <rawval> TokenUint64Val
%token <rawval> TokenDQoutedString TokenIdentifier TokenDoc TokenTypedef

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
    TokenGoImport TokenDQoutedString TokenAs TokenIdentifier TokenSemicolon
    { 
        $$ = &GoImport { BaseImport : BaseImport { Path: $2, Name : $4 } }
    } 
    ;

doc : 
    docRaw
    {
        $$ = Docstring{ Raw: $1 }
    }
    ;

docRaw
    : { $$ = "" }
    | docRaw TokenDoc
    {
        $$ = $1 + $2;
    }
    ;


decorators:
    doc { $$ = &Decorators{ Doc: $1 } }
    ;

typedef:
    decorators TokenTypedef 
    {
        $$ = &Typedef{}
    }
    ;

import: 
    genericImport { $$ = $1 }
    | tsImport    { $$ = $1 }
    | goImport    { $$ = $1 }
    ;

statement:
    import { $$ = $1 }
    | typedef { $$ = $1 }
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