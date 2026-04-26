package service

import (
	"context"

	"github.com/arda-labs/arda/arda-be-go/pkg/middleware"
	pb "github.com/arda-labs/arda/arda-be-go/services/iam-service/api/iam/v1"
	"github.com/arda-labs/arda/arda-be-go/services/iam-service/internal/biz"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IAMService struct {
	pb.UnimplementedIAMServiceServer

	users   *biz.UserUsecase
	tenants *biz.TenantUsecase
	members *biz.MembershipUsecase
	roles   *biz.RoleUsecase
	perms   *biz.PermissionUsecase
	auth    *biz.AuthUsecase
	audit   *biz.AuditUsecase
	groups  *biz.GroupUsecase
	log     *log.Helper
}

func NewIAMService(
	users *biz.UserUsecase,
	tenants *biz.TenantUsecase,
	members *biz.MembershipUsecase,
	roles *biz.RoleUsecase,
	perms *biz.PermissionUsecase,
	auth *biz.AuthUsecase,
	audit *biz.AuditUsecase,
	groups *biz.GroupUsecase,
	logger log.Logger,
) *IAMService {
	return &IAMService{
		users:   users,
		tenants: tenants,
		members: members,
		roles:   roles,
		perms:   perms,
		auth:    auth,
		audit:   audit,
		groups:  groups,
		log:     log.NewHelper(logger),
	}
}

// Auth

func (s *IAMService) CustomLogin(ctx context.Context, req *pb.CustomLoginRequest) (*pb.CustomLoginReply, error) {
	callbackURL, err := s.auth.CustomLogin(ctx, req.Email, req.Password, req.AuthRequestId)
	if err != nil {
		return nil, err
	}
	return &pb.CustomLoginReply{CallbackUrl: callbackURL}, nil
}

func (s *IAMService) GetUserMemberships(ctx context.Context, req *pb.GetUserMembershipsRequest) (*pb.GetUserMembershipsResponse, error) {
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		return nil, errors.Forbidden("UNAUTHORIZED", "missing subject")
	}

	user, err := s.users.GetOrCreateUser(ctx, userID, middleware.GetEmail(ctx), "")
	if err != nil {
		return nil, err
	}

	memberships, err := s.members.ListByUser(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	var tenantMemberships []*pb.TenantMembership
	if len(memberships) > 0 {
		tenantIDs := make([]string, len(memberships))
		for i, m := range memberships {
			tenantIDs[i] = m.TenantID
		}
		tenants, err := s.tenants.GetTenantsByIDs(ctx, tenantIDs)
		if err != nil {
			return nil, err
		}

		tenantMap := make(map[string]*biz.Tenant)
		for _, t := range tenants {
			tenantMap[t.ID] = t
		}

		for _, m := range memberships {
			t := tenantMap[m.TenantID]
			if t == nil {
				continue
			}
			tenantMemberships = append(tenantMemberships, &pb.TenantMembership{
				TenantId:   t.ID,
				TenantName: t.Name,
				TenantSlug: t.Slug,
				Role:       m.Role,
			})
		}
	}

	return &pb.GetUserMembershipsResponse{Memberships: tenantMemberships}, nil
}

func (s *IAMService) GetCurrentUserPermissions(ctx context.Context, req *pb.GetCurrentUserPermissionsRequest) (*pb.ListPermissionsResponse, error) {
	extID := middleware.GetUserID(ctx)
	tenantID := middleware.GetTenantID(ctx)
	if extID == "" || tenantID == "" {
		return nil, errors.Forbidden("UNAUTHORIZED", "missing subject or tenant")
	}

	user, err := s.users.GetOrCreateUser(ctx, extID, "", "")
	if err != nil {
		return nil, err
	}

	perms, err := s.perms.GetUserPermissions(ctx, user.ID, tenantID)
	if err != nil {
		return nil, err
	}

	resp := &pb.ListPermissionsResponse{}
	for _, p := range perms {
		resp.Permissions = append(resp.Permissions, p.Resource+":"+p.Action)
	}
	return resp, nil
}

func (s *IAMService) ForwardAuth(ctx context.Context, req *pb.ForwardAuthRequest) (*pb.ForwardAuthResponse, error) {
	token := req.Token
	if token == "" {
		if tr, ok := transport.FromServerContext(ctx); ok {
			token = tr.RequestHeader().Get("Authorization")
		}
	}

	allowed, userID, tenantID, err := s.auth.ForwardAuth(ctx, req.Method, req.Path, token)
	if err != nil {
		return nil, err
	}
	return &pb.ForwardAuthResponse{
		Allowed:  allowed,
		UserId:   userID,
		TenantId: tenantID,
	}, nil
}

// Users

func (s *IAMService) GetCurrentUser(ctx context.Context, _ *pb.GetCurrentUserRequest) (*pb.GetCurrentUserResponse, error) {
	extID := middleware.GetUserID(ctx)
	if extID == "" {
		return nil, errors.Forbidden("UNAUTHORIZED", "missing subject")
	}

	user, err := s.users.GetOrCreateUser(ctx, extID, middleware.GetEmail(ctx), "")
	if err != nil {
		return nil, err
	}

	memberships, err := s.members.ListByUser(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	var tenantMemberships []*pb.TenantMembership
	if len(memberships) > 0 {
		tenantIDs := make([]string, len(memberships))
		for i, m := range memberships {
			tenantIDs[i] = m.TenantID
		}
		tenants, err := s.tenants.GetTenantsByIDs(ctx, tenantIDs)
		if err != nil {
			return nil, err
		}

		tenantMap := make(map[string]*biz.Tenant)
		for _, t := range tenants {
			tenantMap[t.ID] = t
		}

		for _, m := range memberships {
			t := tenantMap[m.TenantID]
			if t == nil {
				continue
			}
			tenantMemberships = append(tenantMemberships, &pb.TenantMembership{
				TenantId:   t.ID,
				TenantName: t.Name,
				TenantSlug: t.Slug,
				Role:       m.Role,
			})
		}
	}

	return &pb.GetCurrentUserResponse{
		User:        toProtoUser(user),
		Memberships: tenantMemberships,
	}, nil
}

func (s *IAMService) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.User, error) {
	extID := middleware.GetUserID(ctx)
	if extID == "" {
		return nil, errors.Forbidden("UNAUTHORIZED", "missing subject")
	}

	user, err := s.users.GetOrCreateUser(ctx, extID, "", "")
	if err != nil {
		return nil, err
	}

	updated, err := s.users.UpdateProfile(ctx, user.ID, req.DisplayName)
	if err != nil {
		return nil, err
	}

	return toProtoUser(updated), nil
}

func (s *IAMService) ListMyAuditLogs(ctx context.Context, req *pb.ListMyAuditLogsRequest) (*pb.ListAuditLogsResponse, error) {
	extID := middleware.GetUserID(ctx)
	if extID == "" {
		return nil, errors.Forbidden("UNAUTHORIZED", "missing subject")
	}

	user, err := s.users.GetOrCreateUser(ctx, extID, "", "")
	if err != nil {
		return nil, err
	}

	logs, next, err := s.audit.ListMyLogs(ctx, user.ID, int(req.PageSize), req.PageToken)
	if err != nil {
		return nil, err
	}

	resp := &pb.ListAuditLogsResponse{NextPageToken: next}
	for _, l := range logs {
		resp.Logs = append(resp.Logs, &pb.AuditLog{
			Id:         l.ID,
			ActorId:    l.ActorID,
			Action:     l.Action,
			TargetType: l.TargetType,
			TargetId:   l.TargetID,
			CreatedAt:  timestamppb.New(l.CreatedAt),
		})
	}
	return resp, nil
}

func (s *IAMService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	user, err := s.users.GetUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return toProtoUser(user), nil
}

func (s *IAMService) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	users, next, err := s.users.ListUsers(ctx, req.TenantId, int(req.PageSize), req.PageToken)
	if err != nil {
		return nil, err
	}
	resp := &pb.ListUsersResponse{NextPageToken: next}
	for _, u := range users {
		resp.Users = append(resp.Users, toProtoUser(u))
	}
	return resp, nil
}

