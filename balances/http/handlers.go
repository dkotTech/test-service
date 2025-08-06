package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"test-service/balances"
	"test-service/helpers"
)

func Handlers(log *slog.Logger, balancesSrv balances.Service) http.Handler {
	r := chi.NewRouter()

	r.Post("/deposit", func(w http.ResponseWriter, r *http.Request) {
		helpers.LogInfoHTTP(log, r, "deposit", "balances")
		helpers.CallService(w, r, balancesSrv.Deposit)
	})
	r.Post("/withdraw", func(w http.ResponseWriter, r *http.Request) {
		helpers.LogInfoHTTP(log, r, "withdraw", "balances")
		helpers.CallService(w, r, balancesSrv.Withdraw)
	})
	r.Get("/balance/{account_id}", func(w http.ResponseWriter, r *http.Request) {
		helpers.LogInfoHTTP(log, r, "balance", "balances")
		ctx := r.Context()
		var (
			accountID uuid.UUID
			err       error
		)

		if accountIDStr := chi.URLParam(r, "account_id"); accountIDStr != "" {
			accountID, err = uuid.Parse(accountIDStr)
			if err != nil {
				helpers.SendError(ctx, w, err)
			}
		}

		res, err := balancesSrv.CurrentOne(ctx, balances.GetCurrentOneRequest{AccountID: accountID})
		if err != nil {
			helpers.SendError(ctx, w, err)
			return
		}

		_ = helpers.Encode(ctx, w, &res)
	})

	return r
}
