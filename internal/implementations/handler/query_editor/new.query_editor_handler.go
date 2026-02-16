package query_editor

import (
	"github.com/kamil5b/lumen-pg/internal/interfaces/handler"
	"github.com/kamil5b/lumen-pg/internal/interfaces/usecase"
)

type QueryEditorHandlerImplementation struct {
	queryUC usecase.QueryUseCase
	authUC  usecase.AuthenticationUseCase
}

func NewQueryEditorHandlerImplementation(
	queryUC usecase.QueryUseCase,
	authUC usecase.AuthenticationUseCase,
) handler.QueryEditorHandler {
	return &QueryEditorHandlerImplementation{
		queryUC: queryUC,
		authUC:  authUC,
	}
}
