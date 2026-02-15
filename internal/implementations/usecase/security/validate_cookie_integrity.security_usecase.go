package security

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *SecurityUseCaseImplementation) ValidateCookieIntegrity(ctx context.Context, cookieData *domain.CookieData, signature string) (bool, error) {
	return false, errors.New("not implemented")
}
