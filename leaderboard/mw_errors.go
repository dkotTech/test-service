package leaderboard

import (
	"context"
	"test-service/helpers"
)

type mwErrors struct {
	inner Service
}

func NewMWErrors(inner Service) Service {
	return &mwErrors{inner: inner}
}

func (m *mwErrors) GetLeaders(ctx context.Context, r GetLeadersRequest) (GetLeadersResponse, error) {
	return helpers.ToServiceErrorWrap(ctx, r, m.inner.GetLeaders)
}

func (m *mwErrors) GetByAccount(ctx context.Context, r GetByAccountRequest) (GetByAccountResponse, error) {
	return helpers.ToServiceErrorWrap(ctx, r, m.inner.GetByAccount)
}

func (m *mwErrors) CreateRecord(ctx context.Context, r CreateRecordRequest) (CreateRecordResponse, error) {
	return helpers.ToServiceErrorWrap(ctx, r, m.inner.CreateRecord)
}
