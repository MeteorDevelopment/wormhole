package jwt

import (
	"testing"
	"time"
	"wormhole/pkg/snowflake"
)

type RandomClaims struct {
	Id        uint64 `json:"id"`
	CreatedAt int64  `json:"time"`
	Username  string `json:"username"`
}

func TestJwt(t *testing.T) {
	SigningKey = []byte("2hmr82u39nxyet2r693b78ny9pu1m2n")
	ExpirationTime = 1 * 60

	claims := &RandomClaims{
		Id:        snowflake.NextId(),
		CreatedAt: time.Now().Add(30 * time.Minute).Unix(),
		Username:  "test",
	}

	token, err := MakeToken(claims)
	if err != nil {
		t.Error(err)
	}

	var parsed RandomClaims
	err = ParseToken(token, &parsed)
	if err != nil {
		t.Error(err)
	}

	if claims.Id != parsed.Id || claims.CreatedAt != parsed.CreatedAt || claims.Username != parsed.Username {
		t.Error("parsed data mismatch")
	}

	if !IsTokenExpired(parsed.CreatedAt) {
		t.Error("token should be expired")
	}
}
