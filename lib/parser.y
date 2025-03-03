
%{
package lib

import (
    "fmt"
    "strconv"
)

%}

%union {
    root     Root
    uniqueId UniqueID
    uint     uint64
    rawval   string
    stmts    []Statement
    stmt     Statement
    imprt    Import
    dec      Decorators
    doc      Docstring
    docRaw   string
    ident    Identifier
    typ      Type
    num      int
    field    Field
    fields   []Field
    cases    []Case
    cas      Case
    intp     *int
    caseLabel  CaseLabel
    caseLabels []CaseLabel
    enumValues []EnumValue
    enumValue  EnumValue
    protoModifiers []ProtocolModifier
    protoModifier  ProtocolModifier
    params []Param
    param Param
    method Method
    methods []Method

}

%type <root> top
%type <uniqueId> fileID uniqueID uniqueIDOpt
%type <stmts> statements
%type <stmt> statement typedef struct variant enum protocol
%type <imprt> import genericImport tsImport goImport
%type <dec> decorators
%type <doc> doc 
%type <docRaw> docRaw
%type <ident> identifier argTypeOpt
%type <typ> list type simpleType typeOrFuture blob dottedIdentifier future typeOrOptional optionalType typeOrVoid returnOpt
%type <num> countOpt number position
%type <field> field
%type <fields> fields
%type <cases> cases
%type <cas> case normalCase defaultCase
%type <caseLabel> caseLabel
%type <caseLabels> caseLabels
%type <enumValues> enumValues
%type <enumValue> enumValue
%type <protoModifiers> protoModifiers
%type <protoModifier> protoModifier
%type <params> paramList paramsOpt params
%type <param> param
%type <intp> positionOpt
%type <rawval> uintConstant
%type <method> method
%type <methods> methods

%token TokenAt TokenSemicolon TokenAs TokenEquals TokenDot
%token TokenImport TokenTypeScriptImport TokenGoImport
%token TokenList TokenLParen TokenRParen TokenText TokenUint TokenInt TokenBool TokenBlob TokenFuture
%token TokenLBrace TokenRBrace TokenStruct TokenOption TokenColon TokenVariant TokenSwitch TokenCase
%token TokenTrue TokenFalse TokenDefault TokenVoid TokenEnum 
%token TokenProtocol TokenErrors TokenArgHeader TokenResHeader
%token TokenArrow TokenComma

%token <rawval> TokenUint64Val TokenIntVal TokenUint32Val 
%token <rawval> TokenDQoutedString TokenIdentifier TokenDoc TokenTypedef 

%%

top: 
    fileID statements
    {
        $$ = Root{ Id: $1, Stmts : $2 }
        top = &$$ // Set the global top variable
    }
    ;

statements:
    /* empty */ { $$ = []Statement{} }
    | statements statement { $$ = append($1, $2) }
    ;

genericImport: 
    TokenImport TokenDQoutedString TokenAs TokenIdentifier TokenSemicolon
    { 
        $$ = Import{ Path: $2, Name : $4, Lang : LangGeneric }
    }
    ;

tsImport: 
    TokenTypeScriptImport TokenDQoutedString TokenAs TokenIdentifier TokenSemicolon
    { 
        $$ = Import { Path: $2, Name : $4, Lang : LangTypeScript }
    } 
    ;

