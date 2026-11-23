import { Component, inject, signal } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { AuthenticationService } from '../../../../core/services/authentication';
import { Router } from '@angular/router';

@Component({
  imports: [
    MatButtonModule,
  ],
  selector: 'home-page',
  styleUrl: './home.scss',
  templateUrl: './home.html',
})
export class HomePage {
  private readonly authenticationService = inject(AuthenticationService);
  private readonly router = inject(Router);

  protected readonly user = this.authenticationService.user;
  protected readonly signingOut = signal(false);

  protected async signOut(): Promise<void> {
    this.signingOut.set(true);
    try {
      await this.authenticationService.signOut();
      await this.router.navigateByUrl('/authentication');
    } finally {
      this.signingOut.set(false);
    }
  }
}
