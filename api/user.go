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

type CreateUserResponse struct {
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

	var response CreateUserResponse
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

	ctx.JSON(http.StatusCreated, response)
}
