package pegman

import (
	"log"
	"sort"
	"unicode"
	"unicode/utf8"

	"golang.org/x/exp/ebnf"
)

func Parse(files ...string) (Grammar, error) {
	//grammar, err := ebnf.Parse(path, bufio.NewReader(file))
	//if err != nil {
	//	return fmt.Errorf("error parsing %q: %w", path, err)
	//}
	return nil, nil
}

// First returns the first token set of a given expression
func First(g ebnf.Grammar, x ebnf.Expression) []string {
	return first(g, x, make(map[ebnf.Expression]bool))
}

func first(g ebnf.Grammar, x ebnf.Expression, v map[ebnf.Expression]bool) []string {
	switch x := x.(type) {
	case *ebnf.Token:
		return []string{x.String}
	case ebnf.Sequence:
		var ret []string
		if len(x) > 0 {
			ret = append(ret, first(g, x[0], v)...)
		}
		if _, ok := x[0].(*ebnf.Option); ok && len(x) > 1 {
			ret = append(ret, first(g, x[1], v)...)
		}
		return ret
	case *ebnf.Name:
		if !v[x] {
			v[x] = true
			return first(g, g[x.String], v)
		}
		return nil

	case *ebnf.Option:
		return first(g, x.Body, v)
	case *ebnf.Repetition:
		return first(g, x.Body, v)
	case ebnf.Alternative:
		var ret []string
		for _, alt := range x {
			ret = append(ret, first(g, alt, v)...)
		}
		return ret
	case *ebnf.Group:
		return first(g, x.Body, v)
	case *ebnf.Production:
		if x.Expr == nil {
			return nil
		}
		return first(g, x.Expr, v)
	default:
		log.Printf("first: unhandled expression type: %T\n", x)
		return nil
	}
}

// Productions returns the productions of the grammar in the order they appear
// in the source file.
func Productions(g ebnf.Grammar) []*ebnf.Production {
	ret := make([]*ebnf.Production, 0, len(g))
	for _, prod := range g {
		ret = append(ret, prod)
	}
	sort.SliceStable(ret, func(i, j int) bool {
		return ret[i].Pos().Offset < ret[j].Pos().Offset
	})
	return ret
}

// Tokens returns an alphabetically sorted list of unique tokens that are used
// in the grammar.
func Tokens(g ebnf.Grammar) []string {
	m := make(map[string]int)
	for _, prod := range g {
		Inspect(prod.Expr, func(e ebnf.Expression) bool {
			if tok, ok := e.(*ebnf.Token); ok {
				m[tok.String]++
			}
			return true
		})
	}
	ret := make([]string, 0, len(m))
	for k := range m {
		ret = append(ret, k)
	}
	sort.Strings(ret)
	return ret
}

// Inspect traverses the given expression and calls the given function for each.
func Inspect(e ebnf.Expression, fn func(e ebnf.Expression) bool) bool {
	if e == nil {
		return true
	}

	if !fn(e) {
		return false
	}

	switch e := e.(type) {
	case ebnf.Alternative:
		for _, alt := range e {
			if !Inspect(alt, fn) {
				return false
			}
		}
	case ebnf.Sequence:
		for _, seq := range e {
			if !Inspect(seq, fn) {
				return false
			}
		}

	case *ebnf.Group:
		if e.Body != nil {
			Inspect(e.Body, fn)
		}

	case *ebnf.Option:
		if e.Body != nil {
			Inspect(e.Body, fn)
		}

	case *ebnf.Repetition:
		if e.Body != nil {
			Inspect(e.Body, fn)
		}

	}
	return true
}

// IsLexical returns true, when given name is a lexical production.
func IsLexical(name string) bool {
	ch, _ := utf8.DecodeRuneInString(name)
	return !unicode.IsUpper(ch)
}
