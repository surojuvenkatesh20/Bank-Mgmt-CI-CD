package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/surojuvenkatesh20/bank-mgmt/db/sqlc"
	"github.com/surojuvenkatesh20/bank-mgmt/utils"
)

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type UserResponse struct {
	ID                int64        `json:"id"`
	Username          string       `json:"username"`
	FullName          string       `json:"full_name"`
	Email             string       `json:"email"`
	CreatedAt         sql.NullTime `json:"created_at"`
	PasswordChangedAt time.Time    `json:"password_changed_at"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req CreateUserRequest
	var err error
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	req.Password, err = utils.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: req.Password,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		fmt.Println(err)
		if pqErr, ok := err.(*pq.Error); ok {
			fmt.Println(pqErr.Code.Name())
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusBadRequest, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var response UserResponse
	response = mapDBToResponseStruct(user, response)

	ctx.JSON(http.StatusCreated, response)
}

type LoginUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginUserResponse struct {
	AccessToken string       `json:"access_token"`
	UserDetails UserResponse `json:"user_details"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req LoginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = utils.CheckPassword(user.HashedPassword, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	token, err := server.tokenMaker.CreateToken(req.Username, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	var userDetails UserResponse
	userDetails = mapDBToResponseStruct(user, userDetails)
	response := LoginUserResponse{
		AccessToken: token,
		UserDetails: userDetails,
	}

	ctx.JSON(http.StatusOK, response)

}

func mapDBToResponseStruct(user db.User, response UserResponse) UserResponse {
	dbVal := reflect.ValueOf(&user).Elem()
	dbType := dbVal.Type()

	responseVal := reflect.ValueOf(&response).Elem()

	for i := 0; i < dbVal.NumField(); i++ {
		field := dbVal.Field(i)
		fieldName := dbType.Field(i).Name

		if field.IsValid() && !field.IsZero() {
			responseField := responseVal.FieldByName(fieldName)
			if responseField.CanSet() {
				responseField.Set(field)
			}
		}
	}
	return response
}
