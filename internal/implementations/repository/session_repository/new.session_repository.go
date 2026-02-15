package session_repository

import (
	"sync"

	"github.com/kamil5b/lumen-pg/internal/domain"
	"github.com/kamil5b/lumen-pg/internal/interfaces/repository"
)

type SessionRepositoryImplementation struct {
	mu       sync.RWMutex
	sessions map[string]*domain.Session
}

func NewSessionRepository() repository.SessionRepository {
	return &SessionRepositoryImplementation{
		sessions: make(map[string]*domain.Session),
	}
}
