package delete_url

import (
	"log/slog"
	"net/http"
	"url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	Alias string `json:"alias,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.36.1 --name=URLSaver
type URLDeleter interface {
	DeleteURL(alias string) error
}

func New(logger *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"

		logger = logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			logger.Error("Failed to decode request body", sl.Err(err))

			render.JSON(w, r, response.Error("Failed to decode request"))

			return
		}

		logger.Info("Request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateError := err.(validator.ValidationErrors)

			logger.Error("Invalid request", sl.Err(err))

			render.JSON(w, r, response.ValidationError(validateError))

			return
		}

		alias := req.Alias
		if alias == "" {
			logger.Error("Alias is empty")

			render.JSON(w, r, response.Error("Invalid request"))
		}

		err = urlDeleter.DeleteURL(alias)

		if err != nil {
			logger.Error("Error while URL deleting", sl.Err(err))

			render.JSON(w, r, response.Error("Error while URL deleting"))

			return
		}

		logger.Info("URL deleted")

		render.JSON(w, r, response.OK())
	}
}
