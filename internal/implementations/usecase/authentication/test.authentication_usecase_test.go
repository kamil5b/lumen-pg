package authentication

import (
	"testing"

	testRunner "github.com/kamil5b/lumen-pg/internal/testrunners/usecase"
)

func TestAuthenticationUsecase(t *testing.T) {
	testRunner.AuthenticationUsecaseRunner(t, NewAuthenticationUseCaseImplementation)
}
