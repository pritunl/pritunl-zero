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
	Key        string
	Secret     string
	PrivateKey string
	PublicKey  string
}

type AsymNaclHmac struct {
	key          *[32]byte
	secret       *[32]byte
	publicKey    *[32]byte
	privateKey   *[64]byte
	nonceHandler func(nonce []byte) error
}

func (a *AsymNaclHmac) RegisterNonce(handler func(nonce []byte) error) {
	a.nonceHandler = handler
}

func (a *AsymNaclHmac) Seal(input any) (msg *Message, err error) {
	if a.key == nil || a.secret == nil || a.privateKey == nil {
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

	encByt := secretbox.Seal(nil, data, nonce, a.key)
	sigEncByt := sign.Sign(nil, encByt, a.privateKey)
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
	if a.key == nil || a.secret == nil || a.publicKey == nil {
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

	encByt, valid := sign.Open(nil, sigEncByt, a.publicKey)
	if !valid {
		err = &errortypes.ParseError{
			errors.Wrap(err, "crypto: Failed to verify message signature"),
		}
		return
	}

	decByt, ok := secretbox.Open(nil, encByt, nonce, a.key)
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
		Key:        base64.StdEncoding.EncodeToString(a.key[:]),
		Secret:     base64.StdEncoding.EncodeToString(a.secret[:]),
		PublicKey:  base64.StdEncoding.EncodeToString(a.publicKey[:]),
		PrivateKey: base64.StdEncoding.EncodeToString(a.privateKey[:]),
	}
}

func (a *AsymNaclHmac) Import(key AsymNaclHmacKey) (err error) {
	keyByt, err := base64.StdEncoding.DecodeString(key.Key)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "crypto: Failed to decode key"),
		}
		return
	}

	if len(keyByt) != 32 {
		err = &errortypes.ParseError{
			errors.New("crypto: Invalid key length"),
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

	pubKeyByt, err := base64.StdEncoding.DecodeString(
		key.PublicKey)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "crypto: Failed to decode public key"),
		}
		return
	}

	if len(pubKeyByt) != 32 {
		err = &errortypes.ParseError{
			errors.New("crypto: Invalid public key length"),
		}
		return
	}

	if key.PrivateKey != "" {
		privKeyByt, e := base64.StdEncoding.DecodeString(
			key.PrivateKey)
		if e != nil {
			err = &errortypes.ParseError{
				errors.Wrap(e, "crypto: Failed to decode private key"),
			}
			return
		}

		if len(privKeyByt) != 64 {
			err = &errortypes.ParseError{
				errors.New("crypto: Invalid private key length"),
			}
			return
		}

		if a.privateKey == nil {
			a.privateKey = new([64]byte)
		}
		copy(a.privateKey[:], privKeyByt)
	}

	if a.key == nil {
		a.key = new([32]byte)
	}
	if a.secret == nil {
		a.secret = new([32]byte)
	}
	if a.publicKey == nil {
		a.publicKey = new([32]byte)
	}

	copy(a.key[:], keyByt)
	copy(a.secret[:], secrByt)
	copy(a.publicKey[:], pubKeyByt)

	return
}

func (a *AsymNaclHmac) Generate() (err error) {
	key := new([32]byte)
	_, err = io.ReadFull(rand.Reader, key[:])
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "crypto: Failed to generate key"),
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

	a.key = key
	a.secret = secKey
	a.publicKey = signPubKey
	a.privateKey = signPrivKey

	return
}
