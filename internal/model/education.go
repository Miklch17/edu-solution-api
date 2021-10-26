package model

import "fmt"

type solution struct {
	ID     uint64
	TaskID uint64 `json:"task_id"`
	StudentID   uint64 `json:"student_id"`
	Description string `json:"description"`
}

func (c *solution) String() string{
	return fmt.Sprintf("ID: %d TaskID: %d StudentID: %d Description: %s", c.ID, c.TaskID, c.StudentID, c.Description)
}

type EventType uint8

type EventStatus uint8

const (
	Created EventType = iota
	Updated
	Removed

	Deferred EventStatus = iota
	Processed
)

type SolutionEvent struct {
	ID     uint64
	Type   EventType
	Status EventStatus
	Entity *solution
}
func (s * SolutionEvent) SetStatus(typeSE EventType, status EventStatus) {
	s.Type = typeSE
	s.Status = status
}
func (s * SolutionEvent) Lock(n uint64) ([]SolutionEvent, error){
	dat := make([]SolutionEvent, 100)
	return dat, nil
}

func (s * SolutionEvent) Unlock(eventIDs []uint64) error{
	return nil
}

func (s * SolutionEvent) Add(event []SolutionEvent) error{
	return nil
}

func (s * SolutionEvent) Remove(eventIDs []uint64) error{
	return nil
}
