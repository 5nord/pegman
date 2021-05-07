package ttcn3

func Parse(text string) *Tree {
	p := parser{
		text: text,
	}
	p.parse()
	return &Tree{}
}

// Tree represents a parsed input
type Tree struct {
}

// Ein Token hat nur eine Laenge und weiss was es ist.
// SEMICOLON  ==> (";",   1)
// IDENTIFIER ==> ("Foo", 3)

type Node interface {
	Pos() int
	Kind() Kind
	Value() string
	Children() []Node
}

type Kind int32

type kindNode struct {
	kind Kind
}

type valueNode struct {
	kindNode
	value string
}

type node struct {
	valueNode
	children []Node
}

type parser struct {
	text string
	pos  int
	buf  []Node
}

func (p *parser) next() {
}

func (p *parser) parse() {
	return
}
