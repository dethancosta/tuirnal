package auth

import (
	"fmt"

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