func (s *IAMService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	user, err := s.users.CreateUser(ctx, req.Email, req.DisplayName, req.Password, req.TenantId)
	if err != nil {
		return nil, err
	}
	return toProtoUser(user), nil
}

// Tenants

func (s *IAMService) CreateTenant(ctx context.Context, req *pb.CreateTenantRequest) (*pb.Tenant, error) {
	sub := middleware.GetUserID(ctx)
	t, err := s.tenants.CreateTenant(ctx, req.Name, req.Slug, sub)
	if err != nil {
		return nil, err
	}
	return toProtoTenant(t), nil
}

func (s *IAMService) GetTenant(ctx context.Context, req *pb.GetTenantRequest) (*pb.Tenant, error) {
	t, err := s.tenants.GetTenant(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return toProtoTenant(t), nil
}

func (s *IAMService) UpdateTenant(ctx context.Context, req *pb.UpdateTenantRequest) (*pb.Tenant, error) {
	t, err := s.tenants.UpdateTenant(ctx, req.Id, req.Name, req.Slug)
	if err != nil {
		return nil, err
	}
	return toProtoTenant(t), nil
}

func (s *IAMService) DeleteTenant(ctx context.Context, req *pb.DeleteTenantRequest) (*pb.DeleteTenantResponse, error) {
	return &pb.DeleteTenantResponse{}, s.tenants.DeleteTenant(ctx, req.Id)
}

// Membership

func (s *IAMService) InviteMember(ctx context.Context, req *pb.InviteMemberRequest) (*pb.Membership, error) {
	m, err := s.members.InviteMember(ctx, req.TenantId, req.ExternalId, req.Role)
	if err != nil {
		return nil, err
	}
	return toProtoMembership(m), nil
}

func (s *IAMService) ListMembers(ctx context.Context, req *pb.ListMembersRequest) (*pb.ListMembersResponse, error) {
	list, next, err := s.members.ListMembers(ctx, req.TenantId, int(req.PageSize), req.PageToken)
	if err != nil {
		return nil, err
	}
	resp := &pb.ListMembersResponse{NextPageToken: next}
	for _, m := range list {
		resp.Memberships = append(resp.Memberships, toProtoMembership(m))
	}
	return resp, nil
}

func (s *IAMService) RemoveMember(ctx context.Context, req *pb.RemoveMemberRequest) (*pb.RemoveMemberResponse, error) {
	return &pb.RemoveMemberResponse{}, s.members.RemoveMember(ctx, req.UserId, req.TenantId)
}

// Roles

func (s *IAMService) CreateRole(ctx context.Context, req *pb.CreateRoleRequest) (*pb.Role, error) {
	role, err := s.roles.CreateRole(ctx, req.TenantId, req.Name, req.Description, req.Permissions)
	if err != nil {
		return nil, err
	}
	return toProtoRole(role), nil
}

func (s *IAMService) GetRole(ctx context.Context, req *pb.GetRoleRequest) (*pb.Role, error) {
	role, err := s.roles.GetRole(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return toProtoRole(role), nil
}

func (s *IAMService) ListRoles(ctx context.Context, req *pb.ListRolesRequest) (*pb.ListRolesResponse, error) {
	list, next, err := s.roles.ListRoles(ctx, req.TenantId, int(req.PageSize), req.PageToken)
	if err != nil {
		return nil, err
	}
	resp := &pb.ListRolesResponse{NextPageToken: next}
	for _, r := range list {
		resp.Roles = append(resp.Roles, toProtoRole(r))
	}
	return resp, nil
}

func (s *IAMService) UpdateRole(ctx context.Context, req *pb.UpdateRoleRequest) (*pb.Role, error) {
	role, err := s.roles.UpdateRole(ctx, req.Id, req.Name, req.Description, req.Permissions)
	if err != nil {
		return nil, err
	}
	return toProtoRole(role), nil
}

func (s *IAMService) DeleteRole(ctx context.Context, req *pb.DeleteRoleRequest) (*pb.DeleteRoleResponse, error) {
	return &pb.DeleteRoleResponse{}, s.roles.DeleteRole(ctx, req.Id)
}

// Role assignment

func (s *IAMService) AssignRole(ctx context.Context, req *pb.AssignRoleRequest) (*pb.UserRole, error) {
	actorID := middleware.GetUserID(ctx)
	if err := s.roles.AssignRole(ctx, req.UserId, req.RoleId, req.TenantId, actorID); err != nil {
		return nil, err
	}
	return &pb.UserRole{UserId: req.UserId, RoleId: req.RoleId, TenantId: req.TenantId}, nil
}

func (s *IAMService) RevokeRole(ctx context.Context, req *pb.RevokeRoleRequest) (*pb.RevokeRoleResponse, error) {
	actorID := middleware.GetUserID(ctx)
	return &pb.RevokeRoleResponse{}, s.roles.RevokeRole(ctx, req.UserId, req.RoleId, req.TenantId, actorID)
}

// Groups

func (s *IAMService) CreateGroup(ctx context.Context, req *pb.CreateGroupRequest) (*pb.Group, error) {
	g, err := s.groups.CreateGroup(ctx, req.TenantId, req.Name, req.Description)
	if err != nil {
		return nil, err
	}
	return toProtoGroup(g), nil
}

func (s *IAMService) GetGroup(ctx context.Context, req *pb.GetGroupRequest) (*pb.Group, error) {
	g, err := s.groups.GetGroup(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return toProtoGroup(g), nil
}

func (s *IAMService) UpdateGroup(ctx context.Context, req *pb.UpdateGroupRequest) (*pb.Group, error) {
	g, err := s.groups.UpdateGroup(ctx, req.Id, req.Name, req.Description)
	if err != nil {
		return nil, err
	}
	return toProtoGroup(g), nil
}

func (s *IAMService) DeleteGroup(ctx context.Context, req *pb.DeleteGroupRequest) (*pb.DeleteGroupResponse, error) {
	return &pb.DeleteGroupResponse{}, s.groups.DeleteGroup(ctx, req.Id)
}

func (s *IAMService) ListGroups(ctx context.Context, req *pb.ListGroupsRequest) (*pb.ListGroupsResponse, error) {
	list, next, err := s.groups.ListGroups(ctx, req.TenantId, int(req.PageSize), req.PageToken)
	if err != nil {
		return nil, err
	}
	resp := &pb.ListGroupsResponse{NextPageToken: next}
	for _, g := range list {
		resp.Groups = append(resp.Groups, toProtoGroup(g))
	}
	return resp, nil
}

func (s *IAMService) AddGroupMember(ctx context.Context, req *pb.AddGroupMemberRequest) (*pb.AddGroupMemberResponse, error) {
	actorID := middleware.GetUserID(ctx)
	return &pb.AddGroupMemberResponse{}, s.groups.AddMember(ctx, req.GroupId, req.UserId, actorID)
}

func (s *IAMService) RemoveGroupMember(ctx context.Context, req *pb.RemoveGroupMemberRequest) (*pb.RemoveGroupMemberResponse, error) {
	actorID := middleware.GetUserID(ctx)
	return &pb.RemoveGroupMemberResponse{}, s.groups.RemoveMember(ctx, req.GroupId, req.UserId, actorID)
}

func (s *IAMService) ListGroupMembers(ctx context.Context, req *pb.ListGroupMembersRequest) (*pb.ListGroupMembersResponse, error) {
	users, err := s.groups.ListMembers(ctx, req.GroupId)
	if err != nil {
		return nil, err
	}
	resp := &pb.ListGroupMembersResponse{}
	for _, u := range users {
		resp.Users = append(resp.Users, toProtoUser(u))
	}
	return resp, nil
}

func (s *IAMService) AssignGroupRole(ctx context.Context, req *pb.AssignGroupRoleRequest) (*pb.AssignGroupRoleResponse, error) {
	actorID := middleware.GetUserID(ctx)
	return &pb.AssignGroupRoleResponse{}, s.groups.AssignRole(ctx, req.GroupId, req.RoleId, actorID)
}

func (s *IAMService) RevokeGroupRole(ctx context.Context, req *pb.RevokeGroupRoleRequest) (*pb.RevokeGroupRoleResponse, error) {
	actorID := middleware.GetUserID(ctx)
	return &pb.RevokeGroupRoleResponse{}, s.groups.RevokeRole(ctx, req.GroupId, req.RoleId, actorID)
}

func (s *IAMService) ListGroupRoles(ctx context.Context, req *pb.ListGroupRolesRequest) (*pb.ListGroupRolesResponse, error) {
	roles, err := s.groups.ListRoles(ctx, req.GroupId)
	if err != nil {
		return nil, err
	}
	resp := &pb.ListGroupRolesResponse{}
	for _, r := range roles {
		resp.Roles = append(resp.Roles, toProtoRole(r))
	}
	return resp, nil
}

// Permissions

func (s *IAMService) CheckPermission(ctx context.Context, req *pb.CheckPermissionRequest) (*pb.CheckPermissionResponse, error) {
	allowed, source, err := s.perms.CheckPermission(ctx, req.UserId, req.TenantId, req.Resource, req.Action, req.ResourceId)
	if err != nil {
		return nil, err
	}
	return &pb.CheckPermissionResponse{Allowed: allowed, Source: source}, nil
}

func (s *IAMService) ListPermissions(ctx context.Context, req *pb.ListPermissionsRequest) (*pb.ListPermissionsResponse, error) {
	perms, err := s.perms.ListPermissions(ctx, req.TenantId, req.RoleId)
	if err != nil {
		return nil, err
	}
	resp := &pb.ListPermissionsResponse{}
	for _, p := range perms {
		resp.Permissions = append(resp.Permissions, p.Resource+":"+p.Action)
	}
	return resp, nil
}

func (s *IAMService) GrantResourcePermission(ctx context.Context, req *pb.GrantResourcePermissionRequest) (*pb.ResourcePermission, error) {
	actorID := middleware.GetUserID(ctx)
	rp, err := s.perms.GrantResourcePermission(ctx, &biz.ResourcePermission{
		UserID:     req.UserId,
		TenantID:   req.TenantId,
		Resource:   req.Resource,
		Action:     req.Action,
		ResourceID: req.ResourceId,
		Allowed:    req.Allowed,
	}, actorID)
	if err != nil {
		return nil, err
	}
	return toProtoResourcePermission(rp), nil
}

func (s *IAMService) RevokeResourcePermission(ctx context.Context, req *pb.RevokeResourcePermissionRequest) (*pb.RevokeResourcePermissionResponse, error) {
	actorID := middleware.GetUserID(ctx)
	return &pb.RevokeResourcePermissionResponse{}, s.perms.RevokeResourcePermission(ctx, req.Id, actorID)
}

func (s *IAMService) ListPendingApprovals(ctx context.Context, req *pb.ListPendingApprovalsRequest) (*pb.ListPendingApprovalsResponse, error) {
	// TODO: Implement ListByStatus in perms
	return &pb.ListPendingApprovalsResponse{}, nil
}

func (s *IAMService) ApprovePermission(ctx context.Context, req *pb.ApprovePermissionRequest) (*pb.ResourcePermission, error) {
	checkerID := middleware.GetUserID(ctx)
	if checkerID == "" {
		return nil, errors.Forbidden("UNAUTHORIZED", "missing subject")
	}

	err := s.perms.ApprovePermission(ctx, req.Id, checkerID)
	if err != nil {
		return nil, err
	}

	rp, err := s.perms.GetResourcePermission(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return toProtoResourcePermission(rp), nil
}

func (s *IAMService) RejectPermission(ctx context.Context, req *pb.RejectPermissionRequest) (*pb.ResourcePermission, error) {
	checkerID := middleware.GetUserID(ctx)
	if checkerID == "" {
		return nil, errors.Forbidden("UNAUTHORIZED", "missing subject")
	}

	err := s.perms.RejectPermission(ctx, req.Id, checkerID)
	if err != nil {
		return nil, err
	}

	rp, err := s.perms.GetResourcePermission(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return toProtoResourcePermission(rp), nil
}

// Converters

func toProtoUser(u *biz.User) *pb.User {
	return &pb.User{
		Id:          u.ID,
		ExternalId:  u.ExternalID,
		Email:       u.Email,
		DisplayName: u.DisplayName,
		CreatedAt:   timestamppb.New(u.CreatedAt),
		UpdatedAt:   timestamppb.New(u.UpdatedAt),
	}
}

func toProtoTenant(t *biz.Tenant) *pb.Tenant {
	return &pb.Tenant{
		Id:        t.ID,
		Name:      t.Name,
		Slug:      t.Slug,
		OwnerId:   t.OwnerID,
		CreatedAt: timestamppb.New(t.CreatedAt),
		UpdatedAt: timestamppb.New(t.UpdatedAt),
	}
}

func toProtoMembership(m *biz.Membership) *pb.Membership {
	return &pb.Membership{
		Id:        m.ID,
		UserId:    m.UserID,
		TenantId:  m.TenantID,
		Role:      m.Role,
		CreatedAt: timestamppb.New(m.CreatedAt),
	}
}

func toProtoRole(r *biz.Role) *pb.Role {
	permStrs := make([]string, len(r.Permissions))
	for i, p := range r.Permissions {
		permStrs[i] = p.ID
	}
	return &pb.Role{
		Id:          r.ID,
		TenantId:    r.TenantID,
		Name:        r.Name,
		Description: r.Description,
		Permissions: permStrs,
		IsSystem:    r.IsSystem,
		CreatedAt:   timestamppb.New(r.CreatedAt),
		UpdatedAt:   timestamppb.New(r.UpdatedAt),
	}
}

func toProtoResourcePermission(rp *biz.ResourcePermission) *pb.ResourcePermission {
	return &pb.ResourcePermission{
		Id:         rp.ID,
		UserId:     rp.UserID,
		TenantId:   rp.TenantID,
		Resource:   rp.Resource,
		Action:     rp.Action,
		ResourceId: rp.ResourceID,
		Allowed:    rp.Allowed,
		CreatedAt:  timestamppb.New(rp.CreatedAt),
	}
}

func toProtoGroup(g *biz.Group) *pb.Group {
	return &pb.Group{
		Id:          g.ID,
		TenantId:    g.TenantID,
		Name:        g.Name,
		Description: g.Description,
		CreatedAt:   timestamppb.New(g.CreatedAt),
		UpdatedAt:   timestamppb.New(g.UpdatedAt),
	}
}
