package security

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *SecurityUseCaseImplementation) ValidateCookieIntegrity(ctx context.Context, cookieData *domain.CookieData, signature string) (bool, error) {
	// Convert CookieData to JSON string
	jsonBytes, err := json.Marshal(cookieData)
	if err != nil {
		return false, fmt.Errorf("failed to marshal cookie data: %w", err)
	}

	// Generate the expected signature for the cookie data
	expectedSignature, err := u.encryptionRepo.GenerateSignature(ctx, string(jsonBytes))
	if err != nil {
		return false, fmt.Errorf("failed to generate expected signature: %w", err)
	}

	// Compare the provided signature with the expected signature
	return expectedSignature == signature, nil
}
