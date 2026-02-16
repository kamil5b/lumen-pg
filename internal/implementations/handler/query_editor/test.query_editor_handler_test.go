package query_editor_test

import (
	"testing"

	"github.com/kamil5b/lumen-pg/internal/implementations/handler/query_editor"
	"github.com/kamil5b/lumen-pg/internal/interfaces/handler"
	"github.com/kamil5b/lumen-pg/internal/interfaces/usecase"
	handlerTestRunner "github.com/kamil5b/lumen-pg/internal/testrunners/handler"
)

func TestQueryEditorHandler(t *testing.T) {
	constructor := func(
		queryUC usecase.QueryUseCase,
		authUC usecase.AuthenticationUseCase,
	) handler.QueryEditorHandler {
		return query_editor.NewQueryEditorHandlerImplementation(queryUC, authUC)
	}

	handlerTestRunner.QueryEditorHandlerRunner(t, constructor)
}
