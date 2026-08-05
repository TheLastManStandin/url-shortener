package save

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"url-shortener/internal/http-server/handlers/url/save/mocks"
	"url-shortener/internal/storage"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name     string
		url      string
		alias    string
		respErr  string
		mockErr  error
		respCode int
	}{
		{
			name:     "success",
			url:      "http://test.com/",
			alias:    "test_alias",
			respCode: http.StatusOK,
		},
		{
			name:     "already_exists",
			url:      "http://test.com/",
			alias:    "test_alias",
			respErr:  "alias already exists",
			mockErr:  storage.ErrAliasExists,
			respCode: http.StatusConflict,
		},
		{
			name:     "empty_alias",
			url:      "http://test.com/",
			alias:    "",
			respCode: http.StatusOK,
		},
		{
			name:     "empty_URL",
			url:      "",
			alias:    "abj",
			respErr:  "field URL is a required field",
			respCode: http.StatusOK,
		},
		{
			name:     "Invalid URL",
			url:      "some invalid URL",
			alias:    "some_alias",
			respErr:  "field URL is not a valid URL",
			respCode: http.StatusOK,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			urLSaverMock := mocks.NewURLSaver(t)

			if tc.respErr == "" || tc.mockErr != nil {
				urLSaverMock.On("SaveURL", tc.url, mock.AnythingOfType("string")).
					Return(tc.mockErr).
					Once()
			}

			handler := New(slog.New(slog.DiscardHandler), urLSaverMock)

			input := fmt.Sprintf(`{"url":"%s", "alias":"%s"}`, tc.url, tc.alias)

			req, err := http.NewRequest(http.MethodPost, "/save", strings.NewReader(input))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, tc.respCode, rr.Code)

			body := rr.Body.String()

			var resp Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respErr, resp.Error)
		})
	}
}
