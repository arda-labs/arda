package data

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/arda-labs/arda/arda-be-go/pkg/pagination"
	"github.com/arda-labs/arda/arda-be-go/services/bpm-service/internal/biz"
)

type instanceRepo struct {
	data *Data
}

func NewInstanceRepo(data *Data) biz.InstanceRepo {
	return &instanceRepo{data: data}
}

func (r *instanceRepo) List(ctx context.Context, filter biz.InstanceFilter) ([]*biz.ProcessInstance, string, error) {
	params := pagination.Normalize(filter.PageSize, filter.PageToken)

	where := []string{"pi.deleted_at IS NULL"}
	args := []interface{}{}
	argIdx := 1

	if filter.ProcessDefinitionID != "" {
		where = append(where, fmt.Sprintf("pi.process_definition_id = $%d", argIdx))
		args = append(args, filter.ProcessDefinitionID)
		argIdx++
	}
	if filter.Status != "" {
		where = append(where, fmt.Sprintf("pi.status = $%d", argIdx))
		args = append(args, filter.Status)
		argIdx++
	}
	if filter.Module != "" {
		where = append(where, fmt.Sprintf("pd.module = $%d", argIdx))
		args = append(args, filter.Module)
		argIdx++
	}
	if filter.Keyword != "" {
		where = append(where, fmt.Sprintf("(pd.name ILIKE $%d OR pi.id::text ILIKE $%d)", argIdx, argIdx))
		args = append(args, "%"+filter.Keyword+"%")
		argIdx++
	}

	query := fmt.Sprintf(`
		SELECT pi.id, pi.zeebe_instance_key, pi.process_definition_id, pi.status,
		       pi.current_step, pi.assigned_agent, pi.sla_status, pi.created_at, pi.completed_at
		FROM process_instances pi
		JOIN process_definitions pd ON pd.id = pi.process_definition_id
		WHERE %s
		ORDER BY pi.created_at DESC
		LIMIT $%d OFFSET $%d`, strings.Join(where, " AND "), argIdx, argIdx+1)

	args = append(args, params.Limit, params.Offset)

	rows, err := r.data.DB(ctx).Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, "", fmt.Errorf("query instances: %w", err)
	}
	defer rows.Close()

	var instances []*biz.ProcessInstance
	for rows.Next() {
		inst := &biz.ProcessInstance{}
		var completedAt *time.Time
		if err := rows.Scan(&inst.ID, &inst.ZeebeInstanceKey, &inst.ProcessDefinitionID,
			&inst.Status, &inst.CurrentStep, &inst.AssignedAgent, &inst.SLAStatus,
			&inst.CreatedAt, &completedAt); err != nil {
			return nil, "", fmt.Errorf("scan instance: %w", err)
		}
		inst.CompletedAt = completedAt
		instances = append(instances, inst)
	}

	nextToken := ""
	if len(instances) > 0 {
		nextToken = pagination.NextOffsetToken(len(instances), params.Limit, params.Offset)
	}

	return instances, nextToken, rows.Err()
}

func (r *instanceRepo) ListByIDs(ctx context.Context, ids []string) ([]*biz.ProcessInstance, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	// Build parameterized query for list of IDs
	query := `SELECT id, zeebe_instance_key, process_definition_id, status, current_step,
	                  assigned_agent, sla_status, created_at, completed_at
	           FROM process_instances WHERE id = ANY($1) AND deleted_at IS NULL`

	rows, err := r.data.DB(ctx).Pool.Query(ctx, query, ids)
	if err != nil {
		return nil, fmt.Errorf("query instances by ids: %w", err)
	}
	defer rows.Close()

	var instances []*biz.ProcessInstance
	for rows.Next() {
		inst := &biz.ProcessInstance{}
		var completedAt *time.Time
		if err := rows.Scan(&inst.ID, &inst.ZeebeInstanceKey, &inst.ProcessDefinitionID,
			&inst.Status, &inst.CurrentStep, &inst.AssignedAgent, &inst.SLAStatus,
			&inst.CreatedAt, &completedAt); err != nil {
			return nil, fmt.Errorf("scan instance: %w", err)
		}
		inst.CompletedAt = completedAt
		instances = append(instances, inst)
	}
	return instances, rows.Err()
}

