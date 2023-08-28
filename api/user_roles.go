package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	db "github.com/qwetu_petro/backend/db/sqlc"
	"net/http"
	"time"
)

type createUserRolesRequest struct {
	UserID int64 `json:"user_id" binding:"required"`
	RoleID int64 `json:"role_id" binding:"required"`
}

type createUserRolesResponse struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	RoleID    int64     `json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (server *Server) createUserRole(ctx *gin.Context) {
	var req createUserRolesRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.CreateUserRolesParams{
		UserID: req.UserID,
		RoleID: req.RoleID,
	}

	userRole, err := server.store.CreateUserRoles(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, createUserRolesResponse{
		ID:        userRole.ID,
		UserID:    userRole.UserID,
		RoleID:    userRole.RoleID,
		CreatedAt: userRole.CreatedAt,
	})

}

type deleteUserRolesRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteUserRole(ctx *gin.Context) {
	var req deleteUserRolesRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteUserRoles(ctx, req.ID)
	fmt.Println(err)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User Role deleted successfully"})
}
