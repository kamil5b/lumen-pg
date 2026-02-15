package clock_repository

import (
	"github.com/kamil5b/lumen-pg/internal/interfaces/repository"
)

type ClockRepositoryImplementation struct {
	// clock implementation details will be added here
}

func NewClockRepository() repository.ClockRepository {
	return &ClockRepositoryImplementation{}
}
