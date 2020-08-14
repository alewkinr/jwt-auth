package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"example.com/back/auth/internal/session"
	"example.com/back/auth/internal/user"
)

type mockPhoneSessionPUTGetter struct {
	Session *session.Session
	Error   error
}

func (m *mockPhoneSessionPUTGetter) FindActiveBySessionID(ctx context.Context, sessionID string) (*session.Session, error) {
	return m.Session, m.Error
}

type mockPhoneSessionPUTSetter struct {
	Error error
}

func (m *mockPhoneSessionPUTSetter) Delete(ctx context.Context, sessionID string) error {
	return m.Error
}

func (m *mockPhoneSessionPUTSetter) SessionVerified(ctx context.Context, sessionID string) error {
	return m.Error
}

type mockPhoneSessionPUTValidator struct {
	Error error
}

func (m *mockPhoneSessionPUTValidator) Struct(interface{}) error {
	return m.Error
}

type mockUserManager struct {
	User  *user.User
	Error error
}

func (m *mockUserManager) GetUserByPhone(ctx context.Context, phone string) (*user.User, error) {
	return m.User, m.Error
}
func (m *mockUserManager) GetUserByID(ctx context.Context, id int) (*user.User, error) {
	return m.User, m.Error
}
func (m *mockUserManager) Create(ctx context.Context, name, email, password, phone string, role user.Role) (*user.User, error) {
	return m.User, m.Error
}

type mockPhoneSessionPUTTokenGetter struct {
	AccessToken  string
	RefreshToken string
	Error        error
}

func (m *mockPhoneSessionPUTTokenGetter) GetAccessToken(u *user.User) string {
	return m.AccessToken
}

func (m *mockPhoneSessionPUTTokenGetter) GetRefreshToken(u *user.User) string {
	return m.RefreshToken
}

