package biz

import (
	"context"
	"time"
)

type Membership struct {
	ID        string
	UserID    string
	TenantID  string
	Role      string
	CreatedAt time.Time
}

type MembershipRepo interface {
	Create(ctx context.Context, m *Membership) (*Membership, error)
	ListByTenant(ctx context.Context, tenantID string, pageSize int, cursor string) ([]*Membership, string, error)
	ListByUser(ctx context.Context, userID string) ([]*Membership, error)
	GetByUserAndTenant(ctx context.Context, userID, tenantID string) (*Membership, error)
	SoftDelete(ctx context.Context, userID, tenantID string) error
}

type MembershipUsecase struct {
	repo MembershipRepo
}

func NewMembershipUsecase(repo MembershipRepo) *MembershipUsecase {
	return &MembershipUsecase{repo: repo}
}

func (uc *MembershipUsecase) ListByUser(ctx context.Context, userID string) ([]*Membership, error) {
	return uc.repo.ListByUser(ctx, userID)
}

func (uc *MembershipUsecase) AddMembership(ctx context.Context, userID, tenantID, role string) (*Membership, error) {
	return uc.repo.Create(ctx, &Membership{UserID: userID, TenantID: tenantID, Role: role})
}

func (uc *MembershipUsecase) InviteMember(ctx context.Context, tenantID, userID, role string) (*Membership, error) {
	existing, err := uc.repo.GetByUserAndTenant(ctx, userID, tenantID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}
	return uc.repo.Create(ctx, &Membership{UserID: userID, TenantID: tenantID, Role: role})
}

func (uc *MembershipUsecase) ListMembers(ctx context.Context, tenantID string, pageSize int, cursor string) ([]*Membership, string, error) {
	return uc.repo.ListByTenant(ctx, tenantID, pageSize, cursor)
}

func (uc *MembershipUsecase) RemoveMember(ctx context.Context, userID, tenantID string) error {
	return uc.repo.SoftDelete(ctx, userID, tenantID)
}
