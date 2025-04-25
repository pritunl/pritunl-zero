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
	"golang.org/x/crypto/nacl/sign"
)

type Message struct {
	Nonce     string
	Message   string
	Signature string
}

type AsymNaclHmacKey struct {
	PrivateKey        string
	Secret            string
	SigningPrivateKey string
	SigningPublicKey  string
}

type AsymNaclHmac struct {
	privateKey     *[32]byte
	secret         *[32]byte
	signPublicKey  *[32]byte
	signPrivateKey *[64]byte
	nonceHandler   func(nonce []byte) error
}

func (a *AsymNaclHmac) RegisterNonce(handler func(nonce []byte) error) {
	a.nonceHandler = handler
}

func (a *AsymNaclHmac) Seal(input any) (msg *Message, err error) {
	if a.privateKey == nil || a.secret == nil || a.signPrivateKey == nil {
		err = &errortypes.AuthenticationError{
			errors.New("crypto: Private key and secret not loaded"),
		}
		return
	}

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

	encByt := secretbox.Seal(nil, data, nonce, a.privateKey)
	sigEncByt := sign.Sign(nil, encByt, a.signPrivateKey)
	sigEncStr := base64.StdEncoding.EncodeToString(sigEncByt)

	hashFunc := hmac.New(sha512.New, a.secret[:])
	hashFunc.Write([]byte(sigEncStr))
	rawSignature := hashFunc.Sum(nil)
	sigStr := base64.StdEncoding.EncodeToString(rawSignature)

	msg = &Message{
		Nonce:     nonceStr,
		Message:   sigEncStr,
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
	if a.privateKey == nil || a.secret == nil || a.signPublicKey == nil {
		err = &errortypes.AuthenticationError{
			errors.New("crypto: Private key and secret not loaded"),
		}
		return
	}

	hashFunc := hmac.New(sha512.New, a.secret[:])
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

	sigEncByt, err := base64.StdEncoding.DecodeString(msg.Message)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "crypto: Failed to decode message"),
		}
		return
	}

	encByt, valid := sign.Open(nil, sigEncByt, a.signPublicKey)
	if !valid {
		err = &errortypes.ParseError{
			errors.Wrap(err, "crypto: Failed to verify message signature"),
		}
		return
	}

	decByt, ok := secretbox.Open(nil, encByt, nonce, a.privateKey)
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

func (a *AsymNaclHmac) Export() AsymNaclHmacKey {
	return AsymNaclHmacKey{
		PrivateKey: base64.StdEncoding.EncodeToString(
			a.privateKey[:]),
		Secret: base64.StdEncoding.EncodeToString(
			a.secret[:]),
		SigningPublicKey: base64.StdEncoding.EncodeToString(
			a.signPublicKey[:]),
		SigningPrivateKey: base64.StdEncoding.EncodeToString(
			a.signPrivateKey[:]),
	}
}

func (a *AsymNaclHmac) Import(key AsymNaclHmacKey) (err error) {
	privKeyByt, err := base64.StdEncoding.DecodeString(key.PrivateKey)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "crypto: Failed to decode private key"),
		}
		return
	}

	if len(privKeyByt) != 32 {
		err = &errortypes.ParseError{
			errors.New("crypto: Invalid private key length"),
		}
		return
	}

	secrByt, err := base64.StdEncoding.DecodeString(key.Secret)
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

	signPubKeyByt, err := base64.StdEncoding.DecodeString(
		key.SigningPublicKey)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "crypto: Failed to decode signing public key"),
		}
		return
	}

	if len(signPubKeyByt) != 32 {
		err = &errortypes.ParseError{
			errors.New("crypto: Invalid signing public key length"),
		}
		return
	}

	if key.SigningPrivateKey != "" {
		signPrivKeyByt, e := base64.StdEncoding.DecodeString(
			key.SigningPrivateKey)
		if e != nil {
			err = &errortypes.ParseError{
				errors.Wrap(e, "crypto: Failed to decode signing private key"),
			}
			return
		}

		if len(signPrivKeyByt) != 64 {
			err = &errortypes.ParseError{
				errors.New("crypto: Invalid signing private key length"),
			}
			return
		}

		if a.signPrivateKey == nil {
			a.signPrivateKey = new([64]byte)
		}
		copy(a.signPrivateKey[:], signPrivKeyByt)
	}

	if a.privateKey == nil {
		a.privateKey = new([32]byte)
	}
	if a.secret == nil {
		a.secret = new([32]byte)
	}
	if a.signPublicKey == nil {
		a.signPublicKey = new([32]byte)
	}

	copy(a.privateKey[:], privKeyByt)
	copy(a.secret[:], secrByt)
	copy(a.signPublicKey[:], signPubKeyByt)

	return
}

func (a *AsymNaclHmac) Generate() (err error) {
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

	signPubKey, signPrivKey, err := sign.GenerateKey(rand.Reader)
	if err != nil {
		err = &errortypes.ReadError{
			errors.Wrap(err, "crypto: Failed to generate signing key"),
		}
		return
	}

	a.privateKey = privKey
	a.secret = secKey
	a.signPublicKey = signPubKey
	a.signPrivateKey = signPrivKey

	return
}
