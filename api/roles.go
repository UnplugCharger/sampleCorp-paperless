package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	db "github.com/qwetu_petro/backend/db/sqlc"
	"net/http"
)

type createRoleRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description" binding:"required"`
}

type createRoleResponse struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

func (server *Server) createRole(ctx *gin.Context) {
	var req createRoleRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// TODO: check if the user is an admin

	args := db.CreateRoleParams{
		Name:        req.Name,
		Description: req.Description,
	}

	role, err := server.store.CreateRole(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, createRoleResponse{
		Name:        role.Name,
		Description: role.Description,
	})

}

type getRoleRequest struct {
	id int64 `uri:"id" binding:"required,min=1"`
}

type getRoleResponse struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

func (server *Server) getRole(ctx *gin.Context) {
	var req getRoleRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	role, err := server.store.GetRole(ctx, req.id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, getRoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
	})

}

type deleteRoleRequest struct {
	id int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteRole(ctx *gin.Context) {
	var req deleteRoleRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteRole(ctx, req.id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "role deleted successfully"})

}

type updateRoleRequest struct {
	id          int64   `uri:"id" binding:"required,min=1"`
	name        string  `json:"name"`
	description *string `json:"description"`
}

func (server *Server) updateRole(ctx *gin.Context) {
	var req updateRoleRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.UpdateRoleParams{
		ID:          req.id,
		Name:        req.name,
		Description: req.description,
	}

	role, err := server.store.UpdateRole(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, getRoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
	})

}

type listRolesRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=1"`
}

// listRoles returns all the roles in the database
func (server *Server) listRoles(ctx *gin.Context) {
	var req listRolesRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	args := db.ListRolesParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	roles, err := server.store.ListRoles(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	fmt.Println(roles)

	var response []getRoleResponse

	for _, role := range roles {
		response = append(response, getRoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
		})
	}

	ctx.JSON(http.StatusOK, response)

}
