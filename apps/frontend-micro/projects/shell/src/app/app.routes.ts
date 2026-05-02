import { Routes } from '@angular/router';
import { loadRemoteModule } from '@angular-architects/native-federation';
import { Layout } from './components/layout/layout';
import { LandingPage } from './pages/landing/landing-page';
import { LoginPage } from './pages/login/login-page';
import { CallbackPage } from './pages/auth/callback-page';
import { SelectWorkspace } from './pages/select-workspace/select-workspace';
import { Home } from './pages/home/home';
import { Settings } from './pages/settings/settings';
import { WorkspaceManagement } from './pages/workspaces/workspace-management';
import { authGuard } from './guards/auth.guard';
import { AccessDeniedPage } from './pages/error/access-denied';
import { NotFoundPage } from './pages/error/not-found';

export const routes: Routes = [
  // Landing/Public
  {
    path: '',
    pathMatch: 'full',
    component: LandingPage,
  },

  // Login/Auth
  {
    path: 'login',
    component: LoginPage,
  },
  {
    path: 'ui/v2/login/login',
    component: LoginPage,
  },
  {
    path: 'auth/callback',
    component: CallbackPage,
  },

  // Workspace Selection
  {
    path: 'select-workspace',
    component: SelectWorkspace,
    canActivate: [authGuard],
  },

  // Authenticated Shell
  {
    path: '',
    component: Layout,
    canActivate: [authGuard],
    children: [
      {
        path: 'home',
        component: Home,
      },
      {
        path: 'settings',
        component: Settings,
      },
      {
        path: 'workspaces',
        component: WorkspaceManagement,
      },
      {
        path: 'iam',
        loadChildren: () =>
          loadRemoteModule('iam', './Routes').then((m) => m.remoteRoutes),
      },
      {
        path: 'mdm',
        loadChildren: () =>
          loadRemoteModule('mdm', './Routes').then((m) => m.remoteRoutes),
      },
      {
        path: 'ntf',
        loadChildren: () =>
          loadRemoteModule('ntf', './Routes').then((m) => m.remoteRoutes),
      },
      {
        path: 'crm',
        loadChildren: () =>
          loadRemoteModule('crm', './Routes').then((m) => m.remoteRoutes),
      },
      {
        path: 'hrm',
        loadChildren: () =>
          loadRemoteModule('hrm', './Routes').then((m) => m.remoteRoutes),
      },
      {
        path: 'loan',
        loadChildren: () =>
          loadRemoteModule('loan', './Routes').then((m) => m.remoteRoutes),
      },
      {
        path: 'bpm',
        loadChildren: () =>
          loadRemoteModule('bpm', './Routes').then((m) => m.remoteRoutes),
      },
    ],
  },

  // Error Pages
  {
    path: '403',
    component: AccessDeniedPage,
  },
  {
    path: '404',
    component: NotFoundPage,
  },

  // Wildcard
  {
    path: '**',
    redirectTo: '404',
  },
];
