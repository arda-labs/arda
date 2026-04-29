package biz

import (
	"context"
	"fmt"
	"strings"
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
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) (*User, error)
}

type UserUsecase struct {
	repo        UserRepo
	auth        *AuthUsecase
	tenantUsers *TenantUserUsecase
}

func NewUserUsecase(repo UserRepo, auth *AuthUsecase, tenantUsers *TenantUserUsecase) *UserUsecase {
	return &UserUsecase{repo: repo, auth: auth, tenantUsers: tenantUsers}
}

func (uc *UserUsecase) CreateUser(ctx context.Context, username, email, displayName, password, tenantID string) (*TenantUser, error) {
	if tenantID == "" {
		return nil, fmt.Errorf("tenant_id is required")
	}
	if email == "" {
		return nil, fmt.Errorf("email is required")
	}
	if !uc.auth.HasZitadelManagementPAT() {
		return nil, fmt.Errorf("zitadel management PAT is required to create login users; set ZITADEL_MANAGEMENT_PAT or ZITADEL_PAT")
	}

	user, err := uc.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		identityUsername := email
		externalID, err := uc.auth.CreateZitadelUser(ctx, identityUsername, email, displayName, password)
		if err != nil {
			return nil, err
		}

		user, err = uc.repo.Create(ctx, &User{
			ExternalID:  externalID,
			Email:       email,
			DisplayName: displayName,
		})
		if err != nil {
			return nil, err
		}
	} else if isLocalExternalID(user.ExternalID) {
		externalID, err := uc.auth.CreateZitadelUser(ctx, email, email, displayName, password)
		if err != nil {
			return nil, err
		}
		user.ExternalID = externalID
		if displayName != "" {
			user.DisplayName = displayName
		}
		user, err = uc.repo.Update(ctx, user)
		if err != nil {
			return nil, err
		}
	}

	if username == "" {
		username = defaultTenantUsername(email)
	}
	return uc.tenantUsers.AddTenantUser(ctx, user.ID, tenantID, username, displayName, "MEMBER")
}

func (uc *UserUsecase) GetOrCreateUser(ctx context.Context, externalID, email, displayName string) (*User, error) {
	user, err := uc.repo.GetByExternalID(ctx, externalID)
	if err != nil {
		return nil, err
	}
	if user != nil {
		changed := false
		if email != "" && user.Email != email {
			user.Email = email
			changed = true
		}
		if displayName != "" && user.DisplayName != displayName {
			user.DisplayName = displayName
			changed = true
		}
		if changed {
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

func (uc *UserUsecase) ListUsers(ctx context.Context, tenantID string, pageSize int, cursor string) ([]*TenantUser, string, error) {
	return uc.tenantUsers.ListTenantUsers(ctx, tenantID, pageSize, cursor)
}

func defaultTenantUsername(email string) string {
	local, _, ok := strings.Cut(email, "@")
	if !ok || local == "" {
		return email
	}
	return local
}

func isLocalExternalID(externalID string) bool {
	return strings.HasPrefix(strings.ToLower(strings.TrimSpace(externalID)), "local:")
}
