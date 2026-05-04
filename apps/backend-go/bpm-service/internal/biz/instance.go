package biz

import (
	"context"
	"time"
)

type ProcessInstance struct {
	ID                 string
	ZeebeInstanceKey   int64
	ProcessDefinitionID string
	Status             string
	CurrentStep        string
	Variables          string // JSON string
	AssignedAgent      string
	SLAStatus          string
	CreatedAt          time.Time
	CompletedAt        *time.Time
}

type InstanceFilter struct {
	ProcessDefinitionID string
	Status              string
	Module              string
	Keyword             string
	FromDate            string
	ToDate              string
	PageSize            int
	PageToken           string
}

type InstanceDetail struct {
	ProcessInstance
	ProcessName          string
	BPMNXml              string
	ActiveElementIDs     []string
	CompletedElementIDs  []string
}

type HeatmapStep struct {
	ElementID         string
	ElementName       string
	InstanceCount     int32
	AvgDurationSeconds float64
	Severity          string
}

// InstanceRepo interface for process instance persistence.
type InstanceRepo interface {
	List(ctx context.Context, filter InstanceFilter) ([]*ProcessInstance, string, error)
	ListByIDs(ctx context.Context, ids []string) ([]*ProcessInstance, error)
	GetByID(ctx context.Context, id string) (*ProcessInstance, error)
	GetByZeebeKey(ctx context.Context, key int64) (*ProcessInstance, error)
	Create(ctx context.Context, instance *ProcessInstance) (*ProcessInstance, error)
	Update(ctx context.Context, instance *ProcessInstance) error
}

// EventRepo interface for process event persistence.
type EventRepo interface {
	ListByInstance(ctx context.Context, instanceID string, pageSize int, pageToken string) ([]*ProcessEvent, string, error)
	Create(ctx context.Context, event *ProcessEvent) (*ProcessEvent, error)
	GetHeatmap(ctx context.Context, definitionID string) ([]*HeatmapStep, error)
}

type ProcessEvent struct {
	ID                string
	ProcessInstanceID string
	EventType         string
	Source            string
	Data              string // JSON string
	Timestamp         time.Time
}

// InstanceUseCase handles business logic for process instances.
type InstanceUseCase struct {
	instanceRepo InstanceRepo
	eventRepo    EventRepo
	defRepo      DefinitionRepo
}

func NewInstanceUseCase(instanceRepo InstanceRepo, eventRepo EventRepo, defRepo DefinitionRepo) *InstanceUseCase {
	return &InstanceUseCase{
		instanceRepo: instanceRepo,
		eventRepo:    eventRepo,
		defRepo:      defRepo,
	}
}

func (uc *InstanceUseCase) List(ctx context.Context, filter InstanceFilter) ([]*ProcessInstance, string, error) {
	return uc.instanceRepo.List(ctx, filter)
}

func (uc *InstanceUseCase) Create(ctx context.Context, inst *ProcessInstance) (*ProcessInstance, error) {
	return uc.instanceRepo.Create(ctx, inst)
}

func (uc *InstanceUseCase) GetDetail(ctx context.Context, id string) (*InstanceDetail, error) {
	inst, err := uc.instanceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	def, err := uc.defRepo.GetByID(ctx, inst.ProcessDefinitionID)
	if err != nil {
		// Return instance detail even if we can't fetch the definition name
		def = &ProcessDefinition{Name: "Unknown"}
	}

	detail := &InstanceDetail{
		ProcessInstance: *inst,
		ProcessName:     def.Name,
		BPMNXml:         def.BPMNXml,
	}

	return detail, nil
}

func (uc *InstanceUseCase) GetEvents(ctx context.Context, instanceID string, pageSize int, pageToken string) ([]*ProcessEvent, string, error) {
	return uc.eventRepo.ListByInstance(ctx, instanceID, pageSize, pageToken)
}

func (uc *InstanceUseCase) GetHeatmap(ctx context.Context, definitionID string) ([]*HeatmapStep, error) {
	return uc.eventRepo.GetHeatmap(ctx, definitionID)
}
