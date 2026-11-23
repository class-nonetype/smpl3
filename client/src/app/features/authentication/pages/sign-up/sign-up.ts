import { Component, inject, signal } from '@angular/core';
import { AuthenticationService, SignUpRequest } from '../../../../core/services/authentication';
import { Router, RouterLink } from '@angular/router';
import {
  email,
  form,
  FormField,
  FormRoot,
  maxLength,
  minLength,
  pattern,
  required,
} from '@angular/forms/signals';
import { HttpErrorResponse } from '@angular/common/http';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';

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
  selector: 'sign-up-page',
  styleUrl: './sign-up.scss',
  templateUrl: './sign-up.html',
})
export class SignUpPage {
  private readonly authenticationService = inject(AuthenticationService);
  private readonly router = inject(Router);

  protected readonly hidePassword = signal(true);

  private readonly account = signal<SignUpRequest>({
    email: '',
    username: '',
    password: '',
  });

  protected readonly signUpForm = form(
    this.account,
    (path) => {
      required(path.username, { message: 'Ingresa un nombre de usuario' });
      minLength(path.username, 3, { message: 'Mínimo 3 caracteres' });
      maxLength(path.username, 30, { message: 'Máximo 30 caracteres' });
      pattern(path.username, /^[a-zA-Z0-9_]+$/, { message: 'Solo letras, números y _' });

      required(path.email, { message: 'Ingresa tu correo' });
      email(path.email, { message: 'Correo inválido' });

      required(path.password, { message: 'Ingresa tu contraseña' });
      minLength(path.password, 8, { message: 'Mínimo 8 caracteres' });
      maxLength(path.password, 72, { message: 'Máximo 72 caracteres' });
    },
    {
      submission: {
        action: async (field) => {
          const { email, password } = this.account();

          try {
            await this.authenticationService.signUp(this.account());
          } catch (error) {
            if (error instanceof HttpErrorResponse && error.status === 409) {
              const message: string = error.error?.error ?? '';
              if (message.includes('email')) {
                return { kind: 'server', fieldTree: field.email, message: 'Ese correo ya está registrado' };
              }
              return { kind: 'server', fieldTree: field.username, message: 'Ese nombre de usuario ya existe' };
            }
            return { kind: 'server', message: 'No se pudo crear la cuenta, intenta de nuevo' };
          }

          try {
            await this.authenticationService.signIn({ email, password });
            await this.router.navigateByUrl('/');
          } catch {
            await this.router.navigateByUrl('/authentication/sign-in');
          }
          return undefined;
        },
      },
    },
  );
}
