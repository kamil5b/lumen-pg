package security

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *SecurityUseCaseImplementation) GenerateCookieSignature(ctx context.Context, cookieData *domain.CookieData) (string, error) {
	return "", errors.New("not implemented")
}
