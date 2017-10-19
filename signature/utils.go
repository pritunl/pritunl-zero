package signature

import (
	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-zero/errortypes"
	"strconv"
	"time"
)

func Parse(token, sigStr, timeStr, nonce, method, path string) (
	sig *Signature, err error) {

	timestampInt, _ := strconv.ParseInt(timeStr, 10, 64)
	if timestampInt == 0 {
		err = errortypes.AuthenticationError{
			errors.New("signature: Invalid authentication timestamp"),
		}
		return
	}

	timestamp := time.Unix(timestampInt, 0)

	sig = &Signature{
		Token:     token,
		Nonce:     nonce,
		Timestamp: timestamp,
		Signature: sigStr,
		Method:    method,
		Path:      path,
	}

	return
}
