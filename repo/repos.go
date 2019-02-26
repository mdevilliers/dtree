package repo

import (
	"io/ioutil"
	"path"
)

// Repository encapsulates a golang source code repository.
type Repository struct {
	Path string
}

// FromCheckedOut returns a multiple Repository or an error
func FromCheckedOut(root string) ([]Repository, error) {

	files, err := ioutil.ReadDir(root)

	if err != nil {
		return nil, err
	}

	paths := []Repository{}

	for _, file := range files {

		if file.IsDir() {
			paths = append(paths, Repository{Path: path.Join(root, file.Name())})
		}
	}

	return paths, nil
}
