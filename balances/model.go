package balances

import "github.com/google/uuid"

type (
	// Operation one operation for account with balance change
	Operation struct {
		ID        uuid.UUID `json:"id"`
		AccountID uuid.UUID `json:"account_id"`
		Operation float64   `json:"operation"`
	}

	// Wallet account current balance on wallet
	Wallet struct {
		AccountID uuid.UUID `json:"account_id"`
		Balance   float64   `json:"balance"`
	}
)

type (
	GetCurrentOneRequest struct {
		AccountID uuid.UUID `json:"account_id" validate:"required"`
	}

	GetCurrentOneResponse struct {
		Wallet
	}

	DepositRequest struct {
		AccountID uuid.UUID `json:"account_id" validate:"required"`
		Amount    float64   `json:"amount" validate:"required,gt=0"`
	}
	DepositResponse struct{}

	WithdrawRequest struct {
		AccountID uuid.UUID `json:"account_id" validate:"required"`
		Amount    float64   `json:"amount" validate:"required,gt=0"`
	}
	WithdrawResponse struct{}
)

type (
	Filter struct {
		AccountIDs []uuid.UUID `json:"account_ids"`
	}
)
