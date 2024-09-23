package auth

import (
	"net/http"
	"testing"
)

func TestGetToken(t *testing.T) {
	tests := []struct {
		name        string
		tokenType   string
		headers     http.Header
		expected    string
		expectedErr error
	}{
		{
			name:        "Valid Bearer token",
			tokenType:   AccessTokenType,
			headers:     http.Header{"Authorization": {"Bearer valid_token"}},
			expected:    "valid_token",
			expectedErr: nil,
		},
		{
			name:        "Valid ApiKey token",
			tokenType:   APIKeyTokenType,
			headers:     http.Header{"Authorization": {"ApiKey valid_token"}},
			expected:    "valid_token",
			expectedErr: nil,
		},
		{
			name:        "Missing Authorization header",
			tokenType:   AccessTokenType,
			headers:     http.Header{},
			expected:    "",
			expectedErr: ErrMissingAuthHeader,
		},
		{
			name:        "Invalid Authorization header format",
			tokenType:   AccessTokenType,
			headers:     http.Header{"Authorization": {"InvalidHeader"}},
			expected:    "",
			expectedErr: ErrInvalidAuthHeader,
		},
		{
			name:        "Incorrect token type",
			tokenType:   AccessTokenType,
			headers:     http.Header{"Authorization": {"ApiKey valid_token"}},
			expected:    "",
			expectedErr: ErrInvalidAuthHeader,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GetToken(tt.tokenType, tt.headers)
			if token != tt.expected || err != tt.expectedErr {
				t.Errorf("GetToken() = %v, %v; want %v, %v", token, err, tt.expected, tt.expectedErr)
			}
		})
	}
}
