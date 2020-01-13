// Package token contains the object and logic needed to describe a lexed token in a JSON document
package token

import (
	"fmt"

	"github.com/TykTechnologies/go-json-parser/pkg/ast"
	"github.com/TykTechnologies/go-json-parser/pkg/lexer/keyword"
)

type Token struct {
	Keyword keyword.Keyword
	Literal ast.ByteSliceReference
}

func (t Token) String() string {
	return fmt.Sprintf("token:: Keyword: %s", t.Keyword)
}

func (t *Token) SetStart(inputPosition int) {
	t.Literal.Start = uint32(inputPosition)
}

func (t *Token) SetEnd(inputPosition int) {
	t.Literal.End = uint32(inputPosition)
}
