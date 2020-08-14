package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	val10 "github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/assert"
	mocks_phone_sessions_post "example.com/back/auth/internal/mocks"
	"example.com/back/auth/internal/session"
)

// nolint:
func TestPhoneSessionsPOSTHandler(t *testing.T) {
	type fields struct {
		sm phoneSessionPOSTGetter
		ss phoneSessionPOSTSaver
		v  phoneSessionPOSTValidator
		n  notifier
	}

	mock := gomock.NewController(t)
	defer mock.Finish()

	sg := mocks_phone_sessions_post.NewMockphoneSessionPOSTGetter(mock)
	sgСall := sg.EXPECT().FindLastActiveByUserPhone(gomock.Any(), gomock.Any()).AnyTimes()
	sgСall.DoAndReturn(func(ctx context.Context, usersPhone string) (*session.Session, error) {
		switch usersPhone {
		// limits fired
		case "+79779814063":
			return &session.Session{
				SessionID:        "session id",
				UsersPhone:       "+79779814062",
				VerificationCode: 2345,
				ExpiresAt:        time.Now().Add(time.Second * 60),
				CreatedAt:        time.Now().Add(time.Second * 30),
			}, nil

		// limits not fired
		case "+78888888888":
			return nil, session.ErrRecordNotFound
		}
		return nil, session.ErrRecordNotFound
	})

	ss := mocks_phone_sessions_post.NewMockphoneSessionPOSTSaver(mock)
	ssCall := ss.EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	ssCall.DoAndReturn(func(ctx context.Context, usersPhone string, code int64) (*session.Session, error) {
		switch usersPhone {
		// valid case
		default:
			return &session.Session{
				SessionID:        "session id",
				UsersPhone:       "+79779814062",
				VerificationCode: 2345,
				ExpiresAt:        time.Now().Add(time.Second * 60),
				CreatedAt:        time.Now().Add(time.Second * 30),
			}, nil
		// err duplicate
		case "+78888888888":
			return nil, session.ErrDuplicateKey
		case "+76666666666":
			return nil, errors.New("testing storage save error")
		}
	})
	sdCall := ss.EXPECT().DeleteLastByPhone(gomock.Any(), gomock.Any()).AnyTimes()
	sdCall.DoAndReturn(func(ctx context.Context, usersPhone string) error {
		switch usersPhone {
		case "+73333333333":
			return errors.New("failed to delete session from storage")
		default:
			return nil
		}
	})
	v := val10.New()

	n := mocks_phone_sessions_post.NewMocknotifier(mock)
	nCall := n.EXPECT().POSTV1UsersSendMessageSMS(gomock.Any(), gomock.Any()).AnyTimes()
	nCall.DoAndReturn(func(userPhone, message string) error {
		switch userPhone {
		// valid case
		default:
			return nil
		// err duplicate
		case "+76666666666":
			return errors.New("testing send sms error")
		}
	})

	def := fields{
		sm: sg,
		ss: ss,
		v:  v,
		n:  n,
	}

	tests := []struct {
		name             string
		expectedHTTPCode int
		fields           fields
		payload          string
	}{
		{"not JSON body", http.StatusBadRequest, def, "it's not a JSON"},
		{"missing required field in body: phone", http.StatusBadRequest, def, `
			{
			  "testing": "testing field",
			  "countryCode": "RU"
			}`,
		},
		{"missing required field in body: country code", http.StatusBadRequest, def, `
			{
			  "phone": "89779814062",

			}`,
		},
		{"not valid phone number", http.StatusBadRequest, def, `
			{
			  "phone": "89779814062",
			  "countryCode": "RU"
			}`,
		},
		{"rate limits restriction", http.StatusTooManyRequests, def, `
			{
			  "phone": "+79779814063",
			  "countryCode": "RU"
			}`,
		},
		{"duplicated session id in DB", http.StatusInternalServerError, def, `
			{
			  "phone": "+78888888888",
			  "countryCode": "RU"
			}`,
		},
		{"save failed to delete session", http.StatusInternalServerError, def, `
			{
			  "phone": "+73333333333",
			  "countryCode": "RU"
			}`,
		},
		{"save session to storage error", http.StatusInternalServerError, def, `
			{
			  "phone": "+76666666666",
			  "countryCode": "RU"
			}`,
		},
		{"failed to send sms to phone", http.StatusInternalServerError, def, `
			{
			  "phone": "+76666666666",
			  "countryCode": "RU"
			}`,
		},
		{"valid case", http.StatusOK, def, `
			{
			  "phone": "+79779814062",
			  "countryCode": "RU"
			}`,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			h := &PhoneSessionsPOSTHandler{
				sm: tt.fields.sm,
				ss: tt.fields.ss,
				v:  tt.fields.v,
				n:  tt.fields.n,
			}

			req := httptest.NewRequest(http.MethodPost, "/v1/phone_sessions", strings.NewReader(tt.payload))
			rw := httptest.NewRecorder()

			h.ServeHTTP(rw, req)
			assert.Equal(t, tt.expectedHTTPCode, rw.Code)
		})
	}
}

func countDigits(number int64) int {
	count := 0
	for number != 0 {
		number /= 10
		count++
	}
	return count
}
func Test_generateVerificationCode(t *testing.T) {
	testName := "make sure its generates code with len 4"
	testRuns := 100000 // how many times we will run this test
	wantDigits := 4
	t.Run(testName, func(t *testing.T) {
		for i := 1; i <= testRuns; i++ {
			got := generateVerificationCode()
			gotDigits := countDigits(got)
			if gotDigits != wantDigits {
				t.Errorf("generateVerificationCode() digits = %v, want %v", got, wantDigits)
			}
		}
	})

	testName = "make sure its different codes every time"
	// this might be unsafe, coz it's really can return the same value
	testRuns = 100 // how many times we will run this test
	t.Run(testName, func(t *testing.T) {
		for i := 1; i <= testRuns; i++ {
			gotFirst := generateVerificationCode()
			gotSecond := generateVerificationCode()
			isDifferent := gotFirst != gotSecond
			if !isDifferent {
				t.Errorf("generateVerificationCode() returned same values gotFirst = %v, gotSecond = %v", gotFirst, gotSecond)
			}
		}
	})
}
