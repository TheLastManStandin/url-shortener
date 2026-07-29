package save

import (
	"log/slog"
	"net/http"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/random"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	Alias string `json:"alias,omitempty"`
	URL   string `json:"url" validate:"required, url"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

type URLSaver interface {
	SaveURL(urlToSave string, alias string) error
}

const randAliasLength = 6

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
			log.Error("failed to decode request body", err)
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		log.Debug("request body decoded", slog.Any("req", req))

		if err := validator.New().Struct(req); err != nil {
			responseWithErrors := resp.ValidationError(err.(validator.ValidationErrors))
			log.Error("invalid request", err)
			//log.Error("failed to validate request", sl.Err(err))
			render.JSON(w, r, responseWithErrors)
			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomAlias(randAliasLength)
		}
	}
}
