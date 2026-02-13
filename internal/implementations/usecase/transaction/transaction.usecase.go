package transaction

// TransactionUseCase is a noop implementation of the transaction usecase
type TransactionUseCase struct{}

// NewTransactionUseCase creates a new transaction usecase
func NewTransactionUseCase() *TransactionUseCase {
	return &TransactionUseCase{}
}
