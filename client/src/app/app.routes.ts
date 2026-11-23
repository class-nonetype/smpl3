import { Routes } from '@angular/router';
import { authenticatedGuard, guestGuard } from './core/guards/authentication';

export const routes: Routes = [
  {
    path: 'authentication',
    canMatch: [guestGuard],
    loadChildren: () =>
      import('./features/authentication/authentication.routes').then(
        (m) => m.authenticationRoutes
      ),
  },
  {
    path: '',
    canMatch: [authenticatedGuard],
    loadChildren: () => import('./features/home/home.routes').then((m) => m.homeRoutes),
  },
  { path: '**', redirectTo: '' },
];
