// Package lexer contains the logic to turn an ast.Input into lexed tokens
package lexer

import (
	"github.com/TykTechnologies/go-json-parser/pkg/ast"
	"github.com/TykTechnologies/go-json-parser/pkg/lexer/keyword"
	"github.com/TykTechnologies/go-json-parser/pkg/lexer/runes"
	"github.com/TykTechnologies/go-json-parser/pkg/lexer/token"
)

// Lexer emits tokens from an input reader
type Lexer struct {
	input *ast.Input
}

func (l *Lexer) SetInput(input *ast.Input) {
	l.input = input
}

// Read emits the next token
func (l *Lexer) Read() (tok token.Token) {
	var next byte

	// skips insignificant whitespace
	for {
		tok.SetStart(l.input.InputPosition)
		next = l.readRune()
		if !l.byteIsWhitespace(next) {
			break
		}
	}

	if l.matchSingleRuneToken(next, &tok) {
		return
	}

	switch next {
	case runes.Quote:
		l.readString(&tok)
		return
	}

	if runeIsDigit(next) {
		l.readDigit(&tok)
		return
	}

	tok.SetEnd(l.input.InputPosition)
	return
}

func (l *Lexer) readString(tok *token.Token) {
	tok.Keyword = keyword.Quote

	tok.SetStart(l.input.InputPosition)

	escaped := false
	whitespaceCount := 0
	reachedFirstNonWhitespace := false
	leadingWhitespaceToken := 0

	for {
		next := l.readRune()
		switch next {
		case runes.Space, runes.Tab:
			whitespaceCount++
		case runes.EOF:
			tok.SetEnd(l.input.InputPosition)
			tok.Literal.Start += uint32(leadingWhitespaceToken)
			tok.Literal.End -= uint32(whitespaceCount)
			return
		case runes.Quote, runes.CarriageReturn, runes.LineFeed:
			if escaped {
				escaped = !escaped
				continue
			}

			tok.SetEnd(l.input.InputPosition - 1)
			tok.Literal.Start += uint32(leadingWhitespaceToken)
			tok.Literal.End -= uint32(whitespaceCount)
			return
		case runes.Backslash:
			escaped = !escaped
			whitespaceCount = 0
		default:
			if !reachedFirstNonWhitespace {
				reachedFirstNonWhitespace = true
				leadingWhitespaceToken = whitespaceCount
			}
			escaped = false
			whitespaceCount = 0
		}
	}
}

func (l *Lexer) swallowAmount(amount int) {
	for i := 0; i < amount; i++ {
		l.readRune()
	}
}

func (l *Lexer) peekWhitespaceLength() (amount int) {
	for i := l.input.InputPosition; i < l.input.Length; i++ {
		if !l.byteIsWhitespace(l.input.RawBytes[i]) {
			break
		}
		amount++
	}

	return amount
}

func (l *Lexer) peekEquals(ignoreWhitespace bool, equals ...byte) bool {
	var whitespaceOffset int
	if ignoreWhitespace {
		whitespaceOffset = l.peekWhitespaceLength()
	}

	start := l.input.InputPosition + whitespaceOffset
	end := l.input.InputPosition + len(equals) + whitespaceOffset

	if end > l.input.Length {
		return false
	}

	for i := 0; i < len(equals); i++ {
		if l.input.RawBytes[start+i] != equals[i] {
			return false
		}
	}

	return true
}

func runeIsDigit(r byte) bool {
	switch {
	case r >= '0' && r <= '9':
		return true
	default:
		return false
	}
}

func (l *Lexer) peekRune(ignoreWhitespace bool) (r byte) {
	for i := l.input.InputPosition; i < l.input.Length; i++ {
		r = l.input.RawBytes[i]
		if !ignoreWhitespace {
			return r
		} else if !l.byteIsWhitespace(r) {
			return r
		}
	}

	return runes.EOF
}

func (l *Lexer) readDigit(tok *token.Token) {
	var r byte
	for {
		r = l.peekRune(false)
		if !runeIsDigit(r) {
			break
		}
		l.readRune()
	}

	containsExponent := r == runes.ExponentL || r == runes.ExponentU

	containsFraction := r == runes.Dot || containsExponent
	if containsFraction {
		l.readRune()
		l.readFraction(containsExponent, tok)
		return
	}

	tok.Keyword = keyword.Number
	tok.SetEnd(l.input.InputPosition)
}

func (l *Lexer) readFraction(hasReadExponent bool, tok *token.Token) {
	var r byte
	for {
		r = l.peekRune(false)
		if !runeIsDigit(r) {
			break
		}
		l.readRune()
	}

	if hasReadExponent {
		fraction := keyword.Number
		tok.Keyword = fraction
		tok.SetEnd(l.input.InputPosition)
		return
	}

	optionalExponent := l.peekRune(false)
	if optionalExponent == runes.ExponentL || optionalExponent == runes.ExponentU {
		l.readRune()
	}

	optionalPlusMinus := l.peekRune(false)
	if optionalPlusMinus == runes.Plus || optionalPlusMinus == runes.Minus {
		l.readRune()
	}

	for {
		r = l.peekRune(false)
		if !runeIsDigit(r) {
			break
		}
		l.readRune()
	}

	tok.Keyword = keyword.Number
	tok.SetEnd(l.input.InputPosition)
}

func (l *Lexer) readRune() (r byte) {
	if l.input.InputPosition < l.input.Length {
		r = l.input.RawBytes[l.input.InputPosition]

		l.input.InputPosition++
	} else {
		r = runes.EOF
	}

	return
}

func (l *Lexer) byteIsWhitespace(r byte) bool {
	switch r {
	case runes.Space, runes.Tab, runes.CarriageReturn, runes.LineFeed:
		return true
	default:
		return false
	}
}

func (l *Lexer) matchSingleRuneToken(r byte, tok *token.Token) bool {
	switch r {
	case runes.EOF:
		tok.Keyword = keyword.EOF
	case runes.Colon:
		tok.Keyword = keyword.NameSeparator
	case runes.LeftBrace:
		tok.Keyword = keyword.BeginObject
	case runes.LeftBracket:
		tok.Keyword = keyword.BeginArray
	case runes.RightBrace:
		tok.Keyword = keyword.EndObject
	case runes.RightBracket:
		tok.Keyword = keyword.EndArray
	case runes.Minus:
		tok.Keyword = keyword.Minus
	default:
		return false
	}

	tok.SetEnd(l.input.InputPosition)

	return true
}
