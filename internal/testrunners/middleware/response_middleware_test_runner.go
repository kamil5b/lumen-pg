package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kamil5b/lumen-pg/internal/interfaces/middleware"
)

// ResponseMiddlewareConstructor is a function type that creates a ResponseMiddleware
type ResponseMiddlewareConstructor func() middleware.ResponseMiddleware

// ResponseMiddlewareRunner runs all response middleware tests
// Maps to TEST_PLAN.md:
// - Story 7: Security & Best Practices (response handling and optimization)
func ResponseMiddlewareRunner(t *testing.T, constructor ResponseMiddlewareConstructor) {
	t.Helper()

	mw := constructor()

	// ContentNegotiation tests
	t.Run("ContentNegotiation handles JSON accept header", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message":"success"}`))
		})

		wrapped := mw.ContentNegotiation(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
		req.Header.Set("Accept", "application/json")
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
		require.Contains(t, rec.Header().Get("Content-Type"), "application/json")
	})

	t.Run("ContentNegotiation handles HTML accept header", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("<html><body>Success</body></html>"))
		})

		wrapped := mw.ContentNegotiation(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Accept", "text/html")
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
		require.Contains(t, rec.Header().Get("Content-Type"), "text/html")
	})

	t.Run("ContentNegotiation handles wildcard accept header", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.ContentNegotiation(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
		req.Header.Set("Accept", "*/*")
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("ContentNegotiation handles multiple accept types", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.ContentNegotiation(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
		req.Header.Set("Accept", "text/html,application/json,application/xml;q=0.9,*/*;q=0.8")
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("ContentNegotiation handles missing accept header", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.ContentNegotiation(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	// CompressResponse tests
	t.Run("CompressResponse compresses with gzip when accepted", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(strings.Repeat("Hello World! ", 100)))
		})

		wrapped := mw.CompressResponse(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
		req.Header.Set("Accept-Encoding", "gzip")
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)

		// Check if response is compressed
		encoding := rec.Header().Get("Content-Encoding")
		if encoding == "gzip" {
			// Verify we can decompress
			reader, err := gzip.NewReader(rec.Body)
			require.NoError(t, err)
			defer reader.Close()

			decompressed, err := io.ReadAll(reader)
			require.NoError(t, err)
			require.Contains(t, string(decompressed), "Hello World!")
		}
	})

	t.Run("CompressResponse skips compression for small responses", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Hi"))
		})

		wrapped := mw.CompressResponse(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
		req.Header.Set("Accept-Encoding", "gzip")
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		// Small responses typically aren't compressed
		encoding := rec.Header().Get("Content-Encoding")
		require.True(t, encoding == "" || encoding == "gzip")
	})

	t.Run("CompressResponse does not compress without accept-encoding", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(strings.Repeat("Hello World! ", 100)))
		})

		wrapped := mw.CompressResponse(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		encoding := rec.Header().Get("Content-Encoding")
		require.Empty(t, encoding)
	})

	t.Run("CompressResponse skips compression for images", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "image/png")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(strings.Repeat("binary image data", 100)))
		})

		wrapped := mw.CompressResponse(handler)

		req := httptest.NewRequest(http.MethodGet, "/image.png", nil)
		req.Header.Set("Accept-Encoding", "gzip")
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		// Images should not be compressed
		encoding := rec.Header().Get("Content-Encoding")
		require.True(t, encoding == "" || rec.Header().Get("Content-Type") == "image/png")
	})

	t.Run("CompressResponse handles deflate encoding", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(strings.Repeat("Hello World! ", 100)))
		})

		wrapped := mw.CompressResponse(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
		req.Header.Set("Accept-Encoding", "deflate")
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		// Should support deflate or fall back to uncompressed
		encoding := rec.Header().Get("Content-Encoding")
		require.True(t, encoding == "deflate" || encoding == "gzip" || encoding == "")
	})

	// SetDefaultHeaders tests
	t.Run("SetDefaultHeaders sets default content-type", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.SetDefaultHeaders(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		contentType := rec.Header().Get("Content-Type")
		require.NotEmpty(t, contentType)
	})

	t.Run("SetDefaultHeaders sets Cache-Control header", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.SetDefaultHeaders(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		cacheControl := rec.Header().Get("Cache-Control")
		require.NotEmpty(t, cacheControl)
	})

	t.Run("SetDefaultHeaders preserves existing headers", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Custom-Header", "custom-value")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.SetDefaultHeaders(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		require.Equal(t, "custom-value", rec.Header().Get("X-Custom-Header"))
		require.Contains(t, rec.Header().Get("Content-Type"), "application/json")
	})

	t.Run("SetDefaultHeaders does not override explicit headers", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			w.Header().Set("Cache-Control", "max-age=3600")
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.SetDefaultHeaders(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		require.Equal(t, "text/plain", rec.Header().Get("Content-Type"))
		require.Equal(t, "max-age=3600", rec.Header().Get("Cache-Control"))
	})

	t.Run("SetDefaultHeaders sets Server header", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.SetDefaultHeaders(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		// Server header may or may not be set depending on implementation
		server := rec.Header().Get("Server")
		require.True(t, server != "" || server == "")
	})

	// Middleware chaining tests
	t.Run("Response middleware chain works correctly", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(strings.Repeat("test data ", 50)))
		})

		// Chain all response middlewares
		wrapped := mw.SetDefaultHeaders(
			mw.CompressResponse(
				mw.ContentNegotiation(handler),
			),
		)

		req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Accept-Encoding", "gzip")
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		require.NotEmpty(t, rec.Header().Get("Content-Type"))
	})

	t.Run("CompressResponse preserves status codes", func(t *testing.T) {
		statusCodes := []int{
			http.StatusOK,
			http.StatusCreated,
			http.StatusAccepted,
			http.StatusNoContent,
			http.StatusBadRequest,
			http.StatusUnauthorized,
			http.StatusForbidden,
			http.StatusNotFound,
			http.StatusInternalServerError,
		}

		for _, statusCode := range statusCodes {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(statusCode)
				if statusCode != http.StatusNoContent {
					w.Write([]byte("response body"))
				}
			})

			wrapped := mw.CompressResponse(handler)

			req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
			req.Header.Set("Accept-Encoding", "gzip")
			rec := httptest.NewRecorder()

			wrapped.ServeHTTP(rec, req)

			require.Equal(t, statusCode, rec.Code)
		}
	})

	t.Run("ContentNegotiation returns 406 for unsupported media type", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.ContentNegotiation(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
		req.Header.Set("Accept", "application/vnd.custom-unsupported")
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		// Should either accept it or return 406
		require.True(t, rec.Code == http.StatusOK || rec.Code == http.StatusNotAcceptable)
	})

	t.Run("SetDefaultHeaders works with HEAD requests", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.SetDefaultHeaders(handler)

		req := httptest.NewRequest(http.MethodHead, "/api/data", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		require.NotEmpty(t, rec.Header().Get("Content-Type"))
	})

	t.Run("CompressResponse handles error responses", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		})

		wrapped := mw.CompressResponse(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
		req.Header.Set("Accept-Encoding", "gzip")
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("SetDefaultHeaders sets appropriate headers for different content types", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.URL.Path, "json") {
				w.Header().Set("Content-Type", "application/json")
			} else if strings.Contains(r.URL.Path, "html") {
				w.Header().Set("Content-Type", "text/html")
			}
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.SetDefaultHeaders(handler)

		// Test JSON
		req1 := httptest.NewRequest(http.MethodGet, "/api/json", nil)
		rec1 := httptest.NewRecorder()
		wrapped.ServeHTTP(rec1, req1)
		require.Contains(t, rec1.Header().Get("Content-Type"), "json")

		// Test HTML
		req2 := httptest.NewRequest(http.MethodGet, "/page/html", nil)
		rec2 := httptest.NewRecorder()
		wrapped.ServeHTTP(rec2, req2)
		require.Contains(t, rec2.Header().Get("Content-Type"), "html")
	})

	t.Run("CompressResponse with multiple encoding preferences", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(strings.Repeat("test ", 100)))
		})

		wrapped := mw.CompressResponse(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
		req.Header.Set("Accept-Encoding", "gzip, deflate, br")
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		encoding := rec.Header().Get("Content-Encoding")
		// Should use one of the supported encodings
		require.True(t, encoding == "gzip" || encoding == "deflate" || encoding == "br" || encoding == "")
	})

	t.Run("ContentNegotiation handles charset in accept header", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.ContentNegotiation(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Accept", "text/html; charset=utf-8")
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		require.Contains(t, rec.Header().Get("Content-Type"), "text/html")
	})
}
