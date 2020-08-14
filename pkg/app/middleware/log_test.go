package middleware

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func Test_pathToMetric(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "one test",
			args: args{
				"/v1/phone_sessions/7b12338d-7aea-4de9-8ba7-ae9377902fef/sms_code",
			},
			want: "v1_phone_sessions_7b12338d-7aea-4de9-8ba7-ae9377902fef_sms_code",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, pathToMetric(tt.args.path), tt.want)
		})
	}
}
