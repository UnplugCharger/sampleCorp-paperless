package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/qwetu_petro/backend/utils"

	db "github.com/qwetu_petro/backend/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

type createUserRequest struct {
	Username   string `json:"username" binding:"required,alphanum"`
	FullName   string `json:"full_name" binding:"required"`
	Email      string `json:"email" binding:"required"`
	Department string `json:"department" binding:"required"`
	Password   string `json:"password" binding:"required,min=6"`
}

type rolesResponse struct {
	RoleID      int64  `json:"role_id"`
	RoleName    string `json:"role_name"`
	Description string `json:"description"`
}

func newRolesResponse(role db.Role) rolesResponse {
	return rolesResponse{
		RoleID:      role.ID,
		RoleName:    role.Name,
		Description: *role.Description,
	}
}

type userResponse struct {
	ID                int64     `json:"id"`
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	CreatedAt         time.Time `json:"created_at"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		ID:                user.ID,
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		CreatedAt:         user.CreatedAt,
		PasswordChangedAt: user.PasswordChangedAt,
	}
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, fmt.Sprintf("invalid username or password: %s", err.Error()))
		return
	}
	password, err := utils.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	args := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: password,
		FullName:       req.FullName,
		Email:          req.Email,
		Department:     req.Department,
	}

	user, err := server.store.CreateUser(ctx, args)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusConflict, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := newUserResponse(user)
	ctx.JSON(http.StatusOK, resp)
}

type loginUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	SessionID            uuid.UUID    `json:"session_id"`
	AccessToken          string       `json:"access_token"`
	AccessTokenExpiresAt time.Time    `json:"access_token_expires_at"`
	RefreshToken         string       `json:"refresh_token"`
	RefreshTokenExpires  time.Time    `json:"refresh_token_expires_at"`
	User                 userResponse `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest

	// Bind the request body to the struct
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := req.validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.GetUserByUserNameOrEmailParams{
		Username: req.Username,
		Email:    req.Username,
	}

	user, err := server.store.GetUserByUserNameOrEmail(ctx, args)
	if err != nil {
		fmt.Println(err)
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, "user not found")
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return

	}
	err = utils.CheckPasswordHash(req.Password, user.HashedPassword)
	if err != nil {
		e := errors.New("invalid credentials")
		ctx.JSON(http.StatusUnauthorized, errorResponse(e))
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           uuid.UUID(refreshPayload.ID),
		Username:     user.Username,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		RefreshToken: refreshToken,
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiresAt,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := loginUserResponse{
		SessionID:            session.ID,
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiresAt,
		RefreshToken:         refreshToken,
		RefreshTokenExpires:  refreshPayload.ExpiresAt,
		User:                 newUserResponse(user),
	}

	ctx.JSON(http.StatusOK, resp)

}

type listUserRequest struct {
	PageID   int32 `form:"page_id" binding:"required"`
	PageSize int32 `form:"page_size" binding:"required"`
}

type listUserResponse struct {
	Users []userResponse `json:"users"`
}

func (server *Server) listUsers(ctx *gin.Context) {
	var req listUserRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.GetUsersParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	users, err := server.store.GetUsers(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	log.Info().Msgf("users: %v", users)

	var resp listUserResponse

	for _, user := range users {
		resp.Users = append(resp.Users, newUserResponse(user))
	}

	ctx.JSON(http.StatusOK, resp)
}

func (server *Server) currentUser(ctx *gin.Context) {
	User, _ := ctx.Get("user_name")
	args := db.GetUserByUserNameOrEmailParams{
		Username: User.(string),
	}
	user, err := server.store.GetUserByUserNameOrEmail(ctx, args)
	if err != nil {
		log.Error().Msgf("error getting user: %v", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	resp := newUserResponse(user)
	ctx.JSON(http.StatusOK, resp)

}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (r *loginUserRequest) validate() error {
	if strings.TrimSpace(r.Username) == "" && strings.TrimSpace(r.Email) == "" {
		return errors.New("either username or email is required")
	}

	if strings.TrimSpace(r.Password) == "" || len(r.Password) < 6 {
		return errors.New("password is required and should be at least 6 characters long")
	}

	return nil
}

func getUserIdFromContext(ctx *gin.Context) (int32, error) {
	userId, ok := ctx.Get("user_id")
	if !ok {
		log.Error().Msg("user id not found")
		return 0, errors.New("user id not found")
	}

	userIdInt, ok := userId.(int64)
	if !ok {
		return 0, errors.New("user id is not a string")
	}
	id := int32(userIdInt)

	return id, nil
}
