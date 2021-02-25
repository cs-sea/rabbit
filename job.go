package rabbit

import (
	"context"
)

type Job interface {
	RabbitService
}

type JobService struct {
	m *Message
	RabbitService
}

func NewJobService(m *Message, rabbitService RabbitService) Job {
	return &JobService{m: m, RabbitService: rabbitService}
}

func (j *JobService) Attempt() int {
	return j.m.Attempt
}

func (j *JobService) Ack() error {
	return j.m.D.Ack(false)
}

func (j *JobService) Retry(ctx context.Context) error {
	return nil
}
