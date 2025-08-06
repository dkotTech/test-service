package leaderboard

import "github.com/google/uuid"

type (
	// Record with account current place and score in the leaderboard
	Record struct {
		AccountID uuid.UUID `json:"account_id"`
		Score     float64   `json:"score"`
		Place     int       `json:"place"`
	}
)

type (
	GetLeadersRequest struct {
		Offset int `json:"offset"`
		Limit  int `json:"limit" validate:"required"`
	}

	GetLeadersResponse struct {
		Records []Record `json:"records"`
	}

	GetByAccountRequest struct {
		AccountID uuid.UUID `json:"account_id"`
	}

	GetByAccountResponse struct {
		Record Record `json:"record"`
	}

	CreateRecordRequest struct {
		AccountID uuid.UUID `json:"account_id"`
		Score     float64   `json:"score"`
	}
	CreateRecordResponse struct{}
)

type (
	Filter struct {
		AccountIDs []uuid.UUID `json:"account_ids"`
	}
)
