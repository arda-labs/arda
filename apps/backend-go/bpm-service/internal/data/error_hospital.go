package data

import (
	"context"

	"github.com/arda-labs/arda/arda-be-go/services/bpm-service/internal/biz"
)

type errorHospitalRepo struct {
	data *Data
}

func NewErrorHospitalRepo(data *Data) biz.ErrorHospitalRepo {
	return &errorHospitalRepo{data: data}
}

func (r *errorHospitalRepo) GetFailedTasks(ctx context.Context) ([]*biz.FailedTask, error) {
	return []*biz.FailedTask{}, nil
}

func (r *errorHospitalRepo) GetTaskByID(ctx context.Context, id string) (*biz.FailedTask, error) {
	return nil, nil
}

func (r *errorHospitalRepo) UpdateTask(ctx context.Context, task *biz.FailedTask) error {
	return nil
}
