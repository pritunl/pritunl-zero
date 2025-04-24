package crypto

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"io"

	"github.com/pritunl/tools/errors"
	"github.com/pritunl/tools/errortypes"
	"golang.org/x/crypto/nacl/secretbox"
)

type Message struct {
	Nonce     string
	Message   string
	Signature string
}

type AsymNaclHmac struct {
	PrivateKey   *[32]byte
	Secret       *[32]byte
	nonceHandler func(nonce []byte) error
}

func (a *AsymNaclHmac) RegisterNonce(handler func(nonce []byte) error) {
	a.nonceHandler = handler
}

func (a *AsymNaclHmac) Seal(input any) (msg *Message, err error) {
	nonce := new([24]byte)
	_, err = io.ReadFull(rand.Reader, nonce[:])
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "crypto: Failed to generate nonce"),
		}
		return
	}
	nonceStr := base64.StdEncoding.EncodeToString(nonce[:])

	data, err := json.Marshal(input)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "crypto: Failed to marshal json data"),
		}
		return
	}

	encData := secretbox.Seal(nil, data, nonce, a.PrivateKey)
	encStr := base64.StdEncoding.EncodeToString(encData)

	hashFunc := hmac.New(sha512.New, a.Secret[:])
	hashFunc.Write([]byte(encStr))
	rawSignature := hashFunc.Sum(nil)
	sigStr := base64.StdEncoding.EncodeToString(rawSignature)

	msg = &Message{
		Nonce:     nonceStr,
		Message:   encStr,
		Signature: sigStr,
	}

	return
}

func (a *AsymNaclHmac) SealJson(input any) (output string, err error) {
	msg, err := a.Seal(input)
	if err != nil {
		return
	}

	outputByt, err := json.Marshal(msg)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "crypto: Failed to marshal message"),
		}
		return
	}

	output = string(outputByt)
	return
}

func (a *AsymNaclHmac) Unseal(msg *Message, output any) (err error) {
	hashFunc := hmac.New(sha512.New, a.Secret[:])
	hashFunc.Write([]byte(msg.Message))
	rawSignature := hashFunc.Sum(nil)
	sigStr := base64.StdEncoding.EncodeToString(rawSignature)

	if subtle.ConstantTimeCompare([]byte(sigStr), []byte(msg.Signature)) != 1 {
		err = &errortypes.AuthenticationError{
			errors.New("crypto: Invalid message signature"),
		}
		return
	}

	nonceByt, err := base64.StdEncoding.DecodeString(msg.Nonce)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "crypto: Failed to decode nonce"),
		}
		return
	}

	if len(nonceByt) != 24 {
		err = &errortypes.ParseError{
			errors.New("crypto: Invalid nonce length"),
		}
		return
	}

	if a.nonceHandler != nil {
		err = a.nonceHandler(nonceByt)
		if err != nil {
			err = &errortypes.ParseError{
				errors.Wrap(err, "crypto: Nonce validate failed"),
			}
			return
		}
	}

	nonce := new([24]byte)
	copy(nonce[:], nonceByt)

	encByt, err := base64.StdEncoding.DecodeString(msg.Message)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "crypto: Failed to decode message"),
		}
		return
	}

	decByt, ok := secretbox.Open(nil, encByt, nonce, a.PrivateKey)
	if !ok {
		err = &errortypes.AuthenticationError{
			errors.New("crypto: Failed to decrypt message"),
		}
		return
	}

	err = json.Unmarshal(decByt, output)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "crypto: Failed to unmarshal data"),
		}
		return
	}

	return
}

func (a *AsymNaclHmac) UnsealJson(input string, output any) (err error) {
	msg := &Message{}

	err = json.Unmarshal([]byte(input), msg)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "crypto: Failed to unmarshal message"),
		}
		return
	}

	err = a.Unseal(msg, output)
	if err != nil {
		return
	}

	return
}

func (a *AsymNaclHmac) Export() (keyStr, secrStr string) {
	keyStr = base64.StdEncoding.EncodeToString(a.PrivateKey[:])
	secrStr = base64.StdEncoding.EncodeToString(a.Secret[:])
	return
}

func (a *AsymNaclHmac) Import(keyStr, secrStr string) {
	keyByt, err := base64.StdEncoding.DecodeString(keyStr)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "crypto: Failed to decode private key"),
		}
		return
	}

	if len(keyByt) != 32 {
		err = &errortypes.ParseError{
			errors.New("crypto: Invalid private key length"),
		}
		return
	}

	secrByt, err := base64.StdEncoding.DecodeString(secrStr)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "crypto: Failed to decode secret key"),
		}
		return
	}

	if len(secrByt) != 32 {
		err = &errortypes.ParseError{
			errors.New("crypto: Invalid secret key length"),
		}
		return
	}

	if a.PrivateKey == nil {
		a.PrivateKey = new([32]byte)
	}
	if a.Secret == nil {
		a.Secret = new([32]byte)
	}

	copy(a.PrivateKey[:], keyByt)
	copy(a.Secret[:], secrByt)

	return
}

func New() (asym *AsymNaclHmac, err error) {
	privKey := new([32]byte)
	_, err = io.ReadFull(rand.Reader, privKey[:])
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "crypto: Failed to generate private key"),
		}
		return
	}

	secKey := new([32]byte)
	_, err = io.ReadFull(rand.Reader, secKey[:])
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "crypto: Failed to generate secret key"),
		}
		return
	}

	asym = &AsymNaclHmac{
		PrivateKey: privKey,
		Secret:     secKey,
	}

	return
}
