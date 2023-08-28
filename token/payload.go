package token

import (
	"errors"
	"time"

	"github.com/gofrs/uuid"
)

var (
	ExpiredTokenError = errors.New("token is expired")
	InvalidTokenError = errors.New("token is invalid")
)

// TODO - add more claims to the payload such as user role and permissions
type Payload struct {
	ID        uuid.UUID `json:"id"`
	UserName  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenId, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	payload := &Payload{
		ID:        tokenId,
		UserName:  username,
		IssuedAt:  time.Now().UTC(),
		ExpiresAt: time.Now().UTC().Add(duration),
	}
	return payload, nil
}

func (p *Payload) Valid() error {
	if time.Now().UTC().After(p.ExpiresAt) {
		return ExpiredTokenError
	}
	return nil
}
