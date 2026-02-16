package login

import (
	"github.com/kamil5b/lumen-pg/internal/interfaces/handler"
	"github.com/kamil5b/lumen-pg/internal/interfaces/usecase"
)

type LoginHandlerImplementation struct {
	authUC  usecase.AuthenticationUseCase
	setupUC usecase.SetupUseCase
	rbacUC  usecase.RBACUseCase
}

func NewLoginHandlerImplementation(
	authUC usecase.AuthenticationUseCase,
	setupUC usecase.SetupUseCase,
	rbacUC usecase.RBACUseCase,
) handler.LoginHandler {
	return &LoginHandlerImplementation{
		authUC:  authUC,
		setupUC: setupUC,
		rbacUC:  rbacUC,
	}
}
