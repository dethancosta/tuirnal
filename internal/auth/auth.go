package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"

	"github.com/dethancosta/tuirnal/internal/helpers"
	"github.com/dethancosta/tuirnal/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func Authenticate(app helpers.Application, username, password string) (*models.Author, error) {
	am := app.AuthorModel
	user, err := am.Get(username)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("can't authenticate: %w", err)
	}

	return user, nil
}

func EncryptEntry(password string, entryJson []byte) ([]byte, error) {
	key := []byte(password)

	cipherBlock, err := aes.NewCipher(key)

	if err != nil {
		return nil, fmt.Errorf("Couldn't encrypt entry: %w", err)
	}

	cipherText := make([]byte, aes.BlockSize+len(entryJson))

	iv := cipherText[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, fmt.Errorf("Couldn't encrypt entry: %w", err)
	}

	stream := cipher.NewCFBDecrypter(cipherBlock, iv)

	stream.XORKeyStream(cipherText[aes.BlockSize:], entryJson)

	return cipherText, nil
}

func DecryptEntry(password string, encryptedJson []byte) ([]byte, error) {
	key := []byte(password)

	cipherBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("Can't decrypt entry: %w", err)
	}

	if len(encryptedJson) < aes.BlockSize {
		return nil, fmt.Errorf("Text is too short")
	}

	iv := encryptedJson[:aes.BlockSize]

	encryptedJson = encryptedJson[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(cipherBlock, iv)

	stream.XORKeyStream(encryptedJson, encryptedJson)

	return encryptedJson, nil
}
