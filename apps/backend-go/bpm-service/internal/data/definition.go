package data

import (
	"context"
	"fmt"
	"strings"

	"github.com/arda-labs/arda/arda-be-go/pkg/pagination"
	"github.com/arda-labs/arda/arda-be-go/services/bpm-service/internal/biz"
)

type definitionRepo struct {
	data *Data
}

func NewDefinitionRepo(data *Data) biz.DefinitionRepo {
	return &definitionRepo{data: data}
}

func (r *definitionRepo) List(ctx context.Context, filter biz.DefinitionFilter) ([]*biz.ProcessDefinition, string, error) {
	params := pagination.Normalize(filter.PageSize, filter.PageToken)

	where := []string{"deleted_at IS NULL"}
	args := []interface{}{}
	argIdx := 1

	if filter.Module != "" {
		where = append(where, fmt.Sprintf("module = $%d", argIdx))
		args = append(args, filter.Module)
		argIdx++
	}
	if !filter.IncludeInactive {
		where = append(where, "is_active = true")
	}
	if filter.Keyword != "" {
		where = append(where, fmt.Sprintf("(name ILIKE $%d OR description ILIKE $%d)", argIdx, argIdx))
		args = append(args, "%"+filter.Keyword+"%")
		argIdx++
	}

	query := fmt.Sprintf(`
		SELECT id, process_key, name, description, category, module, version,
		       zeebe_deployment_key, is_active, created_at, updated_at
		FROM process_definitions
		WHERE %s
		ORDER BY updated_at DESC
		LIMIT $%d OFFSET $%d`, strings.Join(where, " AND "), argIdx, argIdx+1)

	args = append(args, params.Limit, params.Offset)

	rows, err := r.data.DB(ctx).Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, "", fmt.Errorf("query definitions: %w", err)
	}
	defer rows.Close()

	var defs []*biz.ProcessDefinition
	for rows.Next() {
		d := &biz.ProcessDefinition{}
		if err := rows.Scan(&d.ID, &d.ProcessKey, &d.Name, &d.Description, &d.Category,
			&d.Module, &d.Version, &d.ZeebeDeploymentKey, &d.IsActive, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, "", fmt.Errorf("scan definition: %w", err)
		}
		defs = append(defs, d)
	}

	nextToken := ""
	if len(defs) > 0 {
		nextToken = pagination.NextOffsetToken(len(defs), params.Limit, params.Offset)
	}

	return defs, nextToken, rows.Err()
}

func (r *definitionRepo) GetByID(ctx context.Context, id string) (*biz.ProcessDefinition, error) {
	query := `SELECT id, process_key, name, description, category, module, bpmn_xml, version,
	                  zeebe_deployment_key, is_active, created_at, updated_at
	           FROM process_definitions WHERE id = $1 AND deleted_at IS NULL`

	d := &biz.ProcessDefinition{}
	err := r.data.DB(ctx).Pool.QueryRow(ctx, query, id).Scan(
		&d.ID, &d.ProcessKey, &d.Name, &d.Description, &d.Category,
		&d.Module, &d.BPMNXml, &d.Version, &d.ZeebeDeploymentKey, &d.IsActive,
		&d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get definition by id: %w", err)
	}
	return d, nil
}

func (r *definitionRepo) Create(ctx context.Context, def *biz.ProcessDefinition) (*biz.ProcessDefinition, error) {
	query := `INSERT INTO process_definitions (process_key, name, description, category, module, bpmn_xml, version, zeebe_deployment_key, is_active)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	          RETURNING id, created_at, updated_at`

	d := &biz.ProcessDefinition{}
	err := r.data.DB(ctx).Pool.QueryRow(ctx, query,
		def.ProcessKey, def.Name, def.Description, def.Category, def.Module,
		def.BPMNXml, def.Version, def.ZeebeDeploymentKey, true).Scan(
		&d.ID, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("create definition: %w", err)
	}

	d.ProcessKey = def.ProcessKey
	d.Name = def.Name
	d.Description = def.Description
	d.Category = def.Category
	d.Module = def.Module
	d.BPMNXml = def.BPMNXml
	d.Version = def.Version
	d.ZeebeDeploymentKey = def.ZeebeDeploymentKey
	d.IsActive = true
	return d, nil
}
