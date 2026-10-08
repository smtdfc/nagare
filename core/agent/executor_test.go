package agent

import (
	"errors"
	"testing"

	"github.com/smtdfc/nagare/core/custom_errors"
	"github.com/smtdfc/nagare/pkgs/messages"
)

func TestExecutor_HandleError(t *testing.T) {
	// Verify HandleError mapping HTTP error codes to custom_errors
	ex := &Executor{}

	tests := []struct {
		code        string
		expectedErr *custom_errors.NagareCoreError
	}{
		{"429", custom_errors.ErrModelQuotaExceed},
		{"402", custom_errors.ErrModelPaymentRequired},
		{"500", custom_errors.ErrLLMProviderAdapter},
		{"503", custom_errors.ErrLLMProviderAdapter},
		{"504", custom_errors.ErrLLMProviderAdapter},
	}

	for _, tc := range tests {
		t.Run(tc.code, func(t *testing.T) {
			errMsg := messages.NewResponseFailedMessage(tc.code, "error details")
			err := ex.HandleError(errMsg)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected %v, got %v", tc.expectedErr, err)
			}
		})
	}
}
