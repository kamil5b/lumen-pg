package erd_viewer

import (
	"github.com/kamil5b/lumen-pg/internal/interfaces/handler"
	"github.com/kamil5b/lumen-pg/internal/interfaces/usecase"
)

type ERDViewerHandlerImplementation struct {
	erdUC  usecase.ERDUseCase
	authUC usecase.AuthenticationUseCase
}

func NewERDViewerHandlerImplementation(
	erdUC usecase.ERDUseCase,
	authUC usecase.AuthenticationUseCase,
) handler.ERDViewerHandler {
	return &ERDViewerHandlerImplementation{
		erdUC:  erdUC,
		authUC: authUC,
	}
}
