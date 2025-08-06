package leaderboard

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"test-service/events"
)

type (
	Service interface {
		// GetLeaders return a list of leaderboard with pagination
		GetLeaders(ctx context.Context, r GetLeadersRequest) (GetLeadersResponse, error)
		// GetByAccount return info about account in leaderboard
		GetByAccount(ctx context.Context, r GetByAccountRequest) (GetByAccountResponse, error)
		// CreateRecord add new record in leader board and recalculate it
		CreateRecord(ctx context.Context, r CreateRecordRequest) (CreateRecordResponse, error)
	}
)

type srv struct {
	repo repository

	eventsSrv events.Service
}

func NewService(repo repository, eventsSrv events.Service, mws ...func(s Service) Service) Service {
	var s Service
	s = &srv{repo: repo, eventsSrv: eventsSrv}
	for _, mw := range mws {
		s = mw(s)
	}

	return s
}

func (s *srv) GetLeaders(ctx context.Context, r GetLeadersRequest) (GetLeadersResponse, error) {
	records, err := s.repo.GetLeaders(ctx, r.Offset, r.Limit)
	if err != nil {
		return GetLeadersResponse{}, err
	}

	return GetLeadersResponse{Records: records}, nil
}

func (s *srv) GetByAccount(ctx context.Context, r GetByAccountRequest) (GetByAccountResponse, error) {
	records, err := s.repo.Get(ctx, Filter{AccountIDs: []uuid.UUID{r.AccountID}})
	if err != nil {
		return GetByAccountResponse{}, err
	}

	if len(records) == 0 {
		return GetByAccountResponse{}, errors.New("record not found")
	}

	return GetByAccountResponse{records[0]}, nil
}

func (s *srv) CreateRecord(ctx context.Context, r CreateRecordRequest) (CreateRecordResponse, error) {
	// set a new record
	err := s.repo.Set(ctx, []Record{{
		AccountID: r.AccountID,
		Score:     r.Score,
	}})
	if err != nil {
		return CreateRecordResponse{}, err
	}

	// update a leaderboard
	changes, err := s.repo.UpdateLeaderBoard(ctx)
	if err != nil {
		return CreateRecordResponse{}, err
	}

	// flush changed records to events listeners
	for _, change := range changes {
		s.eventsSrv.BroadcastEvent(ctx, events.LeaderboardChanges, change)
	}

	return CreateRecordResponse{}, nil
}
