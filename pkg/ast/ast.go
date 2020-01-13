package ast

type Document struct {
	Input Input
}

func NewDocument() *Document {
	return &Document{}
}
