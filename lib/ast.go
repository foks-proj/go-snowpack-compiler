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

type Field struct {
	Ident Identifier
	Pos   int
	Type  Type
}

type Struct struct {
	BaseTypedef
	Fields []Field
}

type CaseLabel interface {
}

type CaseLabelIdentifier struct {
	Ident Identifier
}

type CaseLabelNumber struct {
	Num int
}

type CaseLabelBool struct {
	Bool bool
}

type Case struct {
	Labels   []CaseLabel // nil for default case
	Position *int        // will be nil for void data; 0 is a valid position
	Type     Type
}

type Variant struct {
	BaseTypedef
	SwitchVar  Identifier
	SwitchType Type
	Cases      []Case
}

type Option struct {
	Type Type
}

type Void struct {
}

type Text struct{}
type Uint struct{}
type Int struct{}
type Bool struct{}

var _ CaseLabel = CaseLabelIdentifier{}
var _ CaseLabel = CaseLabelNumber{}
var _ CaseLabel = CaseLabelBool{}
var _ Statement = Typedef{}
var _ Type = List{}
var _ Type = Future{}
var _ Type = Text{}
var _ Type = Uint{}
var _ Type = Int{}
var _ Type = Bool{}
var _ Type = Option{}
var _ Type = Void{}
var _ Importer = GenericImport{}
var _ Importer = TypeScriptImport{}
