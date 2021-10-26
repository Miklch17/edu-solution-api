package repo

import (
	"github.com/ozonmp/edu-solution-api/internal/model"
)

type EventRepo interface {
	Lock(n uint64) ([]model.SolutionEvent, error)
	Unlock(eventIDs []uint64) error

	Add(event []model.SolutionEvent) error
	Remove(eventIDs []uint64) error
}
