package main_view_test

import (
	"testing"

	"github.com/kamil5b/lumen-pg/internal/implementations/handler/main_view"
	"github.com/kamil5b/lumen-pg/internal/interfaces/handler"
	"github.com/kamil5b/lumen-pg/internal/interfaces/usecase"
	handlerTestRunner "github.com/kamil5b/lumen-pg/internal/testrunners/handler"
)

func TestMainViewHandler(t *testing.T) {
	constructor := func(
		dataViewUC usecase.DataViewUseCase,
		authUC usecase.AuthenticationUseCase,
		rbacUC usecase.RBACUseCase,
	) handler.MainViewHandler {
		return main_view.NewMainViewHandlerImplementation(dataViewUC, authUC, rbacUC)
	}

	handlerTestRunner.MainViewHandlerRunner(t, constructor)
}
