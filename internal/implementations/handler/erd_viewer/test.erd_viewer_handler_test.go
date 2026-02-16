package erd_viewer_test

import (
	"testing"

	"github.com/kamil5b/lumen-pg/internal/implementations/handler/erd_viewer"
	"github.com/kamil5b/lumen-pg/internal/interfaces/handler"
	"github.com/kamil5b/lumen-pg/internal/interfaces/usecase"
	handlerTestRunner "github.com/kamil5b/lumen-pg/internal/testrunners/handler"
)

func TestERDViewerHandler(t *testing.T) {
	constructor := func(
		erdUC usecase.ERDUseCase,
		authUC usecase.AuthenticationUseCase,
	) handler.ERDViewerHandler {
		return erd_viewer.NewERDViewerHandlerImplementation(erdUC, authUC)
	}

	handlerTestRunner.ERDViewerHandlerRunner(t, constructor)
}
