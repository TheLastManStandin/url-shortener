package save

import (
	"errors"
	"log/slog"
	"net/http"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/random"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	Alias string `json:"alias,omitempty"`
	URL   string `json:"url" validate:"required,url"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

type URLSaver interface {
	SaveURL(urlToSave string, alias string) error
}

const (
	randAliasLength      = 6
	maxAliasSaveAttempts = 10
)

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			//log.Error("failed to decode request", sl.Err(err))
			log.Error("failed to decode request body", slog.Any("error", err))
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		log.Debug("request body decoded", slog.Any("req", req))

		if err := validator.New().Struct(req); err != nil {
			responseWithErrors := resp.ValidationError(err.(validator.ValidationErrors))
			log.Info("invalid request", slog.Any("error", err))
			//log.Error("failed to validate request", sl.Err(err))
			render.JSON(w, r, responseWithErrors)
			return
		}

		alias := req.Alias
		isGeneratedAlias := alias == ""
		attempts := 1
		if isGeneratedAlias {
			attempts = maxAliasSaveAttempts
		}

		for attempt := 0; attempt < attempts; attempt++ {
			if isGeneratedAlias {
				alias = random.NewRandomAlias(randAliasLength)
			}

			err = urlSaver.SaveURL(req.URL, alias)
			if err == nil {
				render.JSON(w, r, Response{
					Response: resp.Ok(),
					Alias:    alias,
				})
				return
			}

			if !errors.Is(err, storage.ErrAliasExists) {
				log.Error("failed to save url", slog.Any("error", err))
				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, resp.Error("failed to save url"))
				return
			}

			if !isGeneratedAlias {
				log.Info("alias already exists", slog.String("alias", alias))
				render.Status(r, http.StatusConflict)
				render.JSON(w, r, resp.Error("alias already exists"))
				return
			}
		}

		log.Error(
			"failed to generate unique alias",
			slog.Int("attempts", maxAliasSaveAttempts),
		)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, resp.Error("failed to generate unique alias"))
		return
	}
}
