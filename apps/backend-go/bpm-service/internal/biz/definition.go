package biz

import (
	"context"
	"time"
)

// ProcessDefinition represents a BPMN process definition deployed to Zeebe.
type ProcessDefinition struct {
	ID                 string
	ProcessKey         string
	Name               string
	Description        string
	Category           string
	Module             string
	BPMNXml            string
	Version            int32
	ZeebeDeploymentKey int64
	IsActive           bool
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type DefinitionFilter struct {
	Module         string
	Keyword        string
	IncludeInactive bool
	PageSize       int
	PageToken      string
}

// DefinitionRepo interface for process definition persistence.
type DefinitionRepo interface {
	List(ctx context.Context, filter DefinitionFilter) ([]*ProcessDefinition, string, error)
	GetByID(ctx context.Context, id string) (*ProcessDefinition, error)
	Create(ctx context.Context, def *ProcessDefinition) (*ProcessDefinition, error)
}

// DefinitionUseCase handles business logic for process definitions.
type DefinitionUseCase struct {
	repo DefinitionRepo
}

func NewDefinitionUseCase(repo DefinitionRepo) *DefinitionUseCase {
	return &DefinitionUseCase{repo: repo}
}

func (uc *DefinitionUseCase) List(ctx context.Context, filter DefinitionFilter) ([]*ProcessDefinition, string, error) {
	return uc.repo.List(ctx, filter)
}

func (uc *DefinitionUseCase) GetByID(ctx context.Context, id string) (*ProcessDefinition, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *DefinitionUseCase) GetDiagram(ctx context.Context, id string) (*ProcessDefinition, error) {
	def, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return def, nil
}

func (uc *DefinitionUseCase) Deploy(ctx context.Context, def *ProcessDefinition, deployer func(string, string) (int64, error)) (*ProcessDefinition, error) {
	if deployer != nil {
		deploymentKey, err := deployer(def.BPMNXml, def.ProcessKey+".bpmn")
		if err != nil {
			return nil, err
		}
		def.ZeebeDeploymentKey = deploymentKey
	}
	def.Version = 1
	return uc.repo.Create(ctx, def)
}
