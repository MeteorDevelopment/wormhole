package account

import (
	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	Id           uint64
	CreatedAt    int64
	Username     string
	Email        string
	ProfileImage string
	Password     []byte
	Admin        bool
}

func (a *Account) PasswordMatches(password string) bool {
	return bcrypt.CompareHashAndPassword(a.Password, []byte(password)) == nil
}
