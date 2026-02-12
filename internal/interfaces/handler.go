package interfaces

import (
	"github.com/go-chi/chi/v5"
)

// AuthHandler handles authentication HTTP endpoints
type AuthHandler interface {
	RegisterRoutes(r chi.Router)
}

// DataExplorerHandler handles data explorer HTTP endpoints
type DataExplorerHandler interface {
	RegisterRoutes(r chi.Router)
}

// QueryEditorHandler handles query editor HTTP endpoints
type QueryEditorHandler interface {
	RegisterRoutes(r chi.Router)
}

// ERDHandler handles ERD viewer HTTP endpoints
type ERDHandler interface {
	RegisterRoutes(r chi.Router)
}
