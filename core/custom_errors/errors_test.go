package custom_errors

import (
	"testing"
)

func TestSentinelErrors(t *testing.T) {
	tests := []struct {
		name        string
		err         *NagareCoreError
		expectedCode string
	}{
		{"ErrAgentExecuteFailed", ErrAgentExecuteFailed, "AGENT_EXEC_FAILED"},
		{"ErrInvalidToolSchema", ErrInvalidToolSchema, "INVALID_TOOL_SCHEMA"},
		{"ErrSessionNotFound", ErrSessionNotFound, "SESSION_NOT_FOUND"},
		{"ErrPluginNotFound", ErrPluginNotFound, "PLUGIN_NOT_FOUND"},
		{"ErrUnknown", ErrUnknown, "UNKNOWN_ERROR"},
		{"ErrMissingDefaultProvider", ErrMissingDefaultProvider, "MISSING_DEFAULT_PROVIDER"},
		{"ErrGetSessionFailed", ErrGetSessionFailed, "GET_SESSION_FAILED"},
		{"ErrCreateCronJobFailed", ErrCreateCronJobFailed, "CREATE_CRONJOB_FAILED"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.err == nil {
				t.Fatalf("expected non-nil error for %s", tc.name)
			}
			if tc.err.Code != tc.expectedCode {
				t.Errorf("expected code %s, got %s", tc.expectedCode, tc.err.Code)
			}
			if tc.err.Details == "" {
				t.Errorf("expected non-empty details for %s", tc.name)
			}
		})
	}
}
