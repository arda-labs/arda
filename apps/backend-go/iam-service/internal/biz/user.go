package biz

import (
	"context"
	"time"
)

type User struct {
	ID          string
	ExternalID  string
	Email       string
	DisplayName string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UserRepo interface {
	Create(ctx context.Context, user *User) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	GetByExternalID(ctx context.Context, externalID string) (*User, error)
	Update(ctx context.Context, user *User) (*User, error)
	ListByTenant(ctx context.Context, tenantID string, pageSize int, cursor string) ([]*User, string, error)
}

type UserUsecase struct {
	repo    UserRepo
	auth    *AuthUsecase
	members *MembershipUsecase
}

func NewUserUsecase(repo UserRepo, auth *AuthUsecase, members *MembershipUsecase) *UserUsecase {
	return &UserUsecase{repo: repo, auth: auth, members: members}
}

func (uc *UserUsecase) CreateUser(ctx context.Context, email, displayName, password, tenantID string) (*User, error) {
	// 1. Gọi Zitadel để tạo Human User
	externalID, err := uc.auth.CreateZitadelUser(ctx, email, displayName, password)
	if err != nil {
		return nil, err
	}

	// 2. Lưu vào DB nội bộ Arda
	user, err := uc.repo.Create(ctx, &User{
		ExternalID:  externalID,
		Email:       email,
		DisplayName: displayName,
	})
	if err != nil {
		return nil, err
	}

	// 3. Tự động thêm vào Tenant hiện tại với Role mặc định là MEMBER
	if tenantID != "" {
		_, _ = uc.members.AddMembership(ctx, user.ID, tenantID, "MEMBER")
	}

	return user, nil
}

func (uc *UserUsecase) GetOrCreateUser(ctx context.Context, externalID, email, displayName string) (*User, error) {
	user, err := uc.repo.GetByExternalID(ctx, externalID)
	if err != nil {
		return nil, err
	}
	if user != nil {
		if user.Email != email || user.DisplayName != displayName {
			user.Email = email
			user.DisplayName = displayName
			return uc.repo.Update(ctx, user)
		}
		return user, nil
	}
	return uc.repo.Create(ctx, &User{
		ExternalID:  externalID,
		Email:       email,
		DisplayName: displayName,
	})
}

func (uc *UserUsecase) GetUser(ctx context.Context, id string) (*User, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *UserUsecase) UpdateProfile(ctx context.Context, id, displayName string) (*User, error) {
	user, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	user.DisplayName = displayName
	return uc.repo.Update(ctx, user)
}

func (uc *UserUsecase) ListUsers(ctx context.Context, tenantID string, pageSize int, cursor string) ([]*User, string, error) {
	return uc.repo.ListByTenant(ctx, tenantID, pageSize, cursor)
}
