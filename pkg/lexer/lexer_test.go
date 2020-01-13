package lexer

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/TykTechnologies/go-json-parser/pkg/ast"
	"github.com/TykTechnologies/go-json-parser/pkg/lexer/keyword"
)

func TestLexer_Read(t *testing.T) {
	type checkFunc func(lex *Lexer, i int)

	run := func(inStr string, checks ...checkFunc) {
		in := &ast.Input{}
		in.ResetInputBytes([]byte(inStr))
		lexer := &Lexer{}
		lexer.SetInput(in)

		for i := range checks {
			checks[i](lexer, i+1)
		}
	}

	mustRead := func(k keyword.Keyword, wantLiteral string) checkFunc {
		return func(lex *Lexer, i int) {
			tok := lex.Read()
			if k != tok.Keyword {
				panic(fmt.Errorf("mustRead: want(keyword): %s, got: %s [check: %d]", k.String(), tok.String(), i))
			}
			gotLiteral := string(lex.input.ByteSlice(tok.Literal))
			if wantLiteral != gotLiteral {
				panic(fmt.Errorf("mustRead: want(literal): %s, got: %s [check: %d]", wantLiteral, gotLiteral, i))
			}
		}
	}

	resetInput := func(input string) checkFunc {
		return func(lex *Lexer, i int) {
			lex.input.ResetInputBytes([]byte(input))
		}
	}
	_ = resetInput

	mustPeekWhitespaceLength := func(want int) checkFunc {
		return func(lex *Lexer, i int) {
			got := lex.peekWhitespaceLength()
			if want != got {
				panic(fmt.Errorf("mustPeekWhitespaceLength: want: %d, got: %d [check: %d]", want, got, i))
			}
		}
	}

	t.Run("peek whitespace length", func(t *testing.T) {
		run("   {", mustPeekWhitespaceLength(3))
	})
	t.Run("peek whitespace length with tab", func(t *testing.T) {
		run("   \t	{", mustPeekWhitespaceLength(5))
	})
	t.Run("peek whitespace length with line feed", func(t *testing.T) {
		run("   \n{", mustPeekWhitespaceLength(4))
	})
	t.Run("peek whitespace length with carriage return", func(t *testing.T) {
		run("   \r{", mustPeekWhitespaceLength(4))
	})
	t.Run("peek whitespace length with carriage return line feed", func(t *testing.T) {
		run("   \r\n{", mustPeekWhitespaceLength(5))
	})
	t.Run("peek whitespace length with comma", func(t *testing.T) {
		run("   {", mustPeekWhitespaceLength(3))
	})

	// Numbers
	t.Run("read number", func(t *testing.T) {
		run("123", mustRead(keyword.Number, "123"))
	})
	t.Run("read negative number", func(t *testing.T) {
		run("-123", mustRead(keyword.Minus, "-"),
			mustRead(keyword.Number, "123"))
	})
	t.Run("read number with fraction", func(t *testing.T) {
		run("123.45", mustRead(keyword.Number, "123.45"))
	})
	t.Run("read negative number with fraction", func(t *testing.T) {
		run("-123.45", mustRead(keyword.Minus, "-"),
			mustRead(keyword.Number, "123.45"))
	})
	t.Run("read plancks constant", func(t *testing.T) {
		run("6.63E-34", mustRead(keyword.Number, "6.63E-34"))
	})
	t.Run("read electron mass kg", func(t *testing.T) {
		run("9.10938356e-3", mustRead(keyword.Number, "9.10938356e-3"))
	})
	t.Run("read earth mass kg", func(t *testing.T) {
		run("5.9724e24", mustRead(keyword.Number, "5.9724e24"))
	})
	t.Run("read earth circumference m", func(t *testing.T) {
		run("4E7", mustRead(keyword.Number, "4E7"))
	})
	t.Run("read an inch in mm", func(t *testing.T) {
		run("2.54E+1", mustRead(keyword.Number, "2.54E+1"))
	})
	t.Run("read electron charge/mass ratio", func(t *testing.T) {
		run("-1.758E11", mustRead(keyword.Minus, "-"),
			mustRead(keyword.Number, "1.758E11"))
	})

	// Strings
	t.Run("read string", func(t *testing.T) {
		run("\"foo\"", mustRead(keyword.String, "foo"))
	})
	t.Run("read single line string with leading/trailing whitespace", func(t *testing.T) {
		run("\" 	foo	 \"", mustRead(keyword.String, " 	foo	 "))
	})

	// Structural Characters
	t.Run("read begin array", func(t *testing.T) {
		run("[", mustRead(keyword.BeginArray, "["))
	})
	t.Run("read end array", func(t *testing.T) {
		run("]", mustRead(keyword.EndArray, "]"))
	})
	t.Run("read begin object", func(t *testing.T) {
		run("{", mustRead(keyword.BeginObject, "{"))
	})
	t.Run("read end object", func(t *testing.T) {
		run("}", mustRead(keyword.EndObject, "}"))
	})
	t.Run("read true", func(t *testing.T) {
		run("true", mustRead(keyword.LiteralTrue, "true"))
	})
	t.Run("read true with space", func(t *testing.T) {
		run(" true ", mustRead(keyword.LiteralTrue, "true"))
	})
	t.Run("read false", func(t *testing.T) {
		run("false", mustRead(keyword.LiteralFalse, "false"))
	})
	t.Run("read null", func(t *testing.T) {
		run("null", mustRead(keyword.LiteralNull, "null"))
	})
}

