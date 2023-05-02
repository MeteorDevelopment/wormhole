package account

import (
	"github.com/pkg/errors"
	"testing"
)

func TestValidateUsername(t *testing.T) {
	cases := []struct {
		username string
		err      error
	}{
		{username: "user123", err: nil},
		{username: "u$ername", err: errors.New("invalid characters in username")},
		{username: "aa", err: errors.New("username must be between 3 and 16 characters")},
		{username: "averylongusername111", err: errors.New("username must be between 3 and 16 characters")},
		{username: "", err: errors.New("username must not be empty")},
	}

	for _, c := range cases {
		err := ValidateUsername(c.username)
		if (err != nil && c.err == nil) || (err == nil && c.err != nil) || (err != nil && c.err != nil && err.Error() != c.err.Error()) {
			t.Errorf("for username %q, expected error %v but got error %v", c.username, c.err, err)
		}
	}
}

func TestValidatePassword(t *testing.T) {
	cases := []struct {
		password string
		err      error
	}{
		{password: "bw>.893!@#", err: nil},
		{password: "password123", err: errors.New("password is too weak")},
		{password: "averylongpasswordthatismorethan32characterslong", err: errors.New("password cannot be longer than 32 characters")},
		{password: "", err: errors.New("password must not be empty")},
	}

	for _, c := range cases {
		err := ValidatePassword(c.password)
		if (err != nil && c.err == nil) || (err == nil && c.err != nil) || (err != nil && c.err != nil && err.Error() != c.err.Error()) {
			t.Errorf("for password %q, expected error %v but got error %v", c.password, c.err, err)
		}
	}
}
