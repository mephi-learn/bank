package pgp

import (
	"bank/pkg/errors"
	"bytes"
	"io"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
)

//var packetConfig = &packet.Config{
//	DefaultCipher: packet.CipherAES256,
//}

type Config struct {
	PGPkey string `yaml:"pgp_key" json:"pgp_key"`
}

func EncryptPGP(plaintext []byte, config *packet.Config, key string) ([]byte, error) {
	var ciphertext []byte
	encBuf := bytes.NewBuffer(nil)
	w, _ := armor.Encode(encBuf, "PGP MESSAGE", nil)
	pt, _ := openpgp.SymmetricallyEncrypt(w, []byte(key), nil, config)

	_, err := pt.Write(plaintext)
	if err != nil {
		return nil, errors.Errorf("failed encrypt PGP: %w", err)
	}

	_ = pt.Close()
	_ = w.Close()
	ciphertext = encBuf.Bytes()

	return ciphertext, nil
}

func DecryptPGP(ciphertext []byte, config *packet.Config, key string) ([]byte, error) {
	var plaintext []byte

	decBuf := bytes.NewBuffer(ciphertext)
	armorBlock, _ := armor.Decode(decBuf)

	failed := false
	prompt := func(keys []openpgp.Key, symmetric bool) ([]byte, error) {
		if failed {
			return nil, errors.New("decryption failed")
		}
		failed = true

		return []byte(key), nil
	}

	md, err := openpgp.ReadMessage(armorBlock.Body, nil, prompt, config)
	if err != nil {
		return nil, errors.Errorf("failed read PGP: %w", err)
	}

	plaintext, err = io.ReadAll(md.UnverifiedBody)
	if err != nil {
		return nil, errors.Errorf("failed decrypt PGP: %w", err)
	}

	return plaintext, nil
}
