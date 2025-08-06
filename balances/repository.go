package balances

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"sync"
)

var (
	accountAlreadyExist = errors.New("account already exists")
)

type (
	repository interface {
		SetOperations(ctx context.Context, ops []Operation) error
		GetOperations(ctx context.Context, f Filter) ([]Operation, error)
	}
)

type repoMock struct {
	Balances map[uuid.UUID][]Operation

	accessMutex sync.Mutex
}

func NewRepositoryMock() repository {
	return &repoMock{
		Balances:    make(map[uuid.UUID][]Operation),
		accessMutex: sync.Mutex{},
	}
}

func (r *repoMock) SetOperations(ctx context.Context, ops []Operation) error {
	r.accessMutex.Lock()
	defer r.accessMutex.Unlock()

	for _, op := range ops {
		if op.ID == uuid.Nil {
			op.ID = uuid.New()
		}

		r.Balances[op.AccountID] = append(r.Balances[op.AccountID], op)
	}

	return nil
}

func (r *repoMock) GetOperations(ctx context.Context, f Filter) ([]Operation, error) {
	r.accessMutex.Lock()
	defer r.accessMutex.Unlock()

	result := make([]Operation, 0) // its hard to determinate how many operations founded, todo it when using a production database

	for _, accountID := range f.AccountIDs {
		if ops, found := r.Balances[accountID]; found {
			result = append(result, ops...)
		}
	}

	return result, nil
}