goImport: 
    TokenGoImport TokenDQoutedString TokenAs TokenIdentifier TokenSemicolon
    { 
        $$ = Import { Path: $2, Name : $4, Lang : LangGo }
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
    | docRaw TokenDoc { $$ = $1 + $2; }
    ;

decorators:
    doc { $$ = Decorators{ Doc: $1 } }
    ;

uniqueIDOpt
    : /* empty */ { $$ = UniqueID{} }
    | uniqueID { $$ = $1; }
    ;

list:
    TokenList TokenLParen type TokenRParen
    {
        $$ = List{ Type: $3 }
    }
    ;

number: TokenIntVal
    {
        var i int
        i, err := strconv.Atoi($1)
        if err != nil {
            parseErr = err
        } else if i < 0 {
            parseErr = fmt.Errorf("blob byte-count must be greater than 0")
        } else {
            $$ = i
        }
    }
    ;

countOpt:
    /* empty */ { $$ = 0 }
    | TokenLParen number TokenRParen
    {
        if $2 <= 0 {
            parseErr = fmt.Errorf("blob byte-count must be greater than 0")
        } else {
            $$ = $2
        }
    }
    ;

blob
    : TokenBlob countOpt { $$ = Blob{ Count: $2 } }
    ;

dottedIdentifier:
    identifier
    {
        $$ = DerivedType{ Name : $1 }
    }
    | identifier TokenDot identifier
    {
        $$ = DerivedType{ ImportedFrom : $1, Name : $3 }
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
    : TokenFuture TokenLParen simpleType TokenRParen { $$ = Future{ Type: $3 } }
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

position: TokenAt number { $$ = $2 };

positionOpt 
    : /* empty */ { $$ = nil }
    | position
    {
        tmp := $1
        $$ = &tmp
    }
    ;

optionalType
    : TokenOption TokenLParen type TokenRParen
    {
        $$ = Option{ Type: $3 }
    }
    ;

typeOrOptional
    : type
    | optionalType
    ;

field:
    identifier position TokenColon typeOrOptional TokenSemicolon
    {
        $$ = Field{
            Ident : $1,
            Pos : $2,
            Type : $4,
        }
    }
    ;

fields 
    : /* empty */ { $$ = []Field{} }
    | fields field { $$ = append($1, $2) }
    ;

struct:
    decorators TokenStruct identifier uniqueIDOpt TokenLBrace fields TokenRBrace
    {
        $$ = Struct{
            BaseTypedef : BaseTypedef{
                BaseStatement: BaseStatement{ Dec : $1 }, 
                Ident : $3, 
                UniqueID : $4,
            },
            Fields : $6,
        }
    }
    ;

cases
    : case { $$ = []Case{ $1 } }
    | cases case { $$ = append($1, $2) }
    ;

case
    : normalCase  { $$ = $1 }
    | defaultCase { $$ = $1 }
    ;

caseLabels
    : caseLabel { $$ = []CaseLabel{ $1 } }
    | caseLabels TokenComma caseLabel { $$ = append($1, $3) }
    ;

caseLabel 
    : identifier { $$ = CaseLabelIdentifier{ Ident: $1 } }
    | number { $$ = CaseLabelNumber{ Num: $1 } }
    | TokenTrue { $$ = CaseLabelBool{ Bool: true } }
    | TokenFalse { $$ = CaseLabelBool{ Bool: false } }
    ;

typeOrVoid
    : type { $$ = $1 }
    | TokenVoid { $$ = Void{} }
    ;

normalCase
    : TokenCase caseLabels positionOpt TokenColon typeOrVoid TokenSemicolon
    {
        $$ = Case{
            Labels : $2,
            Position : $3,
            Type : $5,
        }
    }
    ;

defaultCase
    : TokenDefault positionOpt TokenColon typeOrVoid TokenSemicolon
    {
        $$ = Case{
            Labels : nil,
            Position : $2,
            Type : $4,
        }
    }
    ;

variant:
    decorators TokenVariant identifier TokenSwitch 
        TokenLParen identifier TokenColon simpleType TokenRParen uniqueIDOpt
        TokenLBrace cases TokenRBrace
    {
        $$ = Variant{
            BaseTypedef : BaseTypedef{
                BaseStatement: BaseStatement{ Dec : $1 }, 
                Ident : $3, 
                UniqueID : $10,
            },
            SwitchVar : $6,
            SwitchType : $8,
            Cases : $12,
        }
    }
    ;

enumValue
    : identifier TokenAt number TokenSemicolon
    {
        $$ = EnumValue{
            Ident : $1,
            Num : $3,
        }
    }
    ;

enumValues
    : enumValue { $$ = []EnumValue{ $1 } }
    | enumValues enumValue { $$ = append($1, $2) }
    ;

enum: 
    decorators TokenEnum identifier TokenLBrace enumValues TokenRBrace
    {
        $$ = Enum{
            BaseTypedef : BaseTypedef{
                BaseStatement: BaseStatement{ Dec : $1 }, 
                Ident : $3, 
            },
            Values : $5,
        }
    }
    ;

import: 
    genericImport { $$ = $1 }
    | tsImport    { $$ = $1 }
    | goImport    { $$ = $1 }
    ;

statement
    : import   { $$ = $1 }
    | typedef  { $$ = $1 }
    | struct   { $$ = $1 }
    | variant  { $$ = $1 }
    | enum     { $$ = $1 }
    | protocol { $$ = $1 }
    ;

identifier:
    TokenIdentifier { $$ = Identifier{ Name : $1 } }
    ;

fileID:
    uniqueID TokenSemicolon { $$ = $1 }
    ;

uintConstant
    : TokenUint64Val { $$ = $1 }
    | TokenUint32Val { $$ = $1 }
    ;

uniqueID:
    TokenAt uintConstant { $$ = UniqueID{ Val: $2 } }
    ;

protoModifier
    : TokenErrors type    { $$ = Errors{ Type: $2 } }
    | TokenArgHeader type { $$ = ArgHeader{ Type: $2 } }
    | TokenResHeader type { $$ = ResHeader{ Type: $2 } }
    ;

protoModifiers
    : /* empty */ { $$ = nil }
    | protoModifiers protoModifier { $$ = append($1, $2) }
    ;

argTypeOpt
    : /* empty */ { $$ = Identifier{} }
    | TokenColon identifier { $$ = $2 }
    ;

paramsOpt
    : /* empty */ { $$ = nil }
    | params
    ;

params
    : param { $$ = []Param{ $1 } }
    | params TokenComma param { $$ = append($1, $3) }
    ;

param
    : identifier position TokenColon typeOrOptional
    {
        $$ = Param{
            Ident : $1,
            Pos : $2,
            Type : $4,
        }
    }
    ;


paramList
    : TokenLParen paramsOpt TokenRParen { $$ = $2 }
    ;

returnOpt
    : /* empty */ { $$ = Void{} }
    | TokenArrow typeOrVoid { $$ = $2 }
    ;

method
    : decorators identifier position paramList argTypeOpt returnOpt TokenSemicolon
        {
            $$ = Method{
                BaseTypedef : BaseTypedef{
                    BaseStatement: BaseStatement{ Dec : $1 }, 
                    Ident : $2, 
                },
                Pos : $3,
                Params : $4,
                ArgType : $5,
                ResType : $6,
            }
        }
    ;

methods
    : /* empty */ { $$ = nil }
    | methods method { $$ = append($1, $2) }
    ;


protocol
    : decorators TokenProtocol identifier protoModifiers uniqueID 
        TokenLBrace methods TokenRBrace 
    {
        pmsp, err := NewProtocolModifiers($4)
        if err != nil {
            parseErr = err
        }
        var pms ProtocolModifiers
        if pmsp != nil {
            pms = *pmsp
        }
        $$ = Protocol{
            BaseTypedef : BaseTypedef{
                BaseStatement: BaseStatement{ Dec : $1 }, 
                Ident : $3, 
                UniqueID : $5,
            },
            Modifiers : pms,
            Methods : $7,
        }
    }
    ;

%%