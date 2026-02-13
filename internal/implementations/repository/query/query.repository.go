package query

// QueryRepository is a noop implementation of the query repository
type QueryRepository struct{}

// NewQueryRepository creates a new query repository
func NewQueryRepository() *QueryRepository {
	return &QueryRepository{}
}
