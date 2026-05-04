package biz

import (
	"context"
	"encoding/json"
	"fmt"
)

// variableResolver implements VariableResolver interface.
type variableResolver struct {
	instanceRepo InstanceRepo
	eventRepo    EventRepo
	customFuncs  map[string]ResolverFunc
}

// ResolverFunc is a custom resolver function.
type ResolverFunc func(ctx context.Context, instanceID string, config map[string]interface{}) (string, error)

func NewVariableResolver(instanceRepo InstanceRepo, eventRepo EventRepo) VariableResolver {
	return &variableResolver{
		instanceRepo: instanceRepo,
		eventRepo:    eventRepo,
		customFuncs:  defaultCustomResolvers(),
	}
}

func (r *variableResolver) Resolve(ctx context.Context, variable *TemplateVariable, instanceID string) (string, error) {
	switch variable.SourceType {
	case "INSTANCE_VAR":
		return r.resolveInstanceVar(ctx, instanceID, variable.SourceField)
	case "EVENT_PAYLOAD":
		return r.resolveEventPayload(ctx, instanceID, variable.SourceField)
	case "DB_LOOKUP":
		return r.resolveDBLookup(ctx, instanceID, variable)
	case "CUSTOM":
		return r.resolveCustom(ctx, instanceID, variable)
	default:
		return "", fmt.Errorf("unknown source type: %s", variable.SourceType)
	}
}

func (r *variableResolver) resolveInstanceVar(ctx context.Context, instanceID, field string) (string, error) {
	inst, err := r.instanceRepo.GetByID(ctx, instanceID)
	if err != nil {
		return "", err
	}
	if inst.Variables == "" {
		return "", nil
	}
	var vars map[string]interface{}
	if err := json.Unmarshal([]byte(inst.Variables), &vars); err != nil {
		return "", err
	}
	val, ok := vars[field]
	if !ok {
		return "", fmt.Errorf("variable field %s not found in instance variables", field)
	}
	return fmt.Sprintf("%v", val), nil
}

func (r *variableResolver) resolveEventPayload(ctx context.Context, instanceID, field string) (string, error) {
	events, _, err := r.eventRepo.ListByInstance(ctx, instanceID, 1, "")
	if err != nil {
		return "", err
	}
	if len(events) == 0 {
		return "", fmt.Errorf("no events found for instance %s", instanceID)
	}
	if events[0].Data == "" {
		return "", nil
	}
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(events[0].Data), &data); err != nil {
		return "", err
	}
	val, ok := data[field]
	if !ok {
		return "", fmt.Errorf("field %s not found in latest event data", field)
	}
	return fmt.Sprintf("%v", val), nil
}

func (r *variableResolver) resolveDBLookup(ctx context.Context, instanceID string, variable *TemplateVariable) (string, error) {
	// DB_LOOKUP resolver_config expects: {"table": "...", "key_field": "...", "value_field": "..."}
	if variable.ResolverConfig == "" {
		return "", fmt.Errorf("resolver_config required for DB_LOOKUP")
	}
	var config struct {
		Table      string `json:"table"`
		KeyField   string `json:"key_field"`
		ValueField string `json:"value_field"`
	}
	if err := json.Unmarshal([]byte(variable.ResolverConfig), &config); err != nil {
		return "", fmt.Errorf("invalid resolver_config: %w", err)
	}
	if config.Table == "" || config.ValueField == "" {
		return "", fmt.Errorf("table and value_field required for DB_LOOKUP")
	}
	return "", fmt.Errorf("DB_LOOKUP requires database access configured in resolver")
}

func (r *variableResolver) resolveCustom(ctx context.Context, instanceID string, variable *TemplateVariable) (string, error) {
	var config map[string]interface{}
	if variable.ResolverConfig != "" {
		if err := json.Unmarshal([]byte(variable.ResolverConfig), &config); err != nil {
			return "", fmt.Errorf("invalid resolver_config: %w", err)
		}
	}
	fn, ok := r.customFuncs[variable.SourceField]
	if !ok {
		return "", fmt.Errorf("custom resolver function %s not found", variable.SourceField)
	}
	return fn(ctx, instanceID, config)
}

func defaultCustomResolvers() map[string]ResolverFunc {
	return map[string]ResolverFunc{}
}
