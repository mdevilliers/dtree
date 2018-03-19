package store

import (
	"fmt"
	"testing"

	"github.com/mdevilliers/dtree"
)

func Test_Group(t *testing.T) {

	// A
	// | \
	// B  C

	a := dtree.RootNode("A")
	b := dtree.NewNode("B", "1.0")
	c := dtree.NewNode("C", "2.0")

	nodes := []dtree.Node{a, b, c}

	e1 := dtree.NewDependancy(a, b, "1.0")
	e2 := dtree.NewDependancy(a, c, "2.0")

	edges := append(e1, e2...)

	s := InMemory(nodes, edges)

	fmt.Println(s)
	fmt.Println(s.Group())
}
