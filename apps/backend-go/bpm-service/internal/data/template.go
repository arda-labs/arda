package data

import (
	"context"
	"fmt"
	"strings"

	"github.com/arda-labs/arda/arda-be-go/pkg/pagination"
	"github.com/arda-labs/arda/arda-be-go/services/bpm-service/internal/biz"
)

type templateRepo struct {
	data *Data
}

func NewTemplateRepo(data *Data) biz.TemplateRepo {
	return &templateRepo{data: data}
}

func (r *templateRepo) List(ctx context.Context, filter biz.TemplateFilter) ([]*biz.Template, string, error) {
	params := pagination.Normalize(filter.PageSize, filter.PageToken)

	where := []string{"t.deleted_at IS NULL"}
	args := []interface{}{}
	argIdx := 1

	if filter.ProcessDefinitionID != "" {
		where = append(where, fmt.Sprintf("t.process_definition_id = $%d", argIdx))
		args = append(args, filter.ProcessDefinitionID)
		argIdx++
	}
	if filter.Module != "" {
		where = append(where, fmt.Sprintf("t.module = $%d", argIdx))
		args = append(args, filter.Module)
		argIdx++
	}
	if filter.Keyword != "" {
		where = append(where, fmt.Sprintf("(t.name ILIKE $%d OR t.template_text ILIKE $%d)", argIdx, argIdx))
		args = append(args, "%"+filter.Keyword+"%")
		argIdx++
	}

	query := fmt.Sprintf(`
		SELECT t.id, t.process_definition_id, t.name, t.template_text, t.module, t.created_at, t.updated_at
		FROM process_templates t
		WHERE %s
		ORDER BY t.updated_at DESC
		LIMIT $%d OFFSET $%d`, strings.Join(where, " AND "), argIdx, argIdx+1)

	args = append(args, params.Limit, params.Offset)

	rows, err := r.data.DB(ctx).Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, "", fmt.Errorf("query templates: %w", err)
	}
	defer rows.Close()

	var templates []*biz.Template
	for rows.Next() {
		t := &biz.Template{}
		if err := rows.Scan(&t.ID, &t.ProcessDefinitionID, &t.Name, &t.TemplateText, &t.Module, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, "", fmt.Errorf("scan template: %w", err)
		}
		// Load variables for each template
		vars, err := r.loadVariables(ctx, t.ID)
		if err != nil {
			return nil, "", err
		}
		t.Variables = vars
		templates = append(templates, t)
	}

	nextToken := ""
	if len(templates) > 0 {
		nextToken = pagination.NextOffsetToken(len(templates), params.Limit, params.Offset)
	}

	return templates, nextToken, rows.Err()
}

func (r *templateRepo) GetByID(ctx context.Context, id string) (*biz.Template, error) {
	query := `SELECT id, process_definition_id, name, template_text, module, created_at, updated_at
	           FROM process_templates WHERE id = $1 AND deleted_at IS NULL`

	t := &biz.Template{}
	err := r.data.DB(ctx).Pool.QueryRow(ctx, query, id).Scan(
		&t.ID, &t.ProcessDefinitionID, &t.Name, &t.TemplateText, &t.Module, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get template by id: %w", err)
	}

	vars, err := r.loadVariables(ctx, t.ID)
	if err != nil {
		return nil, err
	}
	t.Variables = vars
	return t, nil
}

func (r *templateRepo) loadVariables(ctx context.Context, templateID string) ([]*biz.TemplateVariable, error) {
	query := `SELECT id, template_id, variable_name, source_type, source_field, resolver_config, fallback_value
	           FROM template_variables WHERE template_id = $1 ORDER BY variable_name`

	rows, err := r.data.DB(ctx).Pool.Query(ctx, query, templateID)
	if err != nil {
		return nil, fmt.Errorf("query template variables: %w", err)
	}
	defer rows.Close()

	var vars []*biz.TemplateVariable
	for rows.Next() {
		v := &biz.TemplateVariable{}
		if err := rows.Scan(&v.ID, &v.TemplateID, &v.VariableName, &v.SourceType, &v.SourceField, &v.ResolverConfig, &v.FallbackValue); err != nil {
			return nil, fmt.Errorf("scan template variable: %w", err)
		}
		vars = append(vars, v)
	}
	return vars, rows.Err()
}

