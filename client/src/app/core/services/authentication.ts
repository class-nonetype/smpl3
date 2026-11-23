import { HttpClient } from '@angular/common/http';
import { Service, inject, signal } from '@angular/core';
import { firstValueFrom } from 'rxjs';
import { environment } from '../../../environments/environment';

export interface User {
  id: string;
  username: string;
  email: string;
}

export interface SignInRequest {
  email: string;
  password: string;
}

export interface SignUpRequest extends SignInRequest {
  username: string;
}


@Service()
export class AuthenticationService {
  private readonly http = inject(HttpClient);
  private readonly url = `${environment.apiUrl}/authentication`;
  private readonly currentUser = signal<User | null>(null);

  readonly user = this.currentUser.asReadonly();


  async signIn(request: SignInRequest): Promise<User> {
    const user = await firstValueFrom(this.http.post<User>(`${this.url}/sign-in`, request));
    this.currentUser.set(user);
    return user;
  }

  async signUp(request: SignUpRequest): Promise<User> {
    return firstValueFrom(this.http.post<User>(`${this.url}/sign-up`, request));
  }

  async signOut(): Promise<void> {
    await firstValueFrom(this.http.post<void>(`${this.url}/sign-out`, null));
    this.currentUser.set(null);
  }

  /*
  nunca lanza un error
  va a correr dentro de provideAppInitializer
  y si esa promesa falla la app no arranca
  quedaria la pantalla en blanco por un 401 o por el server caido
  cualquier fallo significa simplemente 'no hay sesion'
  */
  async restore(): Promise<void> {
    try {
      const user = await firstValueFrom(this.http.get<User>(`${this.url}/me`));
      this.currentUser.set(user);
    } catch {
      this.currentUser.set(null);
    }
  }
}
