package biz

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewUserUsecase,
	NewTenantUsecase,
	NewTenantUserUsecase,
	NewRoleUsecase,
	NewPermissionUsecase,
	NewAuditUsecase,
	NewAuthUsecase,
	NewMenuUsecase,
	NewGroupUsecase,
)
