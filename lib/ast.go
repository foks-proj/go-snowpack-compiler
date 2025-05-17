package lib

type UniqueID struct {
	Val string // is a number, but we don't bother to parse it
}

type FileNode struct {
	Id UniqueID
}

type Statement interface {
}

type BaseImport interface {
	Statement
}
