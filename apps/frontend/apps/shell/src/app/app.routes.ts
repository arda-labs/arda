import { Route } from '@angular/router';
import { Layout } from './components/layout/layout';
import { Home } from './pages/home/home';
import { LoginPage } from './pages/login/login-page';
import { CallbackPage } from './pages/auth/callback-page';
import { authGuard } from './guards/auth.guard';
import { loadRemote } from '@module-federation/enhanced/runtime';

export const appRoutes: Route[] = [
  // Trang gốc → redirect sang dashboard/home (guard check xảy ra trên các route con)
  {
    path: '',
    pathMatch: 'full',
    redirectTo: 'home'
  },

  // Login
  {
    path: 'login',
    loadComponent: () =>
      import('./pages/login/login-page').then((m) => m.LoginPage),
  },

  // OIDC callback
  {
    path: 'auth/callback',
    component: CallbackPage,
  },

  // Chọn workspace (sau login nếu user thuộc nhiều tenant)
  {
    path: 'select-workspace',
    loadComponent: () =>
      import('./pages/select-workspace/select-workspace').then(
        (m) => m.SelectWorkspace,
      ),
    canActivate: [authGuard],
  },

  // Authenticated app (No /app prefix)
  {
    path: '',
    component: Layout,
    canActivate: [authGuard],
    children: [
      { path: 'home', component: Home },
      {
        path: 'common',
        loadChildren: () =>
          loadRemote<typeof import('common/Routes')>('common/Routes').then(
            (m) => m?.remoteRoutes ?? [],
          ),
      },
      {
        path: 'settings',
        loadComponent: () =>
          import('./pages/settings/settings').then((m) => m.Settings),
      },
      {
        path: 'profile',
        loadComponent: () =>
          import('./pages/profile/profile').then((m) => m.Profile),
      },
    ],
  },

  // Error pages
  {
    path: '403',
    loadComponent: () =>
      import('./pages/error/access-denied').then((m) => m.AccessDeniedPage),
  },
  {
    path: '404',
    loadComponent: () =>
      import('./pages/error/not-found').then((m) => m.NotFoundPage),
  },

  // Wildcard fallback → 404
  { path: '**', redirectTo: '404' },
];
// Trigger CI/CD with new PAT scopes