var jsonDocument = `[
  {
    "id": "0001",
    "type": "donut",
    "name": "Cake",
    "ppu": 0.55,
    "batters": {
      "batter": [
        {
          "id": "1001",
          "type": "Regular"
        },
        {
          "id": "1002",
          "type": "Chocolate"
        },
        {
          "id": "1003",
          "type": "Blueberry"
        },
        {
          "id": "1004",
          "type": "Devil's Food"
        }
      ]
    },
    "topping": [
      {
        "id": "5001",
        "type": "None"
      },
      {
        "id": "5002",
        "type": "Glazed"
      },
      {
        "id": "5005",
        "type": "Sugar"
      },
      {
        "id": "5007",
        "type": "Powdered Sugar"
      },
      {
        "id": "5006",
        "type": "Chocolate with Sprinkles"
      },
      {
        "id": "5003",
        "type": "Chocolate"
      },
      {
        "id": "5004",
        "type": "Maple"
      }
    ]
  },
  {
    "id": "0002",
    "type": "donut",
    "name": "Raised",
    "ppu": 0.55,
    "batters": {
      "batter": [
        {
          "id": "1001",
          "type": "Regular"
        }
      ]
    },
    "topping": [
      {
        "id": "5001",
        "type": "None"
      },
      {
        "id": "5002",
        "type": "Glazed"
      },
      {
        "id": "5005",
        "type": "Sugar"
      },
      {
        "id": "5003",
        "type": "Chocolate"
      },
      {
        "id": "5004",
        "type": "Maple"
      }
    ]
  },
  {
    "id": "0003",
    "type": "donut",
    "name": "Old Fashioned",
    "ppu": 0.55,
    "batters": {
      "batter": [
        {
          "id": "1001",
          "type": "Regular"
        },
        {
          "id": "1002",
          "type": "Chocolate"
        }
      ]
    },
    "topping": [
      {
        "id": "5001",
        "type": "None"
      },
      {
        "id": "5002",
        "type": "Glazed"
      },
      {
        "id": "5003",
        "type": "Chocolate"
      },
      {
        "id": "5004",
        "type": "Maple"
      }
    ]
  }
]`

func BenchmarkLexer(b *testing.B) {
	in := &ast.Input{}
	lexer := &Lexer{}
	lexer.SetInput(in)

	inputBytes := []byte(jsonDocument)

	b.ReportAllocs()
	b.ResetTimer()
	b.SetBytes(int64(len(inputBytes)))

	for i := 0; i < b.N; i++ {

		in.ResetInputBytes(inputBytes)

		var key keyword.Keyword

		for key != keyword.EOF {
			key = lexer.Read().Keyword
		}
	}
}

func TestStdlib(t *testing.T) {
	var i interface{}
	err := json.Unmarshal([]byte("    false"), &i)
	if err != nil {
		t.Log(err.Error())
	}
}
