package leaderboard

import (
	"context"
	"test-service/helpers/validate"
)

type mwValidate struct {
	inner Service
}

func NewMWValidate(inner Service) Service {
	return &mwValidate{inner: inner}
}

func (m *mwValidate) GetLeaders(ctx context.Context, r GetLeadersRequest) (GetLeadersResponse, error) {
	if err := validate.MustValidate(ctx, &r); err != nil {
		return GetLeadersResponse{}, err
	}

	return m.inner.GetLeaders(ctx, r)
}

func (m *mwValidate) GetByAccount(ctx context.Context, r GetByAccountRequest) (GetByAccountResponse, error) {
	if err := validate.MustValidate(ctx, &r); err != nil {
		return GetByAccountResponse{}, err
	}

	return m.inner.GetByAccount(ctx, r)
}

func (m *mwValidate) CreateRecord(ctx context.Context, r CreateRecordRequest) (CreateRecordResponse, error) {
	if err := validate.MustValidate(ctx, &r); err != nil {
		return CreateRecordResponse{}, err
	}

	return m.inner.CreateRecord(ctx, r)
}
