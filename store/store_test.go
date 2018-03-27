package store

import (
	"testing"

	"github.com/mdevilliers/dtree"
	"github.com/stretchr/testify/assert"
)

func Test_FindNodes(t *testing.T) {

	// A
	// | \
	// BB CB

	a := dtree.RootNode("A")
	b := dtree.NewNode("BB", "1.0")
	c := dtree.NewNode("CB", "2.0")

	nodes := []dtree.Node{a, b, c}

	s := InMemory(nodes, []dtree.Edge{})

	n, err := s.FindNodes("A")

	assert.Nil(t, err)
	assert.Len(t, n, 1)

	n, err = s.FindNodes("AB")

	assert.Nil(t, err)
	assert.Len(t, n, 0)

	n, err = s.FindNodes("B")

	assert.Nil(t, err)
	assert.Len(t, n, 2)

	n, err = s.FindNodes("does-not-exist")

	assert.Nil(t, err)
	assert.Len(t, n, 0)

}
