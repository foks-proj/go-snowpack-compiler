package lib

type UniqueID struct {
	Val string // is a number, but we don't bother to parse it
}

func (u UniqueID) IsSet() bool {
	return u.Val != ""
}

type FileNode struct {
	Id    UniqueID
	Stmts []Statement
}

type Statement interface {
}

type BaseStatement struct {
	Dec Decorators
}

type BaseTypedef struct {
	BaseStatement
	Ident    Identifier
	UniqueID UniqueID
}

type Importer interface {
	Statement
}

type BaseImport struct {
	Path string
	Name string
}

type GenericImport struct {
	BaseImport
}

type TypeScriptImport struct {
	BaseImport
}

type GoImport struct {
	BaseImport
}

type Typedef struct {
	BaseTypedef
	Type Type
}

type Decorators struct {
	Doc Docstring
}

type Docstring struct {
	Raw string
}

type Identifier struct {
	Name string
}

type Type interface {
}

type Future struct {
	Type Type
}

type List struct {
	Type Type
}

type Blob struct {
	Count int
}

type DerivedType struct {
	Base  Identifier
	Class Identifier
}

type Struct struct {
	BaseTypedef
}

type Text struct{}
type Uint struct{}
type Int struct{}
type Bool struct{}

var _ Statement = Typedef{}
var _ Type = List{}
var _ Type = Future{}
var _ Type = Text{}
var _ Type = Uint{}
var _ Type = Int{}
var _ Type = Bool{}
var _ Importer = GenericImport{}
var _ Importer = TypeScriptImport{}
