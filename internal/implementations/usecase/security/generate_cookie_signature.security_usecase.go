package security

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *SecurityUseCaseImplementation) GenerateCookieSignature(ctx context.Context, cookieData *domain.CookieData) (string, error) {
	// Convert CookieData to JSON string
	jsonBytes, err := json.Marshal(cookieData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal cookie data: %w", err)
	}

	// Generate signature for the JSON string using the encryption repository
	signature, err := u.encryptionRepo.GenerateSignature(ctx, string(jsonBytes))
	if err != nil {
		return "", fmt.Errorf("failed to generate cookie signature: %w", err)
	}

	return signature, nil
}
