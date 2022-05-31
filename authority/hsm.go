package authority

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/base64"
	"time"

	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/pritunl-zero/errortypes"
)

type SshRequest struct {
	Serial      string `json:"serial"`
	Certificate []byte `json:"certificate"`
}

type SshResponse struct {
	Certificate []byte `json:"certificate"`
}

type HsmStatus struct {
	Status       string `json:"status"`
	SshPublicKey string `json:"ssh_public_key"`
}

type HsmPayload struct {
	Id        string `bson:"id" json:"id"`
	Token     string `bson:"token" json:"token"`
	Nonce     string `bson:"nonce" json:"nonce"`
	Signature string `bson:"signature" json:"signature"`
	Iv        []byte `bson:"iv" json:"iv"`
	Type      string `bson:"type" json:"type"`
	Data      []byte `bson:"data" json:"data"`
}

type HsmEvent struct {
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Channel   string             `bson:"channel" json:"channel"`
	Timestamp time.Time          `bson:"timestamp" json:"timestamp"`
	Data      *HsmPayload        `bson:"data" json:"data"`
}

func (h *HsmEvent) GetId() primitive.ObjectID {
	return h.Id
}

func (h *HsmEvent) GetData() interface{} {
	if h.Data == nil {
		return nil
	}
	return h.Data
}

func UnmarshalPayload(token, secret string, payload *HsmPayload) (
	data []byte, err error) {

	if payload.Id == "" || payload.Token == "" || payload.Iv == nil ||
		payload.Data == nil {

		err = &errortypes.ParseError{
			errors.Wrap(err, "authority: Invalid payload"),
		}
		return
	}

	if subtle.ConstantTimeCompare([]byte(token),
		[]byte(payload.Token)) != 1 {

		err = &errortypes.AuthenticationError{
			errors.Wrap(err, "authority: Invalid token"),
		}
		return
	}

	if secret == "" {
		err = &errortypes.ReadError{
			errors.Wrap(err, "session: Empty secret"),
		}
		return
	}

	cipData := payload.Data
	hashFunc := hmac.New(sha512.New, []byte(secret))
	hashFunc.Write(cipData)
	rawSignature := hashFunc.Sum(nil)
	sig := base64.StdEncoding.EncodeToString(rawSignature)

	if subtle.ConstantTimeCompare([]byte(sig),
		[]byte(payload.Signature)) != 1 {

		err = &errortypes.AuthenticationError{
			errors.Wrap(err, "authority: Invalid signature"),
		}
		return
	}

	encKeyHash := sha256.New()
	encKeyHash.Write([]byte(secret))
	cipKey := encKeyHash.Sum(nil)
	cipIv := payload.Iv

	if len(cipIv) != aes.BlockSize {
		err = &errortypes.ParseError{
			errors.Wrap(err, "authority: Invalid payload iv length"),
		}
		return
	}

	if len(cipData) == 0 || len(cipData)%16 != 0 {
		err = &errortypes.ParseError{
			errors.Wrap(err, "authority: Invalid payload data length"),
		}
		return
	}

	block, err := aes.NewCipher(cipKey)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "authority: Failed to load cipher"),
		}
		return
	}

	mode := cipher.NewCBCDecrypter(block, cipIv)
	mode.CryptBlocks(cipData, cipData)
	cipData = bytes.TrimRight(cipData, "\x00")
	data = cipData

	return
}
