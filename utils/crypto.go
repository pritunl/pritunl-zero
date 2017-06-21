package utils

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/errortypes"
	"math"
	"math/big"
	mathrand "math/rand"
	"regexp"
)

var (
	randRe = regexp.MustCompile("[^a-zA-Z0-9]+")
)

func RandStr(n int) (str string, err error) {
	input, err := RandBytes(int(math.Ceil(float64(n) * 1.5)))
	if err != nil {
		return
	}

	output := base64.RawStdEncoding.EncodeToString(input)
	output = randRe.ReplaceAllString(output, "")

	if len(output) < n {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "utils: Random generate error"),
		}
		return
	}

	str = output[:n]

	return
}

func RandBytes(size int) (bytes []byte, err error) {
	bytes = make([]byte, size)
	_, err = rand.Read(bytes)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "utils: Random read error"),
		}
		return
	}

	return
}

func init() {
	n, err := rand.Int(rand.Reader, big.NewInt(9223372036854775806))
	if err != nil {
		panic(err)
	}

	mathrand.Seed(n.Int64())
}
