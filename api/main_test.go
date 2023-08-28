package api

import (
	"github.com/hibiken/asynq"
	"github.com/qwetu_petro/backend/workers"
	"os"
	"testing"
	"time"

	db "github.com/qwetu_petro/backend/db/sqlc"
	"github.com/qwetu_petro/backend/utils"
	"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := utils.Config{
		TokenSymmetricKey:   utils.RandomString(32),
		AccessTokenDuration: time.Minute,
	}
	redisOpts := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	taskDistributor := workers.NewRedisTaskDistributor(redisOpts)
	server, err := NewServer(config, store, taskDistributor)
	require.NoError(t, err)
	return server

}

func TestMain(m *testing.M) {
	err := os.Setenv("TEST_ENVIRONMENT", "true")
	if err != nil {
		log.Fatal().Msg("cannot set test environment")
	}
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
