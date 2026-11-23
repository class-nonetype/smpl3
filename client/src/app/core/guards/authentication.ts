import { inject } from "@angular/core";
import { CanMatchFn, Router } from "@angular/router";
import { AuthenticationService } from "../services/authentication";



/*

retornar un UrlTree es mejor que hacer router.navigate() y devolver false.
el router cancela la navegación y redirige en un soo paso,
sin dejar una navegación a medias en el historial.

puede ser sincrono gracias a provideAppInitializer: cuando corre cualquier guarda,
restore() ya terminó y user() ya tiene el v alor real.
sin eso habria que esperar la respuesta de /me dentro del guard

inject dentro del guard funciona por lo mismo que en el initializer:
el router lo ejecuta dentro del contexto de inyección.
*/

export const authenticatedGuard: CanMatchFn = () => {
  const user = inject(AuthenticationService).user();
  return user ? true : inject(Router).createUrlTree(['/authentication']);
};

export const guestGuard: CanMatchFn = () => {
  const user = inject(AuthenticationService).user();
  return user ? inject(Router).createUrlTree(['/']) : true;
};
