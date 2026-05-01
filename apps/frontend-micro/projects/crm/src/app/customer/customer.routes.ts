import { Routes } from '@angular/router';
import { CustomerComponent } from './customer';
import { CustomerListComponent } from './info/customer-list';
import { CustomerInfoComponent } from './info/customer-info';
import { InitComponent as RegisterInit } from './register/init';
import { InitComponent as AdjustInit } from './adjust/init';

export const CUSTOMER_ROUTES: Routes = [
  {
    path: '',
    component: CustomerComponent,
    children: [
      {
        path: 'info',
        children: [
          { path: '', component: CustomerListComponent },
          { path: 'details/:id', component: CustomerInfoComponent }
        ]
      },
      {
        path: 'register',
        children: [
          { path: 'init', component: RegisterInit }
        ]
      },
      {
        path: 'adjust',
        children: [
          { path: 'init', component: AdjustInit }
        ]
      },
      {
        path: '',
        redirectTo: 'info',
        pathMatch: 'full'
      }
    ]
  }
];
