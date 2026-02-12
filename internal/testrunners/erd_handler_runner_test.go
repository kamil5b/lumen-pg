package testrunners

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/kamil5b/lumen-pg/internal/interfaces/usecase"
	"github.com/kamil5b/lumen-pg/internal/implementations/mocks"
)

// ERDHandlerConstructor creates an ERD handler with its dependencies
type ERDHandlerConstructor func(metadataUseCase usecase.MetadataUseCase) usecase.ERDHandler

// ERDHandlerRunner runs test specs for ERD handler (Story 3 E2E)
func ERDHandlerRunner(t *testing.T, constructor ERDHandlerConstructor) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMetadataUseCase := mocks.NewMockMetadataUseCase(ctrl)
	handler := constructor(mockMetadataUseCase)

	r := chi.NewRouter()
	handler.RegisterRoutes(r)

	t.Run("E2E-S3-01: ERD Viewer Page Access", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/erd", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Header().Get("Content-Type"), "text/html")
	})

	t.Run("E2E-S3-02: ERD Zoom Controls", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/erd/diagram", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		// Should return diagram data that includes zoom capability
	})

	t.Run("E2E-S3-03: ERD Pan", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/erd/diagram", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		// Should return diagram data that supports panning
	})

	t.Run("E2E-S3-04: Table Click in ERD", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/erd/table/users", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		// Should return table details for display in side panel
	})
}
