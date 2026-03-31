package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockdb "github.com/surojuvenkatesh20/bank-mgmt/db/mock"
	db "github.com/surojuvenkatesh20/bank-mgmt/db/sqlc"
	"github.com/surojuvenkatesh20/bank-mgmt/utils"
)

func TestGetAccountAPI(t *testing.T) {
	account := getRandomAccount()

	testCases := []struct {
		name       string
		accountID  int64
		buildStubs func(store *mockdb.MockStore)
		//building stubs - expectation is that if /accounts/id is called,
		//then GetAccount(ctx, id) will invoke 1 time and send response
		//of account and nil error.
		checkResponse func(t *testing.T, w *httptest.ResponseRecorder)
	}{{
		name:      "OK",
		accountID: account.ID,
		buildStubs: func(mock *mockdb.MockStore) {
			mock.EXPECT().
				GetAccount(gomock.Any(), gomock.Eq(account.ID)).
				Times(1).
				Return(account, nil)
		},
		//server serves the request and stores the response in recorder
		checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			require.Equal(t, http.StatusOK, recorder.Code)
			requireBodyMatchAccount(t, recorder.Body, account)
		},
	}, {
		name:      "NotFound",
		accountID: account.ID,
		buildStubs: func(mock *mockdb.MockStore) {
			mock.EXPECT().
				GetAccount(gomock.Any(), gomock.Eq(account.ID)).
				Times(1).
				Return(db.Account{}, sql.ErrNoRows)
		},
		//server serves the request and stores the response in recorder
		checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			require.Equal(t, http.StatusNotFound, recorder.Code)
		},
	},
		{
			name:      "InternalError",
			accountID: account.ID,
			buildStubs: func(mock *mockdb.MockStore) {
				mock.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			//server serves the request and stores the response in recorder
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "BadRequest",
			accountID: 0,
			buildStubs: func(mock *mockdb.MockStore) {
				mock.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			//server serves the request and stores the response in recorder
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		}}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mock := mockdb.NewMockStore(ctrl)

			//build stubs
			tc.buildStubs(mock)

			//creating a server with mock and calling //accounts/id url
			server := NewServer(mock)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tc.accountID)

			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}

}

func getRandomAccount() (account db.Account) {
	return db.Account{
		Owner:    utils.GenerateRandomOwner(),
		ID:       utils.GenerateRandomInt(1, 1000),
		Currency: utils.GenerateRandomCurrency(),
		Balance:  utils.GenerateRandomMoney(),
		// CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)
	fmt.Println(body)
	fmt.Println(string(data))

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	fmt.Println(err)
	require.NoError(t, err)
	require.Equal(t, gotAccount, account)
}

func TestCreateAccountAPI(t *testing.T) {
	account := getRandomAccount()
	testCases := []struct {
		name          string
		body          gin.H
		buildstubs    func(mock *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "CREATED",
			body: gin.H{
				"owner":    account.Owner,
				"currency": account.Currency,
			},
			buildstubs: func(mock *mockdb.MockStore) {
				arg := db.CreateAccountParams{
					Owner:    account.Owner,
					Currency: account.Currency,
					Balance:  0,
				}
				mock.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.NotEmpty(t, recorder.Body)
				require.Equal(t, http.StatusCreated, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name: "BadRequest-InvalidCurrency",
			body: gin.H{
				"owner":    account.Owner,
				"currency": "INR",
			},
			buildstubs: func(mock *mockdb.MockStore) {
				mock.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.NotEmpty(t, recorder.Body)
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "BadRequest-InvalidOwner",
			body: gin.H{
				"owner":    "",
				"currency": "EUR",
			},
			buildstubs: func(mock *mockdb.MockStore) {
				mock.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.NotEmpty(t, recorder.Body)
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			body: gin.H{
				"owner":    account.Owner,
				"currency": account.Currency,
			},
			buildstubs: func(mock *mockdb.MockStore) {
				arg := db.CreateAccountParams{
					Owner:    account.Owner,
					Currency: account.Currency,
					Balance:  0,
				}
				mock.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.NotEmpty(t, recorder.Body)
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mock := mockdb.NewMockStore(ctrl)

			//build stubs
			tc.buildstubs(mock)
			server := NewServer(mock)
			recorder := httptest.NewRecorder()

			url := "/accounts"

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}

}
