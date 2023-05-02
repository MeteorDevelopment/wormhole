package jwt

import (
	"encoding/json"
	jose "github.com/dvsekhvalnov/jose2go"
	"github.com/pkg/errors"
	"time"
	"wormhole/pkg/config"
)

var (
	SigningKey     []byte
	ExpirationTime int64
)

func Init() {
	SigningKey = []byte(config.Get().JwtSigningKey)
	ExpirationTime = config.Get().JwtExpirationTime * 60
}

func MakeToken(claims any) (string, error) {
	claimsJson, err := json.Marshal(claims)
	if err != nil {
		return "", errors.Wrap(err, "error marshaling claims")
	}

	token, err := jose.SignBytes(claimsJson, jose.HS256, SigningKey)
	if err != nil {
		return "", errors.Wrap(err, "error encrypting token")
	}

	return token, nil
}

func ParseToken(token string, destination any) error {
	claimsJson, _, err := jose.DecodeBytes(token, SigningKey)
	if err != nil {
		return errors.Wrap(err, "error decrypting token")
	}

	err = json.Unmarshal(claimsJson, destination)
	if err != nil {
		return errors.Wrap(err, "error unmarshaling claims")
	}

	return nil
}

func IsTokenExpired(createdAt int64) bool {
	return createdAt+ExpirationTime > time.Now().Unix()
}
