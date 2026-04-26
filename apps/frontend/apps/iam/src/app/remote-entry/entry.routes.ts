import { Route } from '@angular/router';

export const remoteRoutes: Route[] = [
  {
    path: 'users',
    loadComponent: () => import('../pages/iam/users/user-management').then(m => m.UserManagement)
  },
  {
    path: 'roles',
    loadComponent: () => import('../pages/iam/roles/role-management').then(m => m.RoleManagement)
  },
  {
    path: 'groups',
    loadComponent: () => import('../pages/iam/groups/group-management').then(m => m.GroupManagement)
  },
  {
    path: 'profile',
    loadComponent: () => import('../pages/profile/profile').then(m => m.Profile)
  }
];
