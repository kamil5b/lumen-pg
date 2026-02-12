package interfaces

import "github.com/go-chi/chi/v5"

type UserHandler interface {
	RegisterRoutes(r chi.Router)
}