var testsCasesPhoneSessionsPUTHandler = []struct {
	name         string
	sm           phoneSessionPUTGetter
	ss           phoneSessionPUTSetter
	v            phoneSessionPUTValidator
	uc           userManager
	tm           phoneSessionPUTTokenGetter
	reqSessionID string
	reqBody      *strings.Reader
	wantedBody   string
	wantedCode   int
}{
	{
		name: "Указан верный телефон и код. - указанного телефона нет в example.comн, это новый клиент.",
		sm: &mockPhoneSessionPUTGetter{
			Session: &session.Session{
				SessionID:        "7b12338d-7aea-4de9-8ba7-ae9377902fef",
				UsersPhone:       "+79045710785",
				VerificationCode: 123,
				ExpiresAt:        time.Now().UTC().Add(time.Duration(+30) * time.Hour),
				CreatedAt:        time.Now().UTC().Add(time.Duration(-40) * time.Second),
			},
			Error: nil},
		ss:           &mockPhoneSessionPUTSetter{Error: nil},
		v:            &mockPhoneSessionPUTValidator{Error: nil},
		uc:           &mockUserManager{Error: nil},
		reqSessionID: "7b12338d-7aea-4de9-8ba7-ae9377902fef",
		reqBody: strings.NewReader(`{
			 	"phone": "+79045710785",
				"code": 123
		}`),
		wantedBody: `{"sessionId":"7b12338d-7aea-4de9-8ba7-ae9377902fef"}`,
		wantedCode: 200,
	},
	{
		name: "Указан верный телефон и код. - Указанный телефонный номер есть в public.users, это вернувшийся клиент",
		sm: &mockPhoneSessionPUTGetter{
			Session: &session.Session{
				SessionID:        "7b12338d-7aea-4de9-8ba7-ae9377902fef",
				UsersPhone:       "+79045710785",
				VerificationCode: 123,
				ExpiresAt:        time.Now().UTC().Add(time.Duration(+30) * time.Hour),
				CreatedAt:        time.Now().UTC().Add(time.Duration(-40) * time.Second),
			},
			Error: nil},
		ss: &mockPhoneSessionPUTSetter{Error: nil},
		v:  &mockPhoneSessionPUTValidator{Error: nil},
		uc: &mockUserManager{
			User: &user.User{
				ID:        1,
				Name:      "Ivanov",
				Email:     "ivanov@gmail.com",
				Password:  "secret",
				Phone:     "+79045710785",
				Role:      1,
				Status:    1,
				CreatedOn: time.Date(2020, 5, 1, 12, 0, 0, 0, time.UTC),
				LastLogin: &time.Time{},
			},
			Error: nil,
		},
		tm: &mockPhoneSessionPUTTokenGetter{
			AccessToken:  "atoken",
			RefreshToken: "rtoken",
			Error:        nil,
		},
		reqSessionID: "7b12338d-7aea-4de9-8ba7-ae9377902fef",
		reqBody: strings.NewReader(`{
			 	"phone": "+79045710785",
				"code": 123
		}`),
		wantedBody: `{"userId":1,"sessionId":"7b12338d-7aea-4de9-8ba7-ae9377902fef","accessToken":"atoken","refreshToken":"rtoken"}`,
		wantedCode: 200,
	},
	{
		name: "Нет активной сессии в БД",
		sm: &mockPhoneSessionPUTGetter{
			Session: nil,
			Error:   session.ErrRecordNotFound},
		ss:           &mockPhoneSessionPUTSetter{Error: nil},
		v:            &mockPhoneSessionPUTValidator{Error: nil},
		uc:           &mockUserManager{},
		tm:           &mockPhoneSessionPUTTokenGetter{},
		reqSessionID: "7b12338d-7aea-4de9-8ba7-ae9377902fef",
		reqBody: strings.NewReader(`{
			 	"phone": "+79045710785",
				"code": 123
		}`),
		wantedBody: `{"error": "active session not found"}`,
		wantedCode: 204,
	},
	{
		name: "Рейт-лимит превышает лимит",
		sm: &mockPhoneSessionPUTGetter{
			Session: &session.Session{
				SessionID:        "7b12338d-7aea-4de9-8ba7-ae9377902fef",
				UsersPhone:       "+79045710785",
				VerificationCode: 123,
				ExpiresAt:        time.Now().UTC(),
				CreatedAt:        time.Now().UTC(),
			},
			Error: nil},
		ss:           &mockPhoneSessionPUTSetter{Error: nil},
		v:            &mockPhoneSessionPUTValidator{Error: nil},
		uc:           &mockUserManager{},
		tm:           &mockPhoneSessionPUTTokenGetter{},
		reqSessionID: "7b12338d-7aea-4de9-8ba7-ae9377902fef",
		reqBody: strings.NewReader(`{
			 	"phone": "+79045710785",
				"code": 123
		}`),
		wantedBody: `{"error": "too many requests, try again later"}`,
		wantedCode: 429,
	},
	{
		name: "Телефон из запроса не совпадает с телефоном в сесии",
		sm: &mockPhoneSessionPUTGetter{
			Session: &session.Session{
				SessionID:        "7b12338d-7aea-4de9-8ba7-ae9377902fef",
				UsersPhone:       "+79045710785",
				VerificationCode: 123,
				ExpiresAt:        time.Now().UTC().Add(time.Duration(+30) * time.Hour),
				CreatedAt:        time.Now().UTC().Add(time.Duration(-40) * time.Second),
			},
			Error: nil},
		ss:           &mockPhoneSessionPUTSetter{Error: nil},
		v:            &mockPhoneSessionPUTValidator{Error: nil},
		uc:           &mockUserManager{},
		tm:           &mockPhoneSessionPUTTokenGetter{},
		reqSessionID: "7b12338d-7aea-4de9-8ba7-ae9377902fef",
		reqBody: strings.NewReader(`{
			 	"phone": "+79045710786",
				"code": 123
		}`),
		wantedBody: `{"error": "bad phone number"}`,
		wantedCode: 204,
	},
	{
		name: "Код подтверждения из запроса не совпадает с кодом в сесии",
		sm: &mockPhoneSessionPUTGetter{
			Session: &session.Session{
				SessionID:        "7b12338d-7aea-4de9-8ba7-ae9377902fef",
				UsersPhone:       "+79045710785",
				VerificationCode: 123,
				ExpiresAt:        time.Now().UTC().Add(time.Duration(+30) * time.Hour),
				CreatedAt:        time.Now().UTC().Add(time.Duration(-40) * time.Second),
			},
			Error: nil},
		ss:           &mockPhoneSessionPUTSetter{Error: nil},
		v:            &mockPhoneSessionPUTValidator{Error: nil},
		uc:           &mockUserManager{},
		tm:           &mockPhoneSessionPUTTokenGetter{},
		reqSessionID: "7b12338d-7aea-4de9-8ba7-ae9377902fef",
		reqBody: strings.NewReader(`{
			 	"phone": "+79045710785",
				"code": 124
		}`),
		wantedBody: `{"error": "bad code number"}`,
		wantedCode: 204,
	},
	{
		name:         "не валидный sessionId",
		sm:           &mockPhoneSessionPUTGetter{Error: nil},
		ss:           &mockPhoneSessionPUTSetter{Error: nil},
		v:            &mockPhoneSessionPUTValidator{Error: nil},
		uc:           &mockUserManager{},
		tm:           &mockPhoneSessionPUTTokenGetter{},
		reqSessionID: "7b12338d",
		reqBody: strings.NewReader(`{
			 	"phone": "+79045710785",
				"code": 124
		}`),
		wantedBody: `{"error": "bad request"}`,
		wantedCode: 400,
	},
}

