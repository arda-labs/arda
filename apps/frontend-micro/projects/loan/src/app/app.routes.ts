import { Routes } from '@angular/router';

export const remoteRoutes: Routes = [
  {
    path: 'customer',
    loadChildren: () => import('./customer/customer.routes').then(m => m.CUSTOMER_ROUTES)
  },
  {
    path: '',
    redirectTo: 'customer',
    pathMatch: 'full'
  }
];

export const routes: Routes = remoteRoutes;
