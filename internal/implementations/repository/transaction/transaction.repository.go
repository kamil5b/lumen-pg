package transaction

// TransactionRepository is a noop implementation of the transaction repository
type TransactionRepository struct{}

// NewTransactionRepository creates a new transaction repository
func NewTransactionRepository() *TransactionRepository {
	return &TransactionRepository{}
}
