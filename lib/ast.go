package lib

type UniqueID struct {
	Val string // is a number, but we don't bother to parse it
}

type FileNode struct {
	Id    UniqueID
	Stmts []Statement
}

type Statement interface {
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
}

type Decorators struct {
	Doc Docstring
}

type Docstring struct {
	Raw string
}

var _ Statement = &Typedef{}

var _ Importer = &GenericImport{}
var _ Importer = &TypeScriptImport{}
