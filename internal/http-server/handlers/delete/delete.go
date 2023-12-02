package delete

import (
	"log/slog"
	"net/http"
	"url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

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

		alias := chi.URLParam(r, "alias")

		if alias == "" {
			logger.Error("Alias is empty")

			render.JSON(w, r, response.Error("Invalid request"))
		}

		err := urlDeleter.DeleteURL(alias)

		if err != nil {
			logger.Error("Error while URL deleting", sl.Err(err))

			render.JSON(w, r, response.Error("Error while URL deleting"))

			return
		}

		logger.Info("URL deleted")

		render.JSON(w, r, response.OK())
	}
}
