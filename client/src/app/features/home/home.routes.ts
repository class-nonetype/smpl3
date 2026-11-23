import { Routes } from '@angular/router';

export const homeRoutes: Routes = [
  {
    path: '',
    title: 'Kin',
    loadComponent: () => import('./pages/home/home').then((m) => m.HomePage),
  },
];
