import { Routes } from '@angular/router';

export const remoteRoutes: Routes = [
  {
    path: '',
    pathMatch: 'full',
    redirectTo: 'geo/administrative-units',
  },
  {
    path: 'geo/administrative-units',
    loadComponent: () => import('./pages/geo/administrative-units.page').then(m => m.AdministrativeUnitsPage),
  },
  {
    path: 'geo/area-types',
    loadComponent: () => import('./pages/geo/area-types.page').then(m => m.AreaTypesPage),
  },
  {
    path: 'geo/areas',
    loadComponent: () => import('./pages/geo/areas.page').then(m => m.AreasPage),
  },
  {
    path: 'catalog/code-sets',
    loadComponent: () => import('./pages/catalog/code-sets.page').then(m => m.CodeSetsPage),
  },
  {
    path: 'catalog/code-items',
    loadComponent: () => import('./pages/catalog/code-items.page').then(m => m.CodeItemsPage),
  },
  {
    path: 'system/parameters',
    loadComponent: () => import('./pages/system/system-parameters.page').then(m => m.SystemParametersPage),
  },
  {
    path: 'banking/reference',
    loadComponent: () => import('./pages/banking/banking-reference.page').then(m => m.BankingReferencePage),
  },
  {
    path: 'banking/credit-institutions',
    loadComponent: () => import('./pages/banking/credit-institutions.page').then(m => m.CreditInstitutionsPage),
  },
  {
    path: 'banking/business-calendars',
    loadComponent: () => import('./pages/banking/business-calendars.page').then(m => m.BusinessCalendarsPage),
  },
  {
    path: 'banking/pricing-rules',
    loadComponent: () => import('./pages/banking/pricing-rules.page').then(m => m.PricingRulesPage),
  },
];

export const routes: Routes = remoteRoutes;
