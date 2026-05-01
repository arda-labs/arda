import { Routes } from '@angular/router';
import { CustomerComponent } from './customer';
import { CustomerInfoComponent } from './info/customer-info';

export const CUSTOMER_ROUTES: Routes = [
  {
    path: '',
    component: CustomerComponent,
    children: [
      {
        path: 'info',
        component: CustomerInfoComponent
      },
      {
        path: '',
        redirectTo: 'info',
        pathMatch: 'full'
      }
    ]
  }
];
