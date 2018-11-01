package http_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	bHTTP "github.com/mszostok/mocks-playground/http"
	"github.com/mszostok/mocks-playground/http/tests/pegomock/automock"
	. "github.com/petergtz/pegomock"
	"github.com/stretchr/testify/assert"
)

func TestDoerWithBasicAuthDo(t *testing.T) {
	RegisterMockTestingT(t)
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

			clientMock := automock.NewMockHTTPClient()
			When(clientMock.Do(EqReqAuth(fixUser, fixPass))).
				ThenReturn(tc.expResp, tc.fixErr)

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

func EqReqAuth(user, pass string) *http.Request {
	RegisterMatcher(&ReqAuthMatcher{expUsername: user, expPassword: pass})
	var nullValue *http.Request
	return nullValue
}

type ReqAuthMatcher struct {
	expUsername, expPassword string
	gotUsername, gotPassword string
	sync.Mutex
}

func (m *ReqAuthMatcher) Matches(param Param) bool {
	m.Lock()
	defer m.Unlock()

	req := param.(*http.Request)
	var credsProvided bool
	m.gotUsername, m.gotPassword, credsProvided = req.BasicAuth()

	return credsProvided &&
		m.expPassword == m.gotPassword && m.expUsername == m.gotUsername
}

func (m *ReqAuthMatcher) FailureMessage() string {
	return fmt.Sprintf("Expected user: %s and password:  but got user: %s and password: %s",
		m.expUsername, m.gotUsername, m.gotPassword)
}

func (m *ReqAuthMatcher) String() string {
	return fmt.Sprintf("Eq(%s,%s)", m.expUsername, m.expPassword)
}
