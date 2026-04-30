package biz

import (
	"context"
	stderrors "errors"
	"strings"
	"time"
)

var (
	ErrNotFound        = stderrors.New("notification record not found")
	ErrInvalidArgument = stderrors.New("invalid notification request")
)

type PageFilter struct {
	Status    string
	Keyword   string
	PageSize  int
	PageToken string
}

type TemplateVersionFilter struct {
	TemplateID string
	Channel    string
	Language   string
	Status     string
}

type DeliveryFilter struct {
	Status      string
	Channel     string
	RecipientID string
	PageSize    int
	PageToken   string
}

type InAppFilter struct {
	RecipientType string
	RecipientID   string
	Status        string
	PageSize      int
	PageToken     string
}

type NotificationRepo interface {
	ListTemplates(ctx context.Context, filter PageFilter) ([]*NotificationTemplate, string, error)
	GetTemplate(ctx context.Context, id string) (*NotificationTemplate, error)
	CreateTemplate(ctx context.Context, item *NotificationTemplate) (*NotificationTemplate, error)
	UpdateTemplate(ctx context.Context, item *NotificationTemplate) (*NotificationTemplate, error)
	DeleteTemplate(ctx context.Context, id string) error

	ListTemplateVersions(ctx context.Context, filter TemplateVersionFilter) ([]*NotificationTemplateVersion, error)
	CreateTemplateVersion(ctx context.Context, item *NotificationTemplateVersion) (*NotificationTemplateVersion, error)
	ApproveTemplateVersion(ctx context.Context, id, actor, note string) (*NotificationTemplateVersion, error)

	CreateNotificationRequest(ctx context.Context, item *NotificationRequest) (*NotificationRequest, error)
	GetNotificationRequest(ctx context.Context, id string) (*NotificationRequest, error)
	ListDeliveries(ctx context.Context, filter DeliveryFilter) ([]*NotificationDelivery, string, error)
	RetryDelivery(ctx context.Context, id, actor string) (*NotificationDelivery, error)
	ClaimDueDeliveries(ctx context.Context, workerID string, limit int) ([]*NotificationDelivery, error)
	CreateInAppNotification(ctx context.Context, delivery *NotificationDelivery) (*InAppNotification, error)
	MarkDeliveryDelivered(ctx context.Context, id, providerCode, providerMessageID, providerResponseJSON string) (*NotificationDelivery, error)
	MarkDeliveryFailed(ctx context.Context, id, message string, retryAfterSeconds int) (*NotificationDelivery, error)
	ListInAppNotifications(ctx context.Context, filter InAppFilter) ([]*InAppNotification, string, error)
	MarkInAppNotificationRead(ctx context.Context, id, actor string) (*InAppNotification, error)
}

type NotificationUsecase struct {
	repo NotificationRepo
}

func NewNotificationUsecase(repo NotificationRepo) *NotificationUsecase {
	return &NotificationUsecase{repo: repo}
}

