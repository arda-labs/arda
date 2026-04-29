import { Routes } from '@angular/router';

export const remoteRoutes: Routes = [
  {
    path: '',
    pathMatch: 'full',
    redirectTo: 'users'
  },
  {
    path: 'users/:id',
    loadComponent: () => import('./pages/iam/users/user-detail').then(m => m.UserDetail)
  },
  {
    path: 'users',
    loadComponent: () => import('./pages/iam/users/user-management').then(m => m.UserManagement)
  },
  {
    path: 'roles',
    loadComponent: () => import('./pages/iam/roles/role-management').then(m => m.RoleManagement)
  },
  {
    path: 'groups',
    loadComponent: () => import('./pages/iam/groups/group-management').then(m => m.GroupManagement)
  },
  {
    path: 'menus',
    loadComponent: () => import('./pages/iam/menus/menu-management/menu-management').then(m => m.MenuManagement)
  },
  {
    path: 'profile',
    loadComponent: () => import('./pages/profile/profile').then(m => m.Profile)
  }
];

export const routes: Routes = remoteRoutes;
