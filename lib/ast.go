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

type BaseImport interface {
	Statement
}

type GenericImport struct {
	Path string
	Name string
}

var _ BaseImport = &GenericImport{}
