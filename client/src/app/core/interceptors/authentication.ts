import { HttpErrorResponse, HttpInterceptorFn } from '@angular/common/http';
import { inject } from '@angular/core';
import { Router } from '@angular/router';
import { catchError, throwError } from 'rxjs';
import { AuthenticationService } from '../services/authentication';

// middleware del cliente http
export const authenticationInterceptor: HttpInterceptorFn = (request, next) => {
  const authentication = inject(AuthenticationService);
  const router = inject(Router);

  // toda petición pasa por aqui
  return next(request).pipe(
    catchError((error: unknown) => {

      // se excluyen las rutas /authentication/
      // un 401 en sign-in significa credenciales malas
      // y en /me significa que no hay sesión
      const expired =
        error instanceof HttpErrorResponse &&
        error.status === 401 &&
        !request.url.includes('/authentication/');

      if (expired) {
        authentication.clear();
        void router.navigateByUrl('/authentication');
      }
      return throwError(() => error);
    }),
  );
};
