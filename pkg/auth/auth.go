package auth

import (
	"encoding/json"
	"errors"
	"github.com/dvsekhvalnov/jose2go"
	"github.com/segmentio/ksuid"
	"golang.org/x/exp/slices"
	"log"
	"math/rand"
	"sync"
	"time"
	"wormhole/pkg/database"
)

type Claims struct {
	TokenID   int
	AccountID ksuid.KSUID
}

var jwtKey []byte
var tokenCount = 0
var tokens = make(map[ksuid.KSUID][]int)
var mu = sync.RWMutex{}

func Init() {
	jwtKey = make([]byte, 36)
	rand.Seed(time.Now().UnixNano())
	rand.Read(jwtKey)
}

func Register(username string, password string) error {
	if err := validateUsername(username); err != nil {
		log.Printf("Username validation failed: %v", err)
		return err
	}

	if err := validatePassword(password); err != nil {
		log.Printf("Password validation failed: %v", err)
		return err
	}

	_, err := database.GetAccountWithUsername(username)
	if err == nil {
		log.Printf("Account with username '%s' already exists", username)
		return errors.New("account with that username already exists")
	}

	return database.NewAccount(username, password)
}

func Login(name string, password string) (string, error) {
	if name == "" || password == "" {
		return "", errors.New("wrong name or password")
	}

	account, err := database.GetAccountWithUsername(name)
	if err != nil {
		return "", errors.New("no account with that username")
	}

	if !account.PasswordMatches(password) {
		return "", errors.New("incorrect password")
	}

	mu.Lock()

	bytes, err := json.Marshal(Claims{TokenID: tokenCount, AccountID: account.ID})
	if err != nil {
		mu.Unlock()
		return "", err
	}

	token, err := jose.Sign(string(bytes), jose.HS256, jwtKey)
	if err != nil {
		mu.Unlock()
		return "", err
	}

	tokens[account.ID] = append(tokens[account.ID], tokenCount)
	tokenCount++

	mu.Unlock()
	return token, nil
}

func Logout(token string, id ksuid.KSUID) error {
	bytes, _, err := jose.Decode(token, jwtKey)
	if err != nil {
		return err
	}

	var claims Claims
	err = json.Unmarshal([]byte(bytes), &claims)
	if err != nil {
		return err
	}

	mu.Lock()
	tokenIds, exists := tokens[id]
	if exists {
		for i := 0; i < len(tokenIds); i++ {
			if tokenIds[i] == claims.TokenID {
				tokens[id] = append(tokenIds[:i], tokenIds[i+1:]...)
				break
			}
		}
	}
	mu.Unlock()

	return nil
}

func IsTokenValid(token string) (ksuid.KSUID, error) {
	bytes, _, err := jose.Decode(token, jwtKey)
	if err != nil {
		return ksuid.Nil, err
	}

	var claims Claims
	err = json.Unmarshal([]byte(bytes), &claims)
	if err != nil {
		return ksuid.Nil, err
	}

	mu.RLock()
	validTokenIds, exists := tokens[claims.AccountID]
	mu.RUnlock()

	if exists && slices.Contains(validTokenIds, claims.TokenID) {
		return claims.AccountID, nil
	}

	return ksuid.Nil, errors.New("invalid token")
}

func Invalidate(id ksuid.KSUID) {
	delete(tokens, id)
}
