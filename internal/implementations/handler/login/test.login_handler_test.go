package login_test

import (
	"testing"

	"github.com/kamil5b/lumen-pg/internal/implementations/handler/login"
	"github.com/kamil5b/lumen-pg/internal/interfaces/handler"
	"github.com/kamil5b/lumen-pg/internal/interfaces/usecase"
	handlerTestRunner "github.com/kamil5b/lumen-pg/internal/testrunners/handler"
)

func TestLoginHandler(t *testing.T) {
	constructor := func(
		authUC usecase.AuthenticationUseCase,
		setupUC usecase.SetupUseCase,
		rbacUC usecase.RBACUseCase,
	) handler.LoginHandler {
		return login.NewLoginHandlerImplementation(authUC, setupUC, rbacUC)
	}

	handlerTestRunner.AuthenticationHandlerRunner(t, constructor)
}
