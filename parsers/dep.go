package parsers

import (
	"errors"

	"github.com/mdevilliers/dtree"
)

type dep struct{}

func (d dep) Parse(path string) ([]dtree.Node, []dtree.Edge, error) {
	return nil, nil, errors.New("not implemented")
}
