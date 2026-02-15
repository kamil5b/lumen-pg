package session_repository

import (
	"database/sql"
	"sync"

	"github.com/kamil5b/lumen-pg/internal/domain"
	"github.com/kamil5b/lumen-pg/internal/interfaces/repository"
)

type SessionRepositoryImplementation struct {
	mu       sync.RWMutex
	db       *sql.DB
	sessions map[string]*domain.Session
}

func NewSessionRepository(db *sql.DB) repository.SessionRepository {
	return &SessionRepositoryImplementation{
		db:       db,
		sessions: make(map[string]*domain.Session),
	}
}
