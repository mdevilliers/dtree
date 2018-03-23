package parsers

import (
	"errors"
	"os"

	"github.com/mdevilliers/dtree"
)

type parser interface {
	Parse(pth string) ([]dtree.Node, []dtree.Edge, error)
}

func New(path string) (parser, error) {

	dep := Dep()
	if dep.Test(path) {
		return dep, nil
	}

	glide := Glide()

	if glide.Test(path) {
		return glide, nil
	}

	return nil, errors.New("parser not found")
}

func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