func (r *templateRepo) Create(ctx context.Context, tpl *biz.Template) (*biz.Template, error) {
	query := `INSERT INTO process_templates (process_definition_id, name, template_text, module)
	          VALUES ($1, $2, $3, $4)
	          RETURNING id, created_at, updated_at`

	created := &biz.Template{}
	err := r.data.DB(ctx).Pool.QueryRow(ctx, query,
		tpl.ProcessDefinitionID, tpl.Name, tpl.TemplateText, tpl.Module).Scan(
		&created.ID, &created.CreatedAt, &created.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("create template: %w", err)
	}

	created.ProcessDefinitionID = tpl.ProcessDefinitionID
	created.Name = tpl.Name
	created.TemplateText = tpl.TemplateText
	created.Module = tpl.Module

	// Save variables
	for _, v := range tpl.Variables {
		v.TemplateID = created.ID
		saved, err := r.CreateVariableMapping(ctx, v)
		if err != nil {
			return nil, err
		}
		created.Variables = append(created.Variables, saved)
	}

	return created, nil
}

func (r *templateRepo) Update(ctx context.Context, tpl *biz.Template) (*biz.Template, error) {
	query := `UPDATE process_templates SET name = $1, template_text = $2, module = $3, updated_at = now()
	          WHERE id = $4 AND deleted_at IS NULL
	          RETURNING updated_at`

	err := r.data.DB(ctx).Pool.QueryRow(ctx, query, tpl.Name, tpl.TemplateText, tpl.Module, tpl.ID).Scan(&tpl.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("update template: %w", err)
	}

	// Update variables
	// Delete existing then re-insert
	delQuery := `DELETE FROM template_variables WHERE template_id = $1`
	if _, err := r.data.DB(ctx).Pool.Exec(ctx, delQuery, tpl.ID); err != nil {
		return nil, fmt.Errorf("delete old variables: %w", err)
	}

	for _, v := range tpl.Variables {
		v.TemplateID = tpl.ID
		if _, err := r.CreateVariableMapping(ctx, v); err != nil {
			return nil, err
		}
	}

	return tpl, nil
}

func (r *templateRepo) Delete(ctx context.Context, id string) error {
	query := `UPDATE process_templates SET deleted_at = now() WHERE id = $1`
	_, err := r.data.DB(ctx).Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete template: %w", err)
	}
	return nil
}

func (r *templateRepo) ListVariableSources(ctx context.Context) ([]*biz.VariableSource, error) {
	return []*biz.VariableSource{
		{Type: "INSTANCE_VAR", Name: "Biến tiến trình", Description: "Dữ liệu từ biến của process instance (variables JSONB)"},
		{Type: "EVENT_PAYLOAD", Name: "Dữ liệu sự kiện", Description: "Dữ liệu từ sự kiện gần nhất của process instance"},
		{Type: "DB_LOOKUP", Name: "Tra cứu CSDL", Description: "Tra cứu từ cơ sở dữ liệu qua cấu hình SQL"},
		{Type: "CUSTOM", Name: "Tùy chỉnh", Description: "Resolver tùy chỉnh theo code"},
	}, nil
}

func (r *templateRepo) CreateVariableMapping(ctx context.Context, v *biz.TemplateVariable) (*biz.TemplateVariable, error) {
	query := `INSERT INTO template_variables (template_id, variable_name, source_type, source_field, resolver_config, fallback_value)
	          VALUES ($1, $2, $3, $4, $5::jsonb, $6)
	          RETURNING id`

	saved := &biz.TemplateVariable{}
	err := r.data.DB(ctx).Pool.QueryRow(ctx, query,
		v.TemplateID, v.VariableName, v.SourceType, v.SourceField, v.ResolverConfig, v.FallbackValue).Scan(&saved.ID)
	if err != nil {
		return nil, fmt.Errorf("create variable mapping: %w", err)
	}
	saved.TemplateID = v.TemplateID
	saved.VariableName = v.VariableName
	saved.SourceType = v.SourceType
	saved.SourceField = v.SourceField
	saved.ResolverConfig = v.ResolverConfig
	saved.FallbackValue = v.FallbackValue
	return saved, nil
}
