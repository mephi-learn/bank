package models

import (
	"bank/pkg/errors"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"time"
)

const cardSecretKey = "однажды в студёную зимнюю пору"

type Card struct {
	ID        uint      `json:"id"`         // Идентификатор карты
	ClientID  uint      `json:"-"`          // Идентификатор клиента
	AccountID uint      `json:"account_id"` // Привязано к счёту
	Number    string    `json:"number"`     // Номер карты
	Expire    time.Time `json:"-"`          // Время истекания карты
	CVV       string    `json:"-"`          // CVV карты
}

func (c *Card) CryptNum() (string, error) {
	hash := c.hash()
	block, err := aes.NewCipher(hash)
	if err != nil {
		return "", err
	}
	ciphertext := make([]byte, aes.BlockSize+len(c.Number))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(c.Number))

	return hex.EncodeToString(ciphertext), nil
}

func (c *Card) CryptCVV() (string, error) {
	hash := c.hash()
	block, err := aes.NewCipher(hash)
	if err != nil {
		return "", err
	}
	ciphertext := make([]byte, aes.BlockSize+len(c.Number))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(c.Number))

	return hex.EncodeToString(ciphertext), nil
}

func (c *Card) DecryptCVV(number string) error {
	hash := c.hash()

	decoded, err := hex.DecodeString(number)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(hash)
	if err != nil {
		return err
	}

	if len(decoded) < aes.BlockSize {
		return errors.New("ciphertext too short")
	}

	iv := decoded[:aes.BlockSize]
	decoded = decoded[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(decoded, decoded)
	c.Number = fmt.Sprintf("%s", decoded)

	return nil
}

func (c *Card) Expired() bool {
	return time.Now().After(c.Expire)
}

func (c *Card) hash() []byte {
	mac := hmac.New(sha512.New384, []byte(cardSecretKey))

	return mac.Sum(nil)[:32]
}
