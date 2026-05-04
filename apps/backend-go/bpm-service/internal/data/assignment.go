package data

import (
	"context"
	"sync"

	"github.com/arda-labs/arda/arda-be-go/services/bpm-service/internal/biz"
)

type assignmentRepo struct {
	data map[string]*biz.AssignmentRule
	mu   sync.RWMutex
}

func NewAssignmentRepo(data *Data) biz.AssignmentRepo {
	repo := &assignmentRepo{
		data: make(map[string]*biz.AssignmentRule),
	}

	repo.data["CRM"] = &biz.AssignmentRule{
		ID:       "rule-1",
		Module:   "CRM",
		Strategy: biz.RoundRobin,
		Agents: []*biz.Agent{
			{ID: "A01", Name: "Nhân viên A", Active: true, Weight: 1},
			{ID: "A02", Name: "Nhân viên B", Active: true, Weight: 1},
		},
	}

	repo.data["LOAN"] = &biz.AssignmentRule{
		ID:       "rule-2",
		Module:   "LOAN",
		Strategy: biz.LoadBalance,
		Agents: []*biz.Agent{
			{ID: "L01", Name: "Thẩm định viên X", Active: true, Weight: 1, TaskCount: 5},
			{ID: "L02", Name: "Thẩm định viên Y", Active: true, Weight: 1, TaskCount: 2},
		},
	}

	return repo
}

func (r *assignmentRepo) GetRuleByModule(ctx context.Context, module string) (*biz.AssignmentRule, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if rule, ok := r.data[module]; ok {
		return rule, nil
	}

	return &biz.AssignmentRule{
		Module:   module,
		Strategy: biz.RoundRobin,
	}, nil
}

func (r *assignmentRepo) UpdateRule(ctx context.Context, rule *biz.AssignmentRule) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[rule.Module] = rule
	return nil
}