func TestPhoneSessionsPUTHandler_ServeHTTP(t *testing.T) {
	for _, tt := range testsCasesPhoneSessionsPUTHandler {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			h := NewPhoneSessionsPUTHandler(tt.sm, tt.ss, tt.uc, tt.tm, tt.v)
			rr := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodPut,
				fmt.Sprintf("/v1/phone_sessions/%v/sms_code", tt.reqSessionID), tt.reqBody)
			if err != nil {
				assert.FailNow(t, "failed to create test request")
			}
			req = mux.SetURLVars(req, map[string]string{"sessionId": tt.reqSessionID})
			h.ServeHTTP(rr, req)
			assert.Equal(t, tt.wantedCode, rr.Code)
			assert.Equal(t, tt.wantedBody, deleteBodyNewLine(rr.Body.String()))
		})
	}
}

func deleteBodyNewLine(s string) string {
	return strings.ReplaceAll(s, "\n", "")
}

func TestPhoneSessionsPUTHandler_isPhoneValid(t *testing.T) {
	type args struct {
		wantedPhone   string
		expectedPhone string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "phones match",
			args: args{
				wantedPhone:   "+79045710785",
				expectedPhone: "+79045710785",
			},
			want: true,
		},
		{
			name: "phones not match",
			args: args{
				wantedPhone:   "+79045710785",
				expectedPhone: "+79055710785",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			h := &PhoneSessionsPUTHandler{}
			if got := h.isPhoneValid(tt.args.wantedPhone, tt.args.expectedPhone); got != tt.want {
				t.Errorf("isPhoneValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPhoneSessionsPUTHandler_isCodeValid(t *testing.T) {
	type args struct {
		wantedCode   int64
		expectedCode int64
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "code match",
			args: args{
				wantedCode:   123,
				expectedCode: 123,
			},
			want: true,
		},
		{
			name: "code not match",
			args: args{
				wantedCode:   123,
				expectedCode: 23,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			h := &PhoneSessionsPUTHandler{}
			if got := h.isCodeValid(tt.args.wantedCode, tt.args.expectedCode); got != tt.want {
				t.Errorf("isCodeValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
