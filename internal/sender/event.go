package sender

type EventSender interface {
	Send(solution *SolutionEvent) error
}
