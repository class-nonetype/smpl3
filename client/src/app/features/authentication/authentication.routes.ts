import { Routes } from '@angular/router';

export const authenticationRoutes: Routes = [
  {
    path: '',
    loadComponent: () =>
      import('./layout/authentication').then((m) => m.AuthenticationLayout),
    children: [
      { path: '', pathMatch: 'full', redirectTo: 'sign-in' },
      {
        path: 'sign-in',
        title: 'Iniciar sesión · Kin',
        loadComponent: () => import('./pages/sign-in/sign-in').then((m) => m.SignInPage),
      },
      {
        path: 'sign-up',
        title: 'Crea tu cuenta · Kin',
        loadComponent: () => import('./pages/sign-up/sign-up').then((m) => m.SignUpPage),
      }
    ],
  },
];
