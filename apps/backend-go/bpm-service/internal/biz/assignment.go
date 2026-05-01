package biz

import (
	"context"
	"sync"
)

// EventPublisher interface for cross-biz publishing
type EventPublisher interface {
	Publish(ctx context.Context, topic string, event interface{}) error
}

// AssignmentStrategy defines how tasks are distributed
type AssignmentStrategy string

const (
	RoundRobin  AssignmentStrategy = "ROUND_ROBIN"
	LoadBalance AssignmentStrategy = "LOAD_BALANCE"
)

// Agent represents a human worker or system available for tasks
type Agent struct {
	ID         string
	Name       string
	Active     bool
	Weight     int // Priority/Capacity weight
	TaskCount  int // Currently assigned tasks (for Load Balance)
}

// AssignmentRule defines a set of agents and a strategy for a specific module/process
type AssignmentRule struct {
	ID        string
	Module    string // CRM, LOAN, etc.
	Strategy  AssignmentStrategy
	Agents    []*Agent
	lastIndex int // for Round Robin
	mu        sync.Mutex
}

// AssignmentRepo interface for persistence
type AssignmentRepo interface {
	GetRuleByModule(ctx context.Context, module string) (*AssignmentRule, error)
	UpdateRule(ctx context.Context, rule *AssignmentRule) error
}

type AssignmentUseCase struct {
	repo AssignmentRepo
}

func NewAssignmentUseCase(repo AssignmentRepo) *AssignmentUseCase {
	return &AssignmentUseCase{repo: repo}
}

// AssignTask finds the best agent for a task based on the module's rule
func (uc *AssignmentUseCase) AssignTask(ctx context.Context, module string) (string, error) {
	rule, err := uc.repo.GetRuleByModule(ctx, module)
	if err != nil {
		return "", err
	}

	rule.mu.Lock()
	defer rule.mu.Unlock()

	var selectedAgent *Agent

	// Filter active agents
	var activeAgents []*Agent
	for _, a := range rule.Agents {
		if a.Active {
			activeAgents = append(activeAgents, a)
		}
	}

	if len(activeAgents) == 0 {
		return "SYSTEM_QUEUE", nil // Fallback to system queue
	}

	switch rule.Strategy {
	case RoundRobin:
		rule.lastIndex = (rule.lastIndex + 1) % len(activeAgents)
		selectedAgent = activeAgents[rule.lastIndex]

	case LoadBalance:
		// Find agent with minimum tasks
		minTasks := -1
		for _, a := range activeAgents {
			if minTasks == -1 || a.TaskCount < minTasks {
				minTasks = a.TaskCount
				selectedAgent = a
			}
		}
	}

	if selectedAgent != nil {
		selectedAgent.TaskCount++
		return selectedAgent.ID, nil
	}

	return "SYSTEM_QUEUE", nil
}
