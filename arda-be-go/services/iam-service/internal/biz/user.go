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
	repo UserRepo
}

func NewUserUsecase(repo UserRepo) *UserUsecase {
	return &UserUsecase{repo: repo}
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

func (uc *UserUsecase) ListUsers(ctx context.Context, tenantID string, pageSize int, cursor string) ([]*User, string, error) {
	return uc.repo.ListByTenant(ctx, tenantID, pageSize, cursor)
}
