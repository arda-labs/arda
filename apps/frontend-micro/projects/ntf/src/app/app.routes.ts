import { Routes } from '@angular/router';

export const remoteRoutes: Routes = [
  {
    path: '',
    pathMatch: 'full',
    redirectTo: 'operations',
  },
  {
    path: 'operations',
    loadComponent: () => import('./pages/operations/notification-operations.page').then(m => m.NotificationOperationsPage),
  },
];

export const routes: Routes = remoteRoutes;
