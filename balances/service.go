package balances

import (
	"context"
	"errors"
	"math"

	"github.com/google/uuid"

	"test-service/events"
)

type (
	Service interface {
		// CurrentOne request current account balance
		CurrentOne(ctx context.Context, r GetCurrentOneRequest) (GetCurrentOneResponse, error)

		// Deposit make a deposit operation on account
		Deposit(ctx context.Context, r DepositRequest) (DepositResponse, error)
		// Withdraw make a withdrawal operation on account
		Withdraw(ctx context.Context, r WithdrawRequest) (WithdrawResponse, error)
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

func (s *srv) CurrentOne(ctx context.Context, r GetCurrentOneRequest) (GetCurrentOneResponse, error) {
	ops, err := s.repo.GetOperations(ctx, Filter{AccountIDs: []uuid.UUID{r.AccountID}})
	if err != nil {
		return GetCurrentOneResponse{}, err
	}

	if len(ops) == 0 {
		return GetCurrentOneResponse{}, errors.New("operations not found")
	}

	var balance float64

	for _, op := range ops {
		balance += op.Operation
	}

	return GetCurrentOneResponse{Wallet{
		AccountID: r.AccountID,
		Balance:   balance,
	}}, nil
}

func (s *srv) Deposit(ctx context.Context, r DepositRequest) (DepositResponse, error) {
	err := s.repo.SetOperations(ctx, []Operation{{
		AccountID: r.AccountID,
		Operation: r.Amount,
	}})
	if err != nil {
		return DepositResponse{}, err
	}

	// flush event with new deposit
	s.eventsSrv.BroadcastEvent(ctx, events.Deposit, r)

	return DepositResponse{}, nil
}

func (s *srv) Withdraw(ctx context.Context, r WithdrawRequest) (WithdrawResponse, error) {
	// make amount negative
	r.Amount = math.Copysign(r.Amount, -1)

	accountBalance, err := s.CurrentOne(ctx, GetCurrentOneRequest{AccountID: r.AccountID})
	if err != nil {
		return WithdrawResponse{}, err
	}

	// calculate future balance and return error if its negative
	if accountBalance.Balance+r.Amount < 0 {
		return WithdrawResponse{}, errors.New("insufficient funds")
	}

	err = s.repo.SetOperations(ctx, []Operation{{
		AccountID: r.AccountID,
		Operation: r.Amount,
	}})
	if err != nil {
		return WithdrawResponse{}, err
	}

	// flush event with new withdrawal
	s.eventsSrv.BroadcastEvent(ctx, events.Withdraw, r)

	return WithdrawResponse{}, nil
}
