package redirect

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"url-shortener/internal/http-server/handlers/url/redirect/mocks"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		alias     string
		sqlReturn string
		respErr   string
		mockErr   error
		respCode  int
	}{
		{
			name:      "success",
			alias:     "valid_alias",
			respCode:  http.StatusFound,
			sqlReturn: "https://www.valid_url.com",
		},
		{
			name:     "empty_alias",
			alias:    "",
			respCode: http.StatusNotFound,
			respErr:  "not found",
		},
		{
			name:     "empty_url",
			alias:    "aboba",
			respCode: http.StatusOK,
			mockErr:  storage.ErrAliasNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			urlGetterMock := mocks.NewURLGetter(t)

			if tt.respErr == "" || tt.mockErr != nil {
				urlGetterMock.On("GetURL", tt.alias).
					Return(tt.sqlReturn, tt.mockErr).
					Once()
			}
			router := chi.NewRouter()

			router.Get("/{alias}", New(slog.New(slog.DiscardHandler), urlGetterMock))

			req, err := http.NewRequest(http.MethodGet, "/"+tt.alias, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			require.Equal(t, tt.respCode, rr.Code)
			require.Equal(t, tt.sqlReturn, rr.Header().Get("Location"))
		})
	}
}
