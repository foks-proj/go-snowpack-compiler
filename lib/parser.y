
%{
package lib

import (
    "fmt"
    "strconv"
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
    dec      Decorators
    doc      Docstring
    docRaw   string
    ident    Identifier
    typ      Type
    count    int
}

%type <file> top
%type <uniqueId> fileID uniqueID uniqueIDOpt
%type <stmts> statements
%type <stmt> statement typedef
%type <imprt> import genericImport tsImport goImport
%type <dec> decorators
%type <doc> doc 
%type <docRaw> docRaw
%type <ident> identifier
%type <typ> list type simpleType typeOrFuture blob dottedIdentifier future
%type <count> countOpt

%token TokenAt TokenSemicolon TokenAs TokenEquals TokenDot
%token TokenImport TokenTypeScriptImport TokenGoImport
%token TokenList TokenLParen TokenRParen TokenText TokenUint TokenInt TokenBool TokenBlob TokenFuture

%token <rawval> TokenUint64Val TokenIntVal
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
        $$ = GenericImport{ BaseImport : BaseImport { Path: $2, Name : $4 }  }
    }
    ;

tsImport: 
    TokenTypeScriptImport TokenDQoutedString TokenAs TokenIdentifier TokenSemicolon
    { 
        $$ = TypeScriptImport{ BaseImport : BaseImport { Path: $2, Name : $4 } }
    } 
    ;

goImport: 
    TokenGoImport TokenDQoutedString TokenAs TokenIdentifier TokenSemicolon
    { 
        $$ = GoImport { BaseImport : BaseImport { Path: $2, Name : $4 } }
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
    doc { $$ = Decorators{ Doc: $1 } }
    ;

uniqueIDOpt
    : { $$ = UniqueID{} }
    | uniqueID
    {
        $$ = $1;
    }
    ;

list:
    TokenList TokenLParen type TokenRParen
    {
        $$ = List{ Type: $3 }
    }
    ;

countOpt:
    /* empty */ { $$ = 0 }
    | TokenLParen TokenIntVal TokenRParen
    {
        var i int
        i, err := strconv.Atoi($2)
        if err != nil {
            parseErr = err
        } else if i <= 0 {
            parseErr = fmt.Errorf("blob byte-count must be greater than 0")
        } else {
            $$ = i
        }
    }
    ;

blob
    : TokenBlob countOpt
    {
        $$ = Blob{ Count: $2 }
    }
    ;

dottedIdentifier:
    identifier
    {
        $$ = DerivedType{ Base : $1 }
    }
    | identifier TokenDot identifier
    {
        $$ = DerivedType{ Base : $3, Class : $1 }
    }
    ;


simpleType
    : TokenUint { $$ = Uint{} }
    | TokenInt  { $$ = Int{} }
    | TokenText { $$ = Text{} }
    | TokenBool { $$ = Bool{} }
    | blob      { $$ = $1 }
    | dottedIdentifier { $$ = $1 }
    ; 

type:
    simpleType
    | list
    ;

future
    : TokenFuture TokenLParen simpleType TokenRParen
    {
        $$ = Future{ Type: $3 }
    }
    ;

typeOrFuture
    : type
    | future
    ; 

typedef:
    decorators TokenTypedef identifier uniqueIDOpt TokenEquals typeOrFuture TokenSemicolon
    {
        $$ = Typedef{
            BaseTypedef : BaseTypedef{
                BaseStatement: BaseStatement{ Dec : $1 }, 
                Ident : $3, 
                UniqueID : $4,
            },
            Type : $6,
        }
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

identifier:
    TokenIdentifier
    {
        $$ = Identifier{ Name : $1 }
    }
    ;

fileID:
    uniqueID TokenSemicolon identifier
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