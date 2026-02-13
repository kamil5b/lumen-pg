package query

// QueryUseCase is a noop implementation of the query usecase
type QueryUseCase struct{}

// NewQueryUseCase creates a new query usecase
func NewQueryUseCase() *QueryUseCase {
	return &QueryUseCase{}
}
