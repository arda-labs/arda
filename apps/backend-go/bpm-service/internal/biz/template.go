package biz

import (
	"context"
	"regexp"
	"strings"
	"sync"
	"time"
)

var variablePattern = regexp.MustCompile(`@([A-Z_][A-Z0-9_]*)`)

type Template struct {
	ID                  string
	ProcessDefinitionID string
	Name                string
	TemplateText        string
	Module              string
	Variables           []*TemplateVariable
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type TemplateVariable struct {
	ID            string
	TemplateID    string
	VariableName  string
	SourceType    string
	SourceField   string
	ResolverConfig string // JSON string
	FallbackValue string
}

type TemplateFilter struct {
	ProcessDefinitionID string
	Module              string
	Keyword             string
	PageSize            int
	PageToken           string
}

type VariableSource struct {
	Type        string
	Name        string
	Description string
}

type ResolvedVariable struct {
	VariableName  string
	ResolvedValue string
	UsedFallback  bool
	SourceType    string
}

type RenderResult struct {
	RenderedText     string
	ResolvedVariables []*ResolvedVariable
}

// TemplateRepo interface for template persistence.
type TemplateRepo interface {
	List(ctx context.Context, filter TemplateFilter) ([]*Template, string, error)
	GetByID(ctx context.Context, id string) (*Template, error)
	Create(ctx context.Context, tpl *Template) (*Template, error)
	Update(ctx context.Context, tpl *Template) (*Template, error)
	Delete(ctx context.Context, id string) error
	ListVariableSources(ctx context.Context) ([]*VariableSource, error)
	CreateVariableMapping(ctx context.Context, mapping *TemplateVariable) (*TemplateVariable, error)
}

// VariableResolver resolves a single variable based on its source configuration.
type VariableResolver interface {
	Resolve(ctx context.Context, variable *TemplateVariable, instanceID string) (string, error)
}

// TemplateUseCase handles business logic for templates.
type TemplateUseCase struct {
	repo     TemplateRepo
	resolver VariableResolver
	renderCache *RenderCache
}

func NewTemplateUseCase(repo TemplateRepo, resolver VariableResolver) *TemplateUseCase {
	return &TemplateUseCase{
		repo:        repo,
		resolver:    resolver,
		renderCache: NewRenderCache(5 * time.Minute),
	}
}

func (uc *TemplateUseCase) List(ctx context.Context, filter TemplateFilter) ([]*Template, string, error) {
	return uc.repo.List(ctx, filter)
}

func (uc *TemplateUseCase) GetByID(ctx context.Context, id string) (*Template, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *TemplateUseCase) Create(ctx context.Context, tpl *Template) (*Template, error) {
	return uc.repo.Create(ctx, tpl)
}

func (uc *TemplateUseCase) Update(ctx context.Context, tpl *Template) (*Template, error) {
	return uc.repo.Update(ctx, tpl)
}

func (uc *TemplateUseCase) Delete(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *TemplateUseCase) ListVariableSources(ctx context.Context) ([]*VariableSource, error) {
	return uc.repo.ListVariableSources(ctx)
}

func (uc *TemplateUseCase) ConfigureVariableMapping(ctx context.Context, mapping *TemplateVariable) (*TemplateVariable, error) {
	return uc.repo.CreateVariableMapping(ctx, mapping)
}

// Render resolves all variables in a template and returns the rendered text.
func (uc *TemplateUseCase) Render(ctx context.Context, templateID, instanceID string) (*RenderResult, error) {
	// Check cache
	cacheKey := templateID + ":" + instanceID
	if cached, ok := uc.renderCache.Get(cacheKey); ok {
		return cached, nil
	}

	tpl, err := uc.repo.GetByID(ctx, templateID)
	if err != nil {
		return nil, err
	}

	result := &RenderResult{
		RenderedText: tpl.TemplateText,
	}

	matches := variablePattern.FindAllStringSubmatch(tpl.TemplateText, -1)
	seen := make(map[string]bool)

	for _, match := range matches {
		varName := match[1]
		if seen[varName] {
			continue
		}
		seen[varName] = true

		resolved := &ResolvedVariable{
			VariableName: varName,
		}

		// Find variable mapping
		var mapping *TemplateVariable
		for _, v := range tpl.Variables {
			if v.VariableName == varName {
				mapping = v
				break
			}
		}

		if mapping != nil {
			resolved.SourceType = mapping.SourceType
			value, err := uc.resolver.Resolve(ctx, mapping, instanceID)
			if err != nil || value == "" {
				resolved.ResolvedValue = mapping.FallbackValue
				resolved.UsedFallback = true
			} else {
				resolved.ResolvedValue = value
			}
		} else {
			// No mapping configured, use fallback
			resolved.SourceType = "UNKNOWN"
			resolved.ResolvedValue = "@" + varName
			resolved.UsedFallback = true
		}

		result.RenderedText = strings.ReplaceAll(result.RenderedText, "@"+varName, resolved.ResolvedValue)
		result.ResolvedVariables = append(result.ResolvedVariables, resolved)
	}

	// Cache result
	uc.renderCache.Set(cacheKey, result)
	return result, nil
}

// RenderCache with TTL eviction.
type RenderCache struct {
	mu    sync.RWMutex
	items map[string]*cacheEntry
	ttl   time.Duration
}

type cacheEntry struct {
	result   *RenderResult
	expiresAt time.Time
}

func NewRenderCache(ttl time.Duration) *RenderCache {
	c := &RenderCache{
		items: make(map[string]*cacheEntry),
		ttl:   ttl,
	}
	// Periodic cleanup
	go func() {
		ticker := time.NewTicker(ttl)
		defer ticker.Stop()
		for range ticker.C {
			c.cleanup()
		}
	}()
	return c
}

func (c *RenderCache) Get(key string) (*RenderResult, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.items[key]
	if !ok || time.Now().After(entry.expiresAt) {
		return nil, false
	}
	return entry.result, true
}

func (c *RenderCache) Set(key string, result *RenderResult) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = &cacheEntry{
		result:    result,
		expiresAt: time.Now().Add(c.ttl),
	}
}

func (c *RenderCache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	for key, entry := range c.items {
		if now.After(entry.expiresAt) {
			delete(c.items, key)
		}
	}
}
