package main_view

import (
	"github.com/kamil5b/lumen-pg/internal/interfaces/handler"
	"github.com/kamil5b/lumen-pg/internal/interfaces/usecase"
)

type MainViewHandlerImplementation struct {
	dataViewUC usecase.DataViewUseCase
	authUC     usecase.AuthenticationUseCase
	rbacUC     usecase.RBACUseCase
}

func NewMainViewHandlerImplementation(
	dataViewUC usecase.DataViewUseCase,
	authUC usecase.AuthenticationUseCase,
	rbacUC usecase.RBACUseCase,
) handler.MainViewHandler {
	return &MainViewHandlerImplementation{
		dataViewUC: dataViewUC,
		authUC:     authUC,
		rbacUC:     rbacUC,
	}
}
