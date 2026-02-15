package transaction_repository

import (
	"database/sql"
	"sync"

	"github.com/kamil5b/lumen-pg/internal/domain"
	"github.com/kamil5b/lumen-pg/internal/interfaces/repository"
)

type TransactionRepositoryImplementation struct {
	mu           sync.RWMutex
	db           *sql.DB
	transactions map[string]*domain.TransactionState
	rowEdits     map[string]map[int]domain.RowEdit
	rowDeletes   map[string][]int
	rowInserts   map[string][]domain.RowInsert
}

func NewTransactionRepository(db *sql.DB) repository.TransactionRepository {
	return &TransactionRepositoryImplementation{
		db:           db,
		transactions: make(map[string]*domain.TransactionState),
		rowEdits:     make(map[string]map[int]domain.RowEdit),
		rowDeletes:   make(map[string][]int),
		rowInserts:   make(map[string][]domain.RowInsert),
	}
}
