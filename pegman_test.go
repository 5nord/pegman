package pegman_test

import (
	"strings"
	"testing"

	"github.com/5nord/pegman"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/ebnf"
)

func TestFirst(t *testing.T) {
	t.Run(t.Name(), func(t *testing.T) {
		g := parse(t, `A = "foo" "bar".`)
		actual := pegman.First(g, g["A"].Expr)
		assert.Equal(t, []string{"foo"}, actual)
	})
	t.Run(t.Name(), func(t *testing.T) {
		g := parse(t, `A = ["foo"] "bar".`)
		actual := pegman.First(g, g["A"].Expr)
		assert.Equal(t, []string{"foo", "bar"}, actual)
	})
	t.Run(t.Name(), func(t *testing.T) {
		g := parse(t, `A = "foo"|"bar".`)
		actual := pegman.First(g, g["A"].Expr)
		assert.Equal(t, []string{"foo", "bar"}, actual)
	})

}

func parse(t *testing.T, input string) ebnf.Grammar {
	g, err := ebnf.Parse("test", strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	return g
}
