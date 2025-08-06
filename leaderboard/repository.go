package leaderboard

import (
	"cmp"
	"context"
	"github.com/google/uuid"
	"slices"
	"sync"
)

type (
	repository interface {
		Get(ctx context.Context, filter Filter) ([]Record, error)
		Set(ctx context.Context, records []Record) error
		GetLeaders(ctx context.Context, offset, limit int) ([]Record, error)

		// UpdateLeaderBoard update a leaderboard database to update places by score, return only place changed records
		UpdateLeaderBoard(ctx context.Context) ([]Record, error)
	}
)

type repoMock struct {
	Accounts    map[uuid.UUID]Record
	LeaderBoard map[int]uuid.UUID

	accessMutex sync.Mutex
}

func NewRepositoryMock() repository {
	return &repoMock{
		Accounts:    make(map[uuid.UUID]Record),
		LeaderBoard: make(map[int]uuid.UUID),
		accessMutex: sync.Mutex{},
	}
}

func (r *repoMock) GetLeaders(ctx context.Context, offset, limit int) ([]Record, error) {
	r.accessMutex.Lock()
	defer r.accessMutex.Unlock()

	result := make([]Record, 0, limit)

	for i := offset; i < offset+limit; i++ {
		if accountID, found := r.LeaderBoard[i]; found {
			if record, found := r.Accounts[accountID]; found {
				result = append(result, record)
			}
		}
	}

	return result, nil
}

func (r *repoMock) Get(ctx context.Context, f Filter) ([]Record, error) {
	r.accessMutex.Lock()
	defer r.accessMutex.Unlock()

	result := make([]Record, 0, len(f.AccountIDs))

	for _, accountID := range f.AccountIDs {
		if record, found := r.Accounts[accountID]; found {
			result = append(result, record)
		}
	}

	return result, nil
}

func (r *repoMock) Set(ctx context.Context, records []Record) error {
	r.accessMutex.Lock()
	defer r.accessMutex.Unlock()

	for _, record := range records {
		r.Accounts[record.AccountID] = record
	}

	return nil
}

func (r *repoMock) UpdateLeaderBoard(ctx context.Context) ([]Record, error) {
	r.accessMutex.Lock()
	defer r.accessMutex.Unlock()

	// get all records
	records := make([]Record, 0, len(r.Accounts))
	for _, record := range r.Accounts {
		records = append(records, record)
	}

	// sort it by score
	slices.SortFunc(records, func(a, b Record) int {
		return cmp.Compare(b.Score, a.Score)
	})

	// alloc a slice for changed place records
	changedRecords := make([]Record, 0, len(r.Accounts))
	for i := 0; i < len(records); i++ {
		newPlace := i + 1 // calculate a new place, places start from 1

		// if the place is changed, update it and append to slice
		if records[i].Place != newPlace {
			records[i].Place = newPlace
			changedRecords = append(changedRecords, records[i])
		}
	}

	// update a database records and a leaderboard with new data
	for _, record := range changedRecords {
		r.Accounts[record.AccountID] = record
		r.LeaderBoard[record.Place] = record.AccountID
	}

	return changedRecords, nil
}
