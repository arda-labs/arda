package service

import (
	"context"
	"encoding/json"

	bpmv1 "github.com/arda-labs/arda/arda-be-go/services/bpm-service/api/bpm/v1"
	"github.com/arda-labs/arda/arda-be-go/services/bpm-service/internal/biz"

	kratoserrors "github.com/go-kratos/kratos/v2/errors"

	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func badRequest(format string) error {
	return kratoserrors.BadRequest("BAD_REQUEST", format)
}

// BPMService implements bpmv1.BPMServiceServer.
type BPMService struct {
	bpmv1.UnimplementedBPMServiceServer

	defUC        *biz.DefinitionUseCase
	instUC       *biz.InstanceUseCase
	eventUC      *biz.EventUseCase
	tplUC        *biz.TemplateUseCase
	deployBPMN   func(ctx context.Context, xml, name string) (int64, error)
}

func NewBPMService(
	defUC *biz.DefinitionUseCase,
	instUC *biz.InstanceUseCase,
	eventUC *biz.EventUseCase,
	tplUC *biz.TemplateUseCase,
) *BPMService {
	return &BPMService{
		defUC:   defUC,
		instUC:  instUC,
		eventUC: eventUC,
		tplUC:   tplUC,
	}
}

// SetDeployer injects the Zeebe deploy function after Wire construction.
func (s *BPMService) SetDeployer(fn func(ctx context.Context, xml, name string) (int64, error)) {
	s.deployBPMN = fn
}

// UseCases returns the underlying use cases for worker integration.
func (s *BPMService) UseCases() (defUC *biz.DefinitionUseCase, instUC *biz.InstanceUseCase, eventUC *biz.EventUseCase, tplUC *biz.TemplateUseCase) {
	return s.defUC, s.instUC, s.eventUC, s.tplUC
}

// ───────────────────────────────────────────────
// Process Definitions
// ───────────────────────────────────────────────

func (s *BPMService) ListDefinitions(ctx context.Context, req *bpmv1.ListDefinitionsRequest) (*bpmv1.ListDefinitionsResponse, error) {
	filter := biz.DefinitionFilter{
		Module:         req.GetModule(),
		Keyword:        req.GetKeyword(),
		IncludeInactive: req.GetIncludeInactive(),
		PageSize:       int(req.GetPageSize()),
		PageToken:      req.GetPageToken(),
	}
	defs, nextToken, err := s.defUC.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	items := make([]*bpmv1.Definition, 0, len(defs))
	for _, d := range defs {
		items = append(items, definitionToProto(d))
	}
	return &bpmv1.ListDefinitionsResponse{
		Definitions:   items,
		NextPageToken: nextToken,
	}, nil
}

func (s *BPMService) GetDefinition(ctx context.Context, req *bpmv1.GetDefinitionRequest) (*bpmv1.Definition, error) {
	d, err := s.defUC.GetByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return definitionToProto(d), nil
}

func (s *BPMService) GetDefinitionDiagram(ctx context.Context, req *bpmv1.GetDefinitionDiagramRequest) (*bpmv1.DefinitionDiagram, error) {
	d, err := s.defUC.GetDiagram(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &bpmv1.DefinitionDiagram{
		Id:         d.ID,
		BpmnXml:    d.BPMNXml,
		ProcessKey: d.ProcessKey,
		Version:    d.Version,
	}, nil
}

func (s *BPMService) DeployDefinition(ctx context.Context, req *bpmv1.DeployDefinitionRequest) (*bpmv1.Definition, error) {
	if req.GetBpmnXml() == "" {
		return nil, badRequest("bpmn_xml is required")
	}
	if req.GetProcessKey() == "" {
		return nil, badRequest("process_key is required")
	}
	if req.GetName() == "" {
		return nil, badRequest("name is required")
	}

	def := &biz.ProcessDefinition{
		ProcessKey:  req.GetProcessKey(),
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Category:    req.GetCategory(),
		Module:      req.GetModule(),
		BPMNXml:     req.GetBpmnXml(),
	}

	var deployFn func(string, string) (int64, error)
	if s.deployBPMN != nil {
		deployFn = func(xml, name string) (int64, error) {
			return s.deployBPMN(ctx, xml, name)
		}
	}

	d, err := s.defUC.Deploy(ctx, def, deployFn)
	if err != nil {
		return nil, err
	}
	return definitionToProto(d), nil
}

// ───────────────────────────────────────────────
// Process Instances
// ───────────────────────────────────────────────

func (s *BPMService) ListInstances(ctx context.Context, req *bpmv1.ListInstancesRequest) (*bpmv1.ListInstancesResponse, error) {
	filter := biz.InstanceFilter{
		ProcessDefinitionID: req.GetProcessDefinitionId(),
		Status:             req.GetStatus(),
		Module:             req.GetModule(),
		Keyword:            req.GetKeyword(),
		FromDate:           req.GetFromDate(),
		ToDate:             req.GetToDate(),
		PageSize:           int(req.GetPageSize()),
		PageToken:          req.GetPageToken(),
	}
	instances, nextToken, err := s.instUC.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	items := make([]*bpmv1.InstanceSummary, 0, len(instances))
	for _, inst := range instances {
		items = append(items, instanceToSummary(inst))
	}
	return &bpmv1.ListInstancesResponse{
		Instances:     items,
		NextPageToken: nextToken,
	}, nil
}

func (s *BPMService) GetInstance(ctx context.Context, req *bpmv1.GetInstanceRequest) (*bpmv1.InstanceDetail, error) {
	detail, err := s.instUC.GetDetail(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	var variables *structpb.Struct
	if detail.Variables != "" {
		var v map[string]interface{}
		if err := json.Unmarshal([]byte(detail.Variables), &v); err == nil {
			variables, _ = structpb.NewStruct(v)
		}
	}

	resp := &bpmv1.InstanceDetail{
		Id:                 detail.ID,
		ZeebeInstanceKey:   detail.ZeebeInstanceKey,
		ProcessDefinitionId: detail.ProcessDefinitionID,
		ProcessName:        detail.ProcessName,
		Status:             detail.Status,
		CurrentStep:        detail.CurrentStep,
		Variables:          variables,
		AssignedAgent:      detail.AssignedAgent,
		SlaStatus:          detail.SLAStatus,
		BpmnXml:            detail.BPMNXml,
		ActiveElementIds:   detail.ActiveElementIDs,
		CompletedElementIds: detail.CompletedElementIDs,
		CreatedAt:          timestamppb.New(detail.CreatedAt),
	}
	if detail.CompletedAt != nil {
		resp.CompletedAt = timestamppb.New(*detail.CompletedAt)
	}
	return resp, nil
}

func (s *BPMService) GetInstanceEvents(ctx context.Context, req *bpmv1.GetInstanceEventsRequest) (*bpmv1.GetInstanceEventsResponse, error) {
	events, nextToken, err := s.instUC.GetEvents(ctx, req.GetId(), int(req.GetPageSize()), req.GetPageToken())
	if err != nil {
		return nil, err
	}
	items := make([]*bpmv1.ProcessEvent, 0, len(events))
	for _, e := range events {
		items = append(items, eventToProto(e))
	}
	return &bpmv1.GetInstanceEventsResponse{
		Events:        items,
		NextPageToken: nextToken,
	}, nil
}

func (s *BPMService) GetHeatmap(ctx context.Context, req *bpmv1.GetHeatmapRequest) (*bpmv1.GetHeatmapResponse, error) {
	steps, err := s.instUC.GetHeatmap(ctx, req.GetProcessDefinitionId())
	if err != nil {
		return nil, err
	}
	items := make([]*bpmv1.HeatmapStep, 0, len(steps))
	for _, st := range steps {
		items = append(items, &bpmv1.HeatmapStep{
			ElementId:         st.ElementID,
			InstanceCount:     st.InstanceCount,
			AvgDurationSeconds: st.AvgDurationSeconds,
			Severity:          st.Severity,
		})
	}
	return &bpmv1.GetHeatmapResponse{Steps: items}, nil
}

// ───────────────────────────────────────────────
// Templates
// ───────────────────────────────────────────────

func (s *BPMService) ListTemplates(ctx context.Context, req *bpmv1.ListTemplatesRequest) (*bpmv1.ListTemplatesResponse, error) {
	filter := biz.TemplateFilter{
		ProcessDefinitionID: req.GetProcessDefinitionId(),
		Module:             req.GetModule(),
		Keyword:            req.GetKeyword(),
		PageSize:           int(req.GetPageSize()),
		PageToken:          req.GetPageToken(),
	}
	templates, nextToken, err := s.tplUC.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	items := make([]*bpmv1.Template, 0, len(templates))
	for _, t := range templates {
		items = append(items, templateToProto(t))
	}
	return &bpmv1.ListTemplatesResponse{
		Templates:     items,
		NextPageToken: nextToken,
	}, nil
}

func (s *BPMService) CreateTemplate(ctx context.Context, req *bpmv1.CreateTemplateRequest) (*bpmv1.Template, error) {
	tpl := &biz.Template{
		ProcessDefinitionID: req.GetProcessDefinitionId(),
		Name:               req.GetName(),
		TemplateText:       req.GetTemplateText(),
		Module:             req.GetModule(),
	}
	for _, v := range req.GetVariables() {
		tpl.Variables = append(tpl.Variables, variableMappingToBiz(v))
	}
	created, err := s.tplUC.Create(ctx, tpl)
	if err != nil {
		return nil, err
	}
	return templateToProto(created), nil
}

func (s *BPMService) UpdateTemplate(ctx context.Context, req *bpmv1.UpdateTemplateRequest) (*bpmv1.Template, error) {
	tpl := &biz.Template{
		ID:           req.GetId(),
		Name:         req.GetName(),
		TemplateText: req.GetTemplateText(),
		Module:       req.GetModule(),
	}
	for _, v := range req.GetVariables() {
		tpl.Variables = append(tpl.Variables, variableMappingToBiz(v))
	}
	updated, err := s.tplUC.Update(ctx, tpl)
	if err != nil {
		return nil, err
	}
	return templateToProto(updated), nil
}

func (s *BPMService) DeleteTemplate(ctx context.Context, req *bpmv1.DeleteTemplateRequest) (*bpmv1.DeleteResponse, error) {
	if err := s.tplUC.Delete(ctx, req.GetId()); err != nil {
		return nil, err
	}
	return &bpmv1.DeleteResponse{}, nil
}

func (s *BPMService) RenderTemplate(ctx context.Context, req *bpmv1.RenderTemplateRequest) (*bpmv1.RenderTemplateResponse, error) {
	result, err := s.tplUC.Render(ctx, req.GetId(), req.GetInstanceId())
	if err != nil {
		return nil, err
	}
	resp := &bpmv1.RenderTemplateResponse{
		RenderedText: result.RenderedText,
	}
	for _, rv := range result.ResolvedVariables {
		resp.ResolvedVariables = append(resp.ResolvedVariables, &bpmv1.ResolvedVariable{
			VariableName:  rv.VariableName,
			ResolvedValue: rv.ResolvedValue,
			UsedFallback:  rv.UsedFallback,
			SourceType:    rv.SourceType,
		})
	}
	return resp, nil
}

func (s *BPMService) ListVariableSources(ctx context.Context, req *bpmv1.ListVariableSourcesRequest) (*bpmv1.ListVariableSourcesResponse, error) {
	sources, err := s.tplUC.ListVariableSources(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]*bpmv1.VariableSource, 0, len(sources))
	for _, src := range sources {
		items = append(items, &bpmv1.VariableSource{
			Type:        src.Type,
			Name:        src.Name,
			Description: src.Description,
		})
	}
	return &bpmv1.ListVariableSourcesResponse{Sources: items}, nil
}

func (s *BPMService) ConfigureVariableMapping(ctx context.Context, req *bpmv1.ConfigureVariableMappingRequest) (*bpmv1.VariableMapping, error) {
	mapping := &biz.TemplateVariable{
		TemplateID:    req.GetTemplateId(),
		VariableName:  req.GetVariableName(),
		SourceType:    req.GetSourceType(),
		SourceField:   req.GetSourceField(),
		FallbackValue: req.GetFallbackValue(),
	}
	if req.GetResolverConfig() != nil {
		b, _ := req.GetResolverConfig().MarshalJSON()
		mapping.ResolverConfig = string(b)
	}
	created, err := s.tplUC.ConfigureVariableMapping(ctx, mapping)
	if err != nil {
		return nil, err
	}
	return &bpmv1.VariableMapping{
		Id:            created.ID,
		TemplateId:    created.TemplateID,
		VariableName:  created.VariableName,
		SourceType:    created.SourceType,
		SourceField:   created.SourceField,
		FallbackValue: created.FallbackValue,
	}, nil
}

// ───────────────────────────────────────────────
// Mappers
// ───────────────────────────────────────────────

func definitionToProto(d *biz.ProcessDefinition) *bpmv1.Definition {
	return &bpmv1.Definition{
		Id:          d.ID,
		ProcessKey:  d.ProcessKey,
		Name:        d.Name,
		Description: d.Description,
		Category:    d.Category,
		Module:      d.Module,
		Version:     d.Version,
		IsActive:    d.IsActive,
		CreatedAt:   timestamppb.New(d.CreatedAt),
		UpdatedAt:   timestamppb.New(d.UpdatedAt),
	}
}

func instanceToSummary(inst *biz.ProcessInstance) *bpmv1.InstanceSummary {
	s := &bpmv1.InstanceSummary{
		Id:                 inst.ID,
		ZeebeInstanceKey:   inst.ZeebeInstanceKey,
		ProcessDefinitionId: inst.ProcessDefinitionID,
		Status:             inst.Status,
		CurrentStep:        inst.CurrentStep,
		AssignedAgent:      inst.AssignedAgent,
		SlaStatus:          inst.SLAStatus,
		CreatedAt:          timestamppb.New(inst.CreatedAt),
	}
	if inst.CompletedAt != nil {
		s.CompletedAt = timestamppb.New(*inst.CompletedAt)
	}
	return s
}

func eventToProto(e *biz.ProcessEvent) *bpmv1.ProcessEvent {
	pe := &bpmv1.ProcessEvent{
		Id:                e.ID,
		ProcessInstanceId: e.ProcessInstanceID,
		EventType:         e.EventType,
		Source:            e.Source,
		Timestamp:         timestamppb.New(e.Timestamp),
	}
	if e.Data != "" {
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(e.Data), &data); err == nil {
			pe.Data, _ = structpb.NewStruct(data)
		}
	}
	return pe
}

func templateToProto(t *biz.Template) *bpmv1.Template {
	pt := &bpmv1.Template{
		Id:                  t.ID,
		ProcessDefinitionId: t.ProcessDefinitionID,
		Name:                t.Name,
		TemplateText:        t.TemplateText,
		Module:              t.Module,
		CreatedAt:           timestamppb.New(t.CreatedAt),
		UpdatedAt:           timestamppb.New(t.UpdatedAt),
	}
	for _, v := range t.Variables {
		pt.Variables = append(pt.Variables, &bpmv1.TemplateVariable{
			Id:            v.ID,
			VariableName:  v.VariableName,
			SourceType:    v.SourceType,
			SourceField:   v.SourceField,
			FallbackValue: v.FallbackValue,
		})
	}
	return pt
}

func variableMappingToBiz(v *bpmv1.TemplateVariable) *biz.TemplateVariable {
	return &biz.TemplateVariable{
		VariableName:  v.GetVariableName(),
		SourceType:    v.GetSourceType(),
		SourceField:   v.GetSourceField(),
		FallbackValue: v.GetFallbackValue(),
	}
}
