package restapireceiver

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
		errMsg  string
	}{
		{
			name:    "ValidConfigWithAuthToken",
			config:  Config{Endpoint: "http://example.com", AuthToken: "someAuthToken"},
			wantErr: false,
		},
		{
			name:    "ValidConfigWithUsernamePassword",
			config:  Config{Endpoint: "http://example.com", Username: "user", Password: "pass"},
			wantErr: false,
		},
		{
			name:    "MissingEndpoint",
			config:  Config{AuthToken: "someAuthToken"},
			wantErr: true,
			errMsg:  "Config validation failed: 'endpoint' is required",
		},
		{
			name:    "MissingAuthTokenAndUsernamePassword",
			config:  Config{Endpoint: "http://example.com"},
			wantErr: true,
			errMsg:  "Config validation failed: either of 'auth_token' or 'username'+'password' are required",
		},
		{
			name:    "MissingUsernameWithPassword",
			config:  Config{Endpoint: "http://example.com", Password: "pass"},
			wantErr: true,
			errMsg:  "Config validation failed: either of 'auth_token' or 'username'+'password' are required",
		},
		{
			name:    "MissingPasswordWithUsername",
			config:  Config{Endpoint: "http://example.com", Username: "user"},
			wantErr: true,
			errMsg:  "Config validation failed: either of 'auth_token' or 'username'+'password' are required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
