import { Component, inject, signal } from '@angular/core';
import { AuthenticationService, SignInRequest } from '../../../../core/services/authentication';
import { email, form, FormField, FormRoot, required } from '@angular/forms/signals';
import { HttpErrorResponse } from '@angular/common/http';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { Router, RouterLink } from '@angular/router';

@Component({
  imports: [
    FormRoot,
    FormField,
    RouterLink,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatIconModule,
  ],
  selector: 'sign-in-page',
  styleUrl: './sign-in.scss',
  templateUrl: './sign-in.html',
})
export class SignInPage {
  private readonly authenticationService = inject(AuthenticationService);
  private readonly router = inject(Router);

  protected readonly hidePassword = signal(true);


  private readonly credentials = signal<SignInRequest>({ email: '', password: '' });

  protected readonly signInForm = form(
    this.credentials,
    (path) => {
      required(path.email, { message: 'Ingresa tu correo' });
      email(path.email, { message: 'Correo inválido' });
      required(path.password, { message: 'Ingresa tu contraseña' });
    },
    {
      submission: {
        action: async () => {
          try {
            await this.authenticationService.signIn(this.credentials());
            await this.router.navigateByUrl('/');
            return undefined;
          } catch (error) {
            if (error instanceof HttpErrorResponse && error.status === 401) {
              return { kind: 'server', message: 'Correo o contraseña incorrectos' };
            }
            return { kind: 'server', message: 'No se pudo iniciar sesión, intenta de nuevo' };
          }
        },
      },
    },
  );

}
