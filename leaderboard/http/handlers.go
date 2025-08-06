package http

import (
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"test-service/helpers"
	"test-service/leaderboard"
)

func Handlers(log *slog.Logger, leaderboardSrv leaderboard.Service) http.Handler {
	r := chi.NewRouter()

	r.Post("/record", func(w http.ResponseWriter, r *http.Request) {
		helpers.LogInfoHTTP(log, r, "record", "leaderboard")
		helpers.CallService(w, r, leaderboardSrv.CreateRecord)
	})
	r.Post("/leaders", func(w http.ResponseWriter, r *http.Request) {
		helpers.LogInfoHTTP(log, r, "leaders", "leaderboard")
		helpers.CallService(w, r, leaderboardSrv.GetLeaders)
	})
	r.Post("/me", func(w http.ResponseWriter, r *http.Request) {
		helpers.LogInfoHTTP(log, r, "me", "leaderboard")
		helpers.CallService(w, r, leaderboardSrv.GetByAccount)
	})

	return r
}
