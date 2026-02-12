package interfaces

import "github.com/go-chi/chi/v5"

// AuthHandler handles authentication HTTP requests
type AuthHandler interface {
	RegisterRoutes(r chi.Router)
}

// MainViewHandler handles main view HTTP requests
type MainViewHandler interface {
	RegisterRoutes(r chi.Router)
}

// QueryEditorHandler handles query editor HTTP requests
type QueryEditorHandler interface {
	RegisterRoutes(r chi.Router)
}

// ERDViewerHandler handles ERD viewer HTTP requests
type ERDViewerHandler interface {
	RegisterRoutes(r chi.Router)
}

// TransactionHandler handles transaction HTTP requests
type TransactionHandler interface {
	RegisterRoutes(r chi.Router)
}
