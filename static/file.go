package static

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/base32"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/errortypes"
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

func NewFile(path string) (file *File, err error) {
	ext := filepath.Ext(path)
	if len(ext) == 0 {
		return
	}

	typ, ok := mimeTypes[ext]
	if !ok {
		return
	}

	data, e := ioutil.ReadFile(path)
	if e != nil {
		err = &errortypes.ReadError{
			errors.Wrap(e, "static: Read error"),
		}
		return
	}

	hash := md5.Sum(data)
	hashStr := base32.StdEncoding.EncodeToString(hash[:])
	hashStr = strings.Replace(hashStr, "=", "", -1)
	hashStr = strings.ToLower(hashStr)

	file = &File{
		Type: typ,
		Hash: hashStr,
		Data: data,
	}

	gzipData := &bytes.Buffer{}

	writer, err := gzip.NewWriterLevel(gzipData, gzip.BestCompression)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "static: Gzip error"),
		}
		return
	}

	writer.Write(file.Data)
	writer.Close()
	file.GzipData = gzipData.Bytes()

	return
}
