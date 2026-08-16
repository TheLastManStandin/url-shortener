package tests

import (
	"net/http"
	"net/url"
	"testing"
	"url-shortener/internal/http-server/handlers/url/save"
	"url-shortener/internal/lib/random"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
)

func TestURLShortener_HappyPath(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   "localhost:8080",
	}

	e := httpexpect.Default(t, u.String())
	alias := random.NewRandomAlias(10)

	e.POST("/url").
		WithJSON(save.Request{
			URL:   gofakeit.URL(),
			Alias: alias,
		}).
		WithBasicAuth("user", "user").
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		ContainsKey("alias").
		Values().
		Value(0).
		IsEqual(alias)

	e.DELETE("/url/"+alias).
		WithBasicAuth("user", "user").
		Expect().
		Status(http.StatusOK)
}
