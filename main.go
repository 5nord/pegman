package main

import (
	"bufio"
	"log"
	"os"
	"text/template"

	"golang.org/x/exp/ebnf"
)

var (
	funcs = template.FuncMap{
		"type": func(v interface{}) string {
			switch v.(type) {
			case *ebnf.Token:
				return "Token"
			case *ebnf.Group:
				return "Group"
			case *ebnf.Name:
				return "Name"
			case *ebnf.Alternative:
				return "Alternative"
			case *ebnf.Option:
				return "Option"
			case *ebnf.Range:
				return "Range"
			case *ebnf.Repetition:
				return "Repetition"
			case *ebnf.Sequence:
				return "Sequence"
			default:
				return "unknown"
			}
		},
	}
)

func main() {
	file, err := os.Open("ttcn3.ebnf")
	if err != nil {
		log.Fatal(err)
	}

	grammar, err := ebnf.Parse("ttcn3.ebnf", bufio.NewReader(file))
	if err != nil {
		log.Fatal(err)
	}

	if err := ebnf.Verify(grammar, "Module"); err != nil {
		log.Fatal(err)
	}

	parser, err := template.New("parser.tmpl").Funcs(funcs).ParseFiles("parser.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	os.Mkdir("ttcn3", 0755)
	w, err := os.Create("ttcn3/parser_gen.go")
	if err != nil {
		log.Fatal(err)
	}
	if err := parser.Execute(w, grammar); err != nil {
		log.Fatal(err)
	}
}
