package dataview

import (
	"testing"

	testRunner "github.com/kamil5b/lumen-pg/internal/testrunners/usecase"
)

func TestDataViewUsecase(t *testing.T) {
	testRunner.DataViewUsecaseRunner(t, NewDataViewUseCaseImplementation)
}
