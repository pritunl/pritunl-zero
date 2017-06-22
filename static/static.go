// Versions static files with hash, replaces references and stores in memory.
package static

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/base32"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/errortypes"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
)

var (
	mimeTypes = map[string]string{
		".js":    "application/javascript",
		".json":  "application/json",
		".css":   "text/css",
		".html":  "text/html",
		".jpg":   "image/jpeg",
		".png":   "image/png",
		".svg":   "image/svg+xml",
		".ico":   "image/vnd.microsoft.icon",
		".otf":   "application/font-sfnt",
		".ttf":   "application/font-sfnt",
		".woff":  "application/font-woff",
		".woff2": "font/woff2",
		".ijmap": "text/plain",
		".eot":   "application/vnd.ms-fontobject",
		".map":   "application/json",
	}
)

type File struct {
	Type     string
	Hash     string
	Data     []byte
	GzipData []byte
}

type Store struct {
	Files map[string]*File
	root  string
}

func (s *Store) addDir(dir string) (err error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}

	for _, info := range files {
		name := info.Name()
		fullPath := path.Join(dir, name)

		if info.IsDir() {
			s.addDir(fullPath)
			continue
		}

		ext := filepath.Ext(name)
		if len(ext) == 0 {
			continue
		}

		typ, ok := mimeTypes[ext]
		if !ok {
			continue
		}

		data, e := ioutil.ReadFile(fullPath)
		if e != nil {
			err = e
			return
		}

		hash := md5.Sum(data)
		hashStr := base32.StdEncoding.EncodeToString(hash[:])
		hashStr = strings.Replace(hashStr, "=", "", -1)
		hashStr = strings.ToLower(hashStr)

		file := &File{
			Type: typ,
			Hash: hashStr,
			Data: data,
		}

		s.Files[fullPath] = file
	}

	return
}

func (s *Store) parseFiles() {
	for _, file := range s.Files {
		data := &bytes.Buffer{}
		writer := gzip.NewWriter(data)
		writer.Write(file.Data)
		writer.Close()
		file.GzipData = data.Bytes()
	}
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

	store.parseFiles()

	return
}

func GetMimeType(name string) string {
	return mimeTypes[path.Ext(name)]
}
