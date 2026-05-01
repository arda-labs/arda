import { Routes } from '@angular/router';

export const routes: Routes = [
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
