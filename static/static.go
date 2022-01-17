// Versions static files with hash, replaces references and stores in memory.
package static

import (
	"io/ioutil"
	"path"

	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/errortypes"
)

type Store struct {
	Files map[string]*File
	root  string
}

func (s *Store) addDir(dir string) (err error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		err = &errortypes.ReadError{
			errors.Wrap(err, "static: Read directory error"),
		}
		return
	}

	for _, info := range files {
		name := info.Name()
		fullPath := path.Join(dir, name)

		if info.IsDir() {
			err = s.addDir(fullPath)
			if err != nil {
				return
			}

			continue
		}

		file, e := NewFile(fullPath)
		if e != nil {
			err = e
			return
		}

		if file != nil {
			s.Files[fullPath] = file
		}
	}

	return
}

func NewStore(root string) (store *Store, err error) {
	store = &Store{
		Files: map[string]*File{},
		root:  root,
	}

	err = store.addDir(root)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "static: Init error"),
		}
		return
	}

	return
}

func GetMimeType(name string) string {
	return mimeTypes[path.Ext(name)]
}
