package repo

import (
	"io/ioutil"
	"path"
)

type Repository struct {
	Path string
}

type from_disk struct {
	root string
}

func FromCheckedOut(root string) from_disk {
	return from_disk{root: root}
}

func (f from_disk) Paths() ([]Repository, error) {

	files, err := ioutil.ReadDir(f.root)

	if err != nil {
		return nil, err
	}

	paths := []Repository{}

	for _, file := range files {

		if file.IsDir() {
			paths = append(paths, Repository{Path: path.Join(f.root, file.Name())})
		}
	}

	return paths, nil
}
