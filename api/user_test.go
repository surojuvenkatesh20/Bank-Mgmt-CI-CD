package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	mockdb "github.com/surojuvenkatesh20/bank-mgmt/db/mock"
	db "github.com/surojuvenkatesh20/bank-mgmt/db/sqlc"
	"github.com/surojuvenkatesh20/bank-mgmt/utils"
)

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	// In case, some value is nil
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := utils.CheckPassword(arg.HashedPassword, e.password)
	if err != nil {
		return false
	}
	e.arg.HashedPassword = arg.HashedPassword

	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("is equal to %v (%T)", e.arg, e.password)
}

func EqCreateUserMatcher(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}
func TestCreateUser(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildstubs    func(mock *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Created",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildstubs: func(mock *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
				}
				//build stubs
				mock.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserMatcher(arg, password)).
					Times(1).
					Return(user, nil)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.NotEmpty(t, recorder.Body)
				require.Equal(t, http.StatusCreated, recorder.Code)
				requireUserBodyMatch(t, recorder.Body, user)
			},
		},
		{
			name: "Invalid Username",
			body: gin.H{
				"username":  "venkate!@#@!#!",
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildstubs: func(mock *mockdb.MockStore) {
				//build stubs
				mock.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
				// .Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.NotEmpty(t, recorder.Body)
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				// requireUserBodyMatch(t, recorder.Body, user)
			},
		},
		{
			name: "Invalid Password",
			body: gin.H{
				"username":  user.Username,
				"password":  utils.GenerateRandomString(5),
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildstubs: func(mock *mockdb.MockStore) {
				//build stubs
				mock.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
				// .Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.NotEmpty(t, recorder.Body)
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				// requireUserBodyMatch(t, recorder.Body, user)
			},
		},
		{
			name: "Invalid Email",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     "sampleemail.com",
			},
			buildstubs: func(mock *mockdb.MockStore) {
				//build stubs
				mock.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
				// .Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.NotEmpty(t, recorder.Body)
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				// requireUserBodyMatch(t, recorder.Body, user)
			},
		},
		{
			name: "Duplicate Username",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildstubs: func(mock *mockdb.MockStore) {
				//build stubs
				mock.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, &pq.Error{Code: "23505"})
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.NotEmpty(t, recorder.Body)
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				// requireUserBodyMatch(t, recorder.Body, user)
			},
		},
		{
			name: "Internal Server Error",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildstubs: func(mock *mockdb.MockStore) {
				//build stubs
				mock.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.NotEmpty(t, recorder.Body)
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
				// requireUserBodyMatch(t, recorder.Body, user)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mock := mockdb.NewMockStore(ctrl)
			//build stubs
			tc.buildstubs(mock)

			server := NewServer(mock)
			recorder := httptest.NewRecorder()
			url := "/users"

			//body is request body sent by user in postman
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
		})
	}

}

func randomUser(t *testing.T) (db.User, string) {
	password := utils.GenerateRandomString(6)
	hashedPassword, err := utils.HashPassword(password)
	require.NoError(t, err)

	user := db.User{
		ID:             utils.GenerateRandomInt(1, 1000),
		Username:       utils.GenerateRandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       utils.GenerateRandomString(10),
		Email:          utils.GenerateRandomEmail(),
	}
	return user, password
}

func requireUserBodyMatch(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)

	require.Equal(t, user.ID, gotUser.ID)
	require.Equal(t, user.Username, gotUser.Username)
	require.Equal(t, user.Email, gotUser.Email)
	require.Equal(t, user.FullName, gotUser.FullName)
	// require.NotZero(t, gotUser.HashedPassword)
	// require.NotZero(t, gotUser.CreatedAt)
	// require.NotZero(t, gotUser.PasswordChangedAt)
}

//arg refers to the struct that server sends to sqldb after encoding from request body
