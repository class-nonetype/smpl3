import { Component, input, signal } from '@angular/core';
import { FieldTree, FormField } from '@angular/forms/signals';
import { MatButtonModule } from '@angular/material/button';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';

@Component({
  imports: [
    FormField,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatIconModule,
  ],
  selector: 'password-field',
  styleUrl: './password-field.scss',
  templateUrl: './password-field.html',
})
export class PasswordField {

  // espera un campo de formulario cuyo valor es string
  readonly field = input.required<FieldTree<string>>();

  readonly autocomplete = input.required<'current-password' | 'new-password'>();
  readonly label = input('Contraseña');
  readonly hint = input<string>();

  protected readonly hidden = signal(true);
}
