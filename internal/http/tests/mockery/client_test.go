package mockery_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	bHTTP "github.com/mszostok/mocks-playground/internal/http"
	"github.com/mszostok/mocks-playground/internal/http/tests/mockery/automock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDoerWithBasicAuthDo(t *testing.T) {
	tests := map[string]struct {
		expResp *http.Response
		fixErr  error
	}{
		"success": {
			expResp: fixResponse(200),
			fixErr:  nil,
		},
		"failure": {
			expResp: nil,
			fixErr:  errors.New("fix ERR"),
		},
	}

	for tn, tc := range tests {
		t.Run(tn, func(t *testing.T) {
			// given
			var (
				fixReq  = httptest.NewRequest(http.MethodGet, "/test", nil)
				fixUser = "user"
				fixPass = "pass"
			)

			clientMock := &automock.HTTPClient{}
			defer clientMock.AssertExpectations(t)
			clientMock.On("Do", mock.AnythingOfType("*http.Request")).Run(func(args mock.Arguments) {
				req := args.Get(0).(*http.Request)
				gotUsername, gotPassword, ok := req.BasicAuth()

				// then
				require.True(t, ok, "basic auth was not provided")
				assert.Equal(t, fixUser, gotUsername, "username mismatch")
				assert.Equal(t, fixPass, gotPassword, "password mismatch")
			}).Return(tc.expResp, tc.fixErr).Once()

			cli := bHTTP.NewClientWithBasicAuth(bHTTP.BasicAuth{Username: fixUser, Password: fixPass}).
				WithClient(clientMock)

			// when
			gotResp, gotErr := cli.Do(fixReq)

			// then
			assert.Equal(t, tc.expResp, gotResp)
			assert.Equal(t, tc.fixErr, gotErr)
		})
	}
}

func fixResponse(status int) *http.Response {
	respRec := httptest.NewRecorder()
	respRec.WriteHeader(status)

	return respRec.Result()
}