func (r *instanceRepo) GetByID(ctx context.Context, id string) (*biz.ProcessInstance, error) {
	query := `SELECT id, zeebe_instance_key, process_definition_id, status, current_step,
	                  variables, assigned_agent, sla_status, created_at, completed_at
	           FROM process_instances WHERE id = $1 AND deleted_at IS NULL`

	inst := &biz.ProcessInstance{}
	var completedAt *time.Time
	err := r.data.DB(ctx).Pool.QueryRow(ctx, query, id).Scan(
		&inst.ID, &inst.ZeebeInstanceKey, &inst.ProcessDefinitionID,
		&inst.Status, &inst.CurrentStep, &inst.Variables, &inst.AssignedAgent,
		&inst.SLAStatus, &inst.CreatedAt, &completedAt)
	if err != nil {
		return nil, fmt.Errorf("get instance by id: %w", err)
	}
	inst.CompletedAt = completedAt
	return inst, nil
}

func (r *instanceRepo) GetByZeebeKey(ctx context.Context, key int64) (*biz.ProcessInstance, error) {
	query := `SELECT id, zeebe_instance_key, process_definition_id, status, current_step,
	                  variables, assigned_agent, sla_status, created_at, completed_at
	           FROM process_instances WHERE zeebe_instance_key = $1 AND deleted_at IS NULL`

	inst := &biz.ProcessInstance{}
	var completedAt *time.Time
	err := r.data.DB(ctx).Pool.QueryRow(ctx, query, key).Scan(
		&inst.ID, &inst.ZeebeInstanceKey, &inst.ProcessDefinitionID,
		&inst.Status, &inst.CurrentStep, &inst.Variables, &inst.AssignedAgent,
		&inst.SLAStatus, &inst.CreatedAt, &completedAt)
	if err != nil {
		return nil, fmt.Errorf("get instance by zeebe key: %w", err)
	}
	inst.CompletedAt = completedAt
	return inst, nil
}

func (r *instanceRepo) Create(ctx context.Context, inst *biz.ProcessInstance) (*biz.ProcessInstance, error) {
	query := `INSERT INTO process_instances (zeebe_instance_key, process_definition_id, status, current_step, variables, assigned_agent, sla_status)
	          VALUES ($1, $2, $3, $4, $5, $6, $7)
	          RETURNING id, created_at`

	created := &biz.ProcessInstance{}
	err := r.data.DB(ctx).Pool.QueryRow(ctx, query,
		inst.ZeebeInstanceKey, inst.ProcessDefinitionID, inst.Status,
		inst.CurrentStep, inst.Variables, inst.AssignedAgent, inst.SLAStatus).Scan(
		&created.ID, &created.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("create instance: %w", err)
	}
	created.ZeebeInstanceKey = inst.ZeebeInstanceKey
	created.ProcessDefinitionID = inst.ProcessDefinitionID
	created.Status = inst.Status
	created.CurrentStep = inst.CurrentStep
	created.AssignedAgent = inst.AssignedAgent
	created.SLAStatus = inst.SLAStatus
	return created, nil
}

func (r *instanceRepo) Update(ctx context.Context, inst *biz.ProcessInstance) error {
	query := `UPDATE process_instances SET status = $1, current_step = $2, variables = $3,
	          assigned_agent = $4, sla_status = $5, completed_at = $6, updated_at = now()
	          WHERE id = $7 AND deleted_at IS NULL`
	_, err := r.data.DB(ctx).Pool.Exec(ctx, query,
		inst.Status, inst.CurrentStep, inst.Variables,
		inst.AssignedAgent, inst.SLAStatus, inst.CompletedAt, inst.ID)
	if err != nil {
		return fmt.Errorf("update instance: %w", err)
	}
	return nil
}
