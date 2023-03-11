package database

import (
	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	ID       ksuid.KSUID `bson:"id" json:"id"`
	Username string      `bson:"username" json:"username"`
	Password []byte      `bson:"password" json:"-"`
	Admin    bool        `bson:"admin" json:"-"`
}

func NewAccount(username string, password string) error {
	pass, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return err
	}

	_, err = accounts.InsertOne(nil, Account{
		ID:       ksuid.New(),
		Username: username,
		Password: pass,
		Admin:    false,
	})

	return err
}

func (a *Account) String() string {
	return a.Username
}

func (a *Account) PasswordMatches(password string) bool {
	return bcrypt.CompareHashAndPassword(a.Password, []byte(password)) == nil
}
