import { Routes } from '@angular/router';

export const remoteRoutes: Routes = [
  {
    path: 'inbound',
    loadComponent: () => import('./pages/inbound/inbound').then(m => m.InboundComponent)
  },
  {
    path: 'outbound',
    loadComponent: () => import('./pages/outbound/outbound').then(m => m.OutboundComponent)
  },
  {
    path: 'definitions',
    loadComponent: () => import('./pages/definitions/definitions').then(m => m.DefinitionsComponent)
  },
  {
    path: 'deploy',
    loadComponent: () => import('./pages/deploy/deploy').then(m => m.DeployComponent)
  },
  {
    path: 'search',
    loadComponent: () => import('./pages/search/search').then(m => m.SearchComponent)
  },
  {
    path: 'monitor',
    loadComponent: () => import('./pages/monitor/monitor').then(m => m.MonitorComponent)
  },
  {
    path: 'error-hospital',
    loadComponent: () => import('./pages/error-hospital/error-hospital').then(m => m.ErrorHospitalComponent)
  },
  {
    path: 'config',
    children: [
      {
        path: 'sla',
        loadComponent: () => import('./pages/config/sla/sla').then(m => m.SlaComponent)
      },
      {
        path: 'description',
        loadComponent: () => import('./pages/config/description/description').then(m => m.DescriptionComponent)
      },
      {
        path: 'assignment',
        loadComponent: () => import('./pages/config/assignment/assignment').then(m => m.AssignmentComponent)
      }
    ]
  },
  {
    path: '',
    redirectTo: 'inbound',
    pathMatch: 'full'
  }
];

export const routes: Routes = remoteRoutes;
