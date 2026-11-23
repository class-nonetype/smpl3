import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';

@Component({
  selector: 'authentication-layout',
  imports: [RouterOutlet],
  templateUrl: './authentication.html',
  styleUrl: './authentication.scss',
})
export class AuthenticationLayout {}
