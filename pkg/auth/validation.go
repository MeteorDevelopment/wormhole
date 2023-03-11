package auth

import (
	"errors"
	"github.com/nbutton23/zxcvbn-go"
	"regexp"
)

var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_\-.]{3,16}$`)
var emailRegex = regexp.MustCompile(`^.+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func validateUsername(username string) error {
	if len(username) < 3 || len(username) > 16 {
		return errors.New("username must be between 3 and 16 characters")
	}

	if !usernameRegex.MatchString(username) {
		return errors.New("invalid characters in username")
	}

	return nil
}

func validateEmail(email string) error {
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}

	return nil
}

func validatePassword(password string) error {
	if len(password) > 32 {
		return errors.New("password cannot be longer than 32 characters")
	}

	passwordStrength := zxcvbn.PasswordStrength(password, nil).Score
	if passwordStrength < 1 {
		return errors.New("password is too weak")
	}

	return nil
}
