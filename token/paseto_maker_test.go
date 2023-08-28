package token

import (
	"errors"
	"testing"
	"time"

	"github.com/qwetu_petro/backend/utils"

	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {

	_, err := NewPasetoMaker(utils.RandomString(16))
	require.Error(t, err)

	maker, err := NewJwtMaker(utils.RandomString(32))
	require.NoError(t, err)

	username := utils.RandomUser()
	duration := time.Minute

	createdAt := time.Now()
	Expires := createdAt.Add(duration)

	token, payload, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.Equal(t, username, payload.UserName)
	require.True(t, createdAt.Before(payload.IssuedAt) || createdAt.Equal(payload.IssuedAt))
	require.True(t, Expires.Before(payload.ExpiresAt) || Expires.Equal(payload.ExpiresAt))
	require.NotZero(t, payload.ID)

}

func TestPasetoExpired(t *testing.T) {

	maker, err := NewPasetoMaker(utils.RandomString(32))
	require.NoError(t, err)

	token, payload, err := maker.CreateToken(utils.RandomUser(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.True(t, errors.Is(err, ExpiredTokenError))
	require.Nil(t, payload)

}
