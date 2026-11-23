import { ApplicationConfig, inject, provideAppInitializer, provideBrowserGlobalErrorListeners } from '@angular/core';
import { provideRouter } from '@angular/router';
import { routes } from './app.routes';
import { MAT_FORM_FIELD_DEFAULT_OPTIONS } from '@angular/material/form-field';
import { AuthenticationService } from './core/services/authentication';

export const appConfig: ApplicationConfig = {
  providers: [
    provideBrowserGlobalErrorListeners(),
    provideRouter(routes),

    /*
    ejecuta la función antes de renderizar la app, y si devuelve una promesa
    la espera.
    asi, cuando se pinta la primera pantalla, ya se sabe si hay sesión,
    y no aparece el formulario por un instante antes de saltar.
    */
    provideAppInitializer(() => inject(AuthenticationService).restore()),
    { provide: MAT_FORM_FIELD_DEFAULT_OPTIONS, useValue: { appearance: 'fill' } }
  ]
};
