package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	val10 "github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"example.com/back/auth/internal/mocks/mocks_v1_login"
	"example.com/back/auth/internal/mocks/mocks_v1_users_post"
	"example.com/back/auth/internal/user"
)

// nolint:
func TestCreateHandler(t *testing.T) {
	type fields struct {
		uc userManager
		tm tokenGetter
		v  requestValidator
		sv sessionsVerificationChecker
		sd sessionsDeleter
	}

	mock := gomock.NewController(t)
	defer mock.Finish()

	uc := mocks_v1_users_post.NewMockuserManager(mock)
	ucCall := uc.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any(),
		gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	ucCall.DoAndReturn(func(ctx context.Context, name, email, password, phone string, role user.Role) (*user.User, error) {
		switch name {
		case "create user duplicate error":
			return nil, user.ErrDuplicateKey
		case "create user other DB error":
			return nil, errors.New("testing other DB error in hanndler")
		default:
			return &user.User{
				ID:     123,
				Role:   user.RoleUser,
				Status: user.StatusActive,
			}, nil
		}
	})

	tm := mocks_v1_login.NewMocktokenGetter(mock)
	tm.EXPECT().GetAccessToken(gomock.Any()).Return("access token").AnyTimes()
	tm.EXPECT().GetRefreshToken(gomock.Any()).Return("refresh token").AnyTimes()

	sv := mocks_v1_users_post.NewMocksessionsVerificationChecker(mock)
	svCall := sv.EXPECT().IsSessionValid(gomock.Any(), gomock.Any()).AnyTimes()
	svCall.DoAndReturn(func(ctx context.Context, sessionID string) bool {
		return sessionID != "d2067508-e490-4bd7-a13a-a8b7c28aa7e0"
	})

	sd := mocks_v1_users_post.NewMocksessionsDeleter(mock)
	sdCall := sd.EXPECT().Delete(gomock.Any(), gomock.Any()).AnyTimes()
	sdCall.DoAndReturn(func(ctx context.Context, sessionID string) error {
		if sessionID == "9a2f99e0-f304-4905-b175-2a664707c2a0" {
			return errors.New("testing delete session from DB error")
		}
		return nil
	})

	v := val10.New()

	def := fields{
		uc: uc,
		tm: tm,
		v:  v,
		sv: sv,
		sd: sd,
	}
	tests := []struct {
		name             string
		expectedHTTPCode int
		fields           fields
		payload          string
	}{
		{"passing not JSON body", http.StatusBadRequest, def, "it's not a JSON"},
		{"validation err: missing required field phone", http.StatusBadRequest, def, `
		{
		    "name": "missing phone user",
			"email":"example@example.ru",
			"sessionId": "5212874b-1990-4360-b2f3-09e0728e1fb4",
		    "role": "client"
		}`,
		},
		{"validation err: missing required field sessionId", http.StatusBadRequest, def, `
		{
		    "name": "missing sessionId user",
			"email":"example@example.ru",
			"phone":"+79779814062",
		    "role": "client"
		}`,
		},
		{"validation err: missing required field role", http.StatusBadRequest, def, `
		{
		    "name": "missing role user",
			"email":"example@example.ru",
			"phone":"+79779814062",
		    "sessionId": "5212874b-1990-4360-b2f3-09e0728e1fb4"
		}`,
		},
		{"validation err: bad format email", http.StatusBadRequest, def, `
		{
		    "name": "bad format email",
			"email":"example",
			"phone":"+79779814062",
		    "sessionId": "5212874b-1990-4360-b2f3-09e0728e1fb4",
			"role": "client"
		}`,
		},
		{"validation err: bad format phone", http.StatusBadRequest, def, `
		{
		    "name": "bad format phone",
			"email":"example@email.ru",
			"phone":"bad format phone",
		    "sessionId": "5212874b-1990-4360-b2f3-09e0728e1fb4",
			"role": "client"
		}`,
		},
		{"validation err: bad format sessionId", http.StatusBadRequest, def, `
		{
		    "name": "bad format sessionId",
			"email":"example@email.ru",
			"phone":"+79779814062",
		    "sessionId": "this is my sessionID",
			"role": "client"
		}`,
		},
		{"passed not valid sessionId", http.StatusForbidden, def, `
		{
		    "name": "passed not valid sessionId",
			"email":"example@email.ru",
			"phone":"+79779814062",
		    "sessionId": "d2067508-e490-4bd7-a13a-a8b7c28aa7e0",
			"role": "client"
		}`,
		},
		{"create user duplicate error", http.StatusConflict, def, `
		{
		    "name": "create user duplicate error",
			"email":"example@email.ru",
			"phone":"+79779814062",
		    "sessionId": "5212874b-1990-4360-b2f3-09e0728e1fb4",
			"role": "client"
		}`,
		},
		{"create user other DB error", http.StatusInternalServerError, def, `
		{
		    "name": "create user other DB error",
			"email":"example@email.ru",
			"phone":"+79779814062",
		    "sessionId": "5212874b-1990-4360-b2f3-09e0728e1fb4",
			"role": "client"
		}`,
		},
		{"delete session DB error", http.StatusInternalServerError, def, `
		{
		    "name": "delete session DB error",
			"email":"example@email.ru",
			"phone":"+79779814062",
		    "sessionId": "9a2f99e0-f304-4905-b175-2a664707c2a0",
			"role": "client"
		}`,
		},
		{"valid case", http.StatusOK, def, `
		{
		    "name": "valid case",
			"email":"example@email.ru",
			"phone":"+79779814062",
		    "sessionId": "5212874b-1990-4360-b2f3-09e0728e1fb4",
			"role": "client"
		}`,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			h := &CreateHandler{
				uc: tt.fields.uc,
				tm: tt.fields.tm,
				v:  tt.fields.v,
				sv: tt.fields.sv,
				sd: tt.fields.sd,
			}

			req := httptest.NewRequest(http.MethodPost, "/v1/users", strings.NewReader(tt.payload))
			rw := httptest.NewRecorder()

			h.ServeHTTP(rw, req)
			assert.Equal(t, tt.expectedHTTPCode, rw.Code)
		})
	}
}
