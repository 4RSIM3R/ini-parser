package main

import (
	"os"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/alecthomas/repr"
)

type Section struct {
	Identifier string      `"[" @Ident "]"`
	Property   []*Property `@@*`
}

type Value interface{ value() }

type String struct {
	String string `@String`
}

func (String) value() {}

type Number struct {
	Number float64 `@Float`
}

func (Number) value() {}

type Boolean struct {
	Bool bool `@Boolean`
}

func (Boolean) value() {}

type Property struct {
	Key   string `@Ident "="`
	Value Value  `@@`
}

type INI struct {
	Properties []*Property `@@*`
	Sections   []*Section  `@@*`
}

func main() {
	var lexer = lexer.MustSimple([]lexer.SimpleRule{
		{`Boolean`, `(?i:true|false)`}, // Prioritize matching booleans
		{`Ident`, `[a-zA-Z][a-zA-Z_\d]*`},
		{`String`, `"(?:\\.|[^"\\])*"|'(?:\\.|[^'\\])*'`},
		{`Float`, `\d+(?:\.\d+)?`},
		{`Punct`, `[][=]`},
		{"comment", `[#;][^\n]*`},
		{"whitespace", `\s+`},
	})

	var parser = participle.MustBuild[INI](
		participle.Lexer(lexer),
		participle.Unquote("String"),
		participle.Union[Value](String{}, Number{}, Boolean{}),
	)

	content, err := os.ReadFile("example.ini")

	if err != nil {
		panic(err)
	}

	ini, err := parser.ParseString("", string(content))

	if err != nil {
		panic(err)
	}

	repr.Println(ini, repr.Indent("  "), repr.OmitEmpty(true))
}