type NotificationTemplate struct {
	ID             string
	Code           string
	Name           string
	Category       string
	DefaultChannel string
	Description    string
	Status         string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type NotificationTemplateVersion struct {
	ID                string
	TemplateID        string
	Version           int
	Channel           string
	Language          string
	Subject           string
	Body              string
	PayloadSchemaJSON string
	ApprovalStatus    string
	ApprovedBy        string
	ApprovedAt        time.Time
	ChangeNote        string
	Status            string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type NotificationRequest struct {
	ID               string
	SourceService    string
	EventType        string
	CorrelationID    string
	IdempotencyKey   string
	TemplateCode     string
	RecipientType    string
	RecipientID      string
	RecipientAddress string
	Channels         []string
	Language         string
	PayloadJSON      string
	Priority         int
	Status           string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Deliveries       []*NotificationDelivery
}

type NotificationDelivery struct {
	ID                   string
	RequestID            string
	TemplateVersionID    string
	Channel              string
	RecipientType        string
	RecipientID          string
	RecipientAddress     string
	Subject              string
	Body                 string
	Status               string
	AttemptCount         int
	MaxAttempts          int
	NextAttemptAt        time.Time
	LockedBy             string
	LockedAt             time.Time
	ProviderCode         string
	ProviderMessageID    string
	ProviderResponseJSON string
	ErrorMessage         string
	Priority             int
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type InAppNotification struct {
	ID            string
	DeliveryID    string
	RecipientType string
	RecipientID   string
	Title         string
	Body          string
	DataJSON      string
	Status        string
	ReadAt        time.Time
	CreatedAt     time.Time
}

func (uc *NotificationUsecase) ListTemplates(ctx context.Context, filter PageFilter) ([]*NotificationTemplate, string, error) {
	normalizePageFilter(&filter)
	return uc.repo.ListTemplates(ctx, filter)
}

func (uc *NotificationUsecase) GetTemplate(ctx context.Context, id string) (*NotificationTemplate, error) {
	return uc.repo.GetTemplate(ctx, strings.TrimSpace(id))
}

func (uc *NotificationUsecase) CreateTemplate(ctx context.Context, item *NotificationTemplate) (*NotificationTemplate, error) {
	normalizeTemplate(item)
	if item.Code == "" || item.Name == "" {
		return nil, ErrInvalidArgument
	}
	return uc.repo.CreateTemplate(ctx, item)
}

func (uc *NotificationUsecase) UpdateTemplate(ctx context.Context, item *NotificationTemplate) (*NotificationTemplate, error) {
	normalizeTemplate(item)
	if item.ID == "" || item.Code == "" || item.Name == "" {
		return nil, ErrInvalidArgument
	}
	return uc.repo.UpdateTemplate(ctx, item)
}

func (uc *NotificationUsecase) DeleteTemplate(ctx context.Context, id string) error {
	return uc.repo.DeleteTemplate(ctx, strings.TrimSpace(id))
}

func (uc *NotificationUsecase) ListTemplateVersions(ctx context.Context, filter TemplateVersionFilter) ([]*NotificationTemplateVersion, error) {
	filter.TemplateID = strings.TrimSpace(filter.TemplateID)
	filter.Channel = upperDefault(filter.Channel, "")
	filter.Language = lowerDefault(filter.Language, "")
	filter.Status = upperDefault(filter.Status, "")
	return uc.repo.ListTemplateVersions(ctx, filter)
}

func (uc *NotificationUsecase) CreateTemplateVersion(ctx context.Context, item *NotificationTemplateVersion) (*NotificationTemplateVersion, error) {
	normalizeTemplateVersion(item)
	if item.TemplateID == "" || item.Channel == "" || item.Language == "" || item.Body == "" {
		return nil, ErrInvalidArgument
	}
	return uc.repo.CreateTemplateVersion(ctx, item)
}

func (uc *NotificationUsecase) ApproveTemplateVersion(ctx context.Context, id, actor, note string) (*NotificationTemplateVersion, error) {
	actor = strings.TrimSpace(actor)
	if actor == "" {
		actor = "SYSTEM"
	}
	return uc.repo.ApproveTemplateVersion(ctx, strings.TrimSpace(id), actor, strings.TrimSpace(note))
}

func normalizePageFilter(filter *PageFilter) {
	filter.Status = upperDefault(filter.Status, "")
	filter.Keyword = strings.TrimSpace(filter.Keyword)
}

func normalizeTemplate(item *NotificationTemplate) {
	item.Code = upperDefault(item.Code, "")
	item.Name = strings.TrimSpace(item.Name)
	item.Category = upperDefault(item.Category, "")
	item.DefaultChannel = upperDefault(item.DefaultChannel, "IN_APP")
	item.Description = strings.TrimSpace(item.Description)
	item.Status = upperDefault(item.Status, "ACTIVE")
}

func normalizeTemplateVersion(item *NotificationTemplateVersion) {
	item.TemplateID = strings.TrimSpace(item.TemplateID)
	item.Channel = upperDefault(item.Channel, "IN_APP")
	item.Language = lowerDefault(item.Language, "vi")
	item.Subject = strings.TrimSpace(item.Subject)
	item.Body = strings.TrimSpace(item.Body)
	item.PayloadSchemaJSON = strings.TrimSpace(item.PayloadSchemaJSON)
	if item.PayloadSchemaJSON == "" {
		item.PayloadSchemaJSON = "{}"
	}
	item.ApprovalStatus = upperDefault(item.ApprovalStatus, "DRAFT")
	item.ApprovedBy = strings.TrimSpace(item.ApprovedBy)
	item.ChangeNote = strings.TrimSpace(item.ChangeNote)
	item.Status = upperDefault(item.Status, "ACTIVE")
	if item.Version <= 0 {
		item.Version = 1
	}
}

func upperDefault(value, fallback string) string {
	value = strings.ToUpper(strings.TrimSpace(value))
	if value == "" {
		return fallback
	}
	return value
}

func lowerDefault(value, fallback string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "" {
		return fallback
	}
	return value
}
