package account

import (
	"errors"
	"github.com/nbutton23/zxcvbn-go"
	"net/mail"
	"regexp"
)

var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_\-.]{3,16}$`)

func ValidateLogin(username, email, password string) error {
	if err := ValidateUsername(username); err != nil {
		return err
	}
	if err := ValidateEmail(email); err != nil {
		return err
	}
	return ValidatePassword(password)
}

func ValidateUsername(username string) error {
	if username == "" {
		return errors.New("username must not be empty")
	}

	if len(username) < 3 || len(username) > 16 {
		return errors.New("username must be between 3 and 16 characters")
	}

	if !usernameRegex.MatchString(username) {
		return errors.New("invalid characters in username")
	}

	return nil
}

func ValidateEmail(email string) error {
	if email == "" {
		return errors.New("email must not be empty")
	}

	if len(email) > 32 {
		return errors.New("email cannot be longer than 32 characters")
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		return errors.New("invalid email address")
	}

	return nil
}

func ValidatePassword(password string) error {
	if password == "" {
		return errors.New("password must not be empty")
	}

	if len(password) > 32 {
		return errors.New("password cannot be longer than 32 characters")
	}

	passwordStrength := zxcvbn.PasswordStrength(password, nil).Score
	if passwordStrength < 1 {
		return errors.New("password is too weak")
	}

	return nil
}
