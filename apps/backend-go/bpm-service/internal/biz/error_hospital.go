package biz

import (
	"context"
	"time"
)

// FailedTask represents a service task that failed and needs manual intervention
type FailedTask struct {
	ID           string    `json:"id"`
	InstanceID   string    `json:"instance_id"`
	ProcessName  string    `json:"process_name"`
	StepName     string    `json:"step_name"`
	Error        string    `json:"error"`
	Payload      string    `json:"payload"` // JSON string of variables
	RetryCount   int       `json:"retry_count"`
	LastAttempt  time.Time `json:"last_attempt"`
	Status       string    `json:"status"` // "PENDING", "RETRYING", "RESOLVED"
}

type ErrorHospitalRepo interface {
	GetFailedTasks(ctx context.Context) ([]*FailedTask, error)
	GetTaskByID(ctx context.Context, id string) (*FailedTask, error)
	UpdateTask(ctx context.Context, task *FailedTask) error
}

type ErrorHospitalUseCase struct {
	repo      ErrorHospitalRepo
	publisher EventPublisher // To re-emit fixed events to Kafka
}

func NewErrorHospitalUseCase(repo ErrorHospitalRepo, pub EventPublisher) *ErrorHospitalUseCase {
	return &ErrorHospitalUseCase{
		repo:      repo,
		publisher: pub,
	}
}

// RetryTask updates the payload and triggers a re-execution by emitting an event
func (uc *ErrorHospitalUseCase) RetryTask(ctx context.Context, taskID string, newPayload string) error {
	task, err := uc.repo.GetTaskByID(ctx, taskID)
	if err != nil {
		return err
	}

	task.Payload = newPayload
	task.Status = "RETRYING"
	task.LastAttempt = time.Now()
	task.RetryCount++

	if err := uc.repo.UpdateTask(ctx, task); err != nil {
		return err
	}

	// Logic: Emit to Kafka 'bpm-retry-topic' or directly back to generic worker
	return uc.publisher.Publish(ctx, "bpm-service-retry", task)
}
