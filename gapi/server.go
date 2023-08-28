package gapi

import (
	"fmt"
	"github.com/gin-gonic/gin"
	db "github.com/qwetu_petro/backend/db/sqlc"
	"github.com/qwetu_petro/backend/pb"
	"github.com/qwetu_petro/backend/token"
	"github.com/qwetu_petro/backend/utils"
	"github.com/rs/zerolog/log"
)

type Server struct {
	pb.UnimplementedQwetuBackendGrpcServiceServer
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
	config     utils.Config
}

func NewServer(config utils.Config, store db.Store) (*Server, error) {

	log.Info().Msg("Hurray  Qwetu Just Gone Live")

	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("can not create token maker: %w", err)
	}

	server := &Server{store: store, tokenMaker: tokenMaker, config: config}

	return server, nil

}
