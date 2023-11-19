package redirect

import (
	"errors"
	"log/slog"
	"net/http"
	"url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

//go:generate go run github.com/vektra/mockery/v2@v2.36.1 --name=URLSaver
type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(logger *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect.New"

		logger = logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")

		if alias == "" {
			logger.Error("Alias is empty")

			render.JSON(w, r, response.Error("Invalid request"))
		}

		url, err := urlGetter.GetURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			logger.Info("URL not found", slog.String("alias", alias))

			render.JSON(w, r, response.Error("URL not found"))

			return
		}

		if err != nil {
			logger.Error("Error while URL getting", sl.Err(err))

			render.JSON(w, r, response.Error("Internal error"))

			return
		}

		logger.Info("URL got", slog.String("url", url))

		http.Redirect(w, r, url, http.StatusFound)
	}
}
