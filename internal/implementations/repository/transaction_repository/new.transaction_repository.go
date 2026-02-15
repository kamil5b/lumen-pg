package transaction_repository

import (
	"sync"

	"github.com/kamil5b/lumen-pg/internal/domain"
	"github.com/kamil5b/lumen-pg/internal/interfaces/repository"
)

type TransactionRepositoryImplementation struct {
	mu           sync.RWMutex
	transactions map[string]*domain.TransactionState
	rowEdits     map[string]map[int]domain.RowEdit
	rowDeletes   map[string][]int
	rowInserts   map[string][]domain.RowInsert
}

func NewTransactionRepository() repository.TransactionRepository {
	return &TransactionRepositoryImplementation{
		transactions: make(map[string]*domain.TransactionState),
		rowEdits:     make(map[string]map[int]domain.RowEdit),
		rowDeletes:   make(map[string][]int),
		rowInserts:   make(map[string][]domain.RowInsert),
	}
}
