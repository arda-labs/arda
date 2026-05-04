package biz

import "context"

// EventUseCase handles business logic for process events.
type EventUseCase struct {
	eventRepo EventRepo
}

func NewEventUseCase(eventRepo EventRepo) *EventUseCase {
	return &EventUseCase{eventRepo: eventRepo}
}

func (uc *EventUseCase) Create(ctx context.Context, event *ProcessEvent) (*ProcessEvent, error) {
	return uc.eventRepo.Create(ctx, event)
}
