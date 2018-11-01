package http_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	bHTTP "github.com/mszostok/mocks-playground/http"
	"github.com/mszostok/mocks-playground/http/tests/gomock/automock"
	"github.com/stretchr/testify/assert"
)

func TestDoerWithBasicAuthDo_use_gomock(t *testing.T) {
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

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockObj := automock.NewMockHTTPClient(mockCtrl)
			mockObj.EXPECT().
				Do(EqReqAuthMatcher(fixUser, fixPass)).
				Return(tc.expResp, tc.fixErr)

			cli := bHTTP.NewClientWithBasicAuth(bHTTP.BasicAuth{Username: fixUser, Password: fixPass}).
				WithClient(mockObj)

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

func EqReqAuthMatcher(user, pass string) *ReqAuthMatcher {
	return &ReqAuthMatcher{user, pass}
}

type ReqAuthMatcher struct {
	expUsername, expPassword string
}

func (m *ReqAuthMatcher) Matches(in interface{}) bool {
	req := in.(*http.Request)
	gotUser, gotPass, credsProvided := req.BasicAuth()

	return credsProvided &&
		m.expPassword == gotPass && m.expUsername == gotUser
}

func (m *ReqAuthMatcher) String() string {
	return fmt.Sprintf("basic auth is set [user: %s pass: %s]", m.expUsername, m.expPassword)
}
