package token

import "time"

type Maker interface {
	// CreateToken creates a new token for the given username and signs it with the secret.
	CreateToken(username string, duration time.Duration) (string, *Payload, error)

	// VerifyToken verifies the given token and returns the payload if the token is valid.
	VerifyToken(token string) (*Payload, error)
}
