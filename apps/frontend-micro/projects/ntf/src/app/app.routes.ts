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
  {
    path: 'template-upload',
    loadComponent: () => import('./pages/operations/template-upload.page').then(m => m.TemplateUploadPage),
  },
];

export const routes: Routes = remoteRoutes;
