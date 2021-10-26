package sender

import "github.com/ozonmp/edu-solution-api/internal/model"

type EventSender interface {
	Send(solution *model.SolutionEvent) error
}
