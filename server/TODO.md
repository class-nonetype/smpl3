# Kin server — TODO

Última sesión: 2026-11-22

## Cómo trabajamos

- Yo escribo el código. Claude me guía **un archivo por vez**, bloques cortos, sin comentarios en el código y explicando buenas prácticas.
- Claude solo edita archivos cuando se lo pido explícitamente.
- Tests al final, cuando estén cerradas las funcionalidades.

## Filosofía

**El objetivo no es solo terminar Kin: es convertirme en desarrollador Go.** Kin es el vehículo para aprender.

### Del código
- **Purista: librería estándar primero.** `net/http` en vez de Gin/Chi, `slog` en vez de zap, `os.Getenv` en vez de librerías de config. Solo se agrega una dependencia si resuelve algo que no vale la pena escribir (pgx, sqlc, goose, x/crypto, uuid).
- **Si es corto y lo entiendo, lo escribo yo.** Ej: el router con grupos son ~30 líneas propias en vez de una dependencia.
- **SQL real y visible.** Nada de ORM: migraciones `.sql` versionadas (goose) y queries `.sql` que sqlc convierte en Go tipado.
- **Idiomático, no mágico.** Errores explícitos con `%w`, dependencias por constructor (sin frameworks de DI), conexión de todo en `main` (composition root), `context` para pasar datos por petición.
- **Estructura:** `cmd/` + `internal/`; funcionalidades en `internal/feature/`, lo técnico plano en `internal/`. Sin `utils`, `models`, `controllers`.
- **Capas solo cuando aportan.** Servicio cuando hay lógica de negocio (auth); handler directo a `store.Queries` cuando es CRUD (posts).
- **Seguridad y producción desde el inicio:** timeouts, graceful shutdown, Argon2id, sesiones con token hasheado, cookies HttpOnly, validación en la app y `CHECK`s en la BD.
- **Mejores prácticas actuales (noviembre 2026),** no tutoriales viejos.

### Del aprendizaje
- **Entiendo el código cuando lo veo, pero todavía no lo escribo solo** (reconocer vs recordar). Es normal; la meta es pasar a recordar.
- **Por eso:** en las próximas funcionalidades (ej: comentarios) **yo intento primero** copiando el patrón existente, y Claude revisa y corrige después. No darme la solución de entrada.
- **Ejercicio útil:** reescribir de memoria un archivo ya hecho y ver dónde me trabo.
- Saber **qué existe y por qué** importa más que memorizar sintaxis.

### Del estilo de Claude
- Respuestas cortas y directas, en español, sin rodeos.
- **Un archivo por mensaje,** bloques de código chicos, 3-5 bullets explicando el porqué. Nada de pantallas llenas de código.
- **Sin comentarios en el código** que me pasa (los comentarios en mis archivos son apuntes míos; no tocarlos).
- Opinión clara cuando pregunto "qué es mejor": recomendar una opción, no listar todas.
- No crear ni editar archivos salvo que lo pida explícitamente.
- **Verificar antes de recomendar:** revisar la versión instalada (`node_modules`, `go.mod`) y la documentación actual. Nada de memoria ni de APIs deprecadas.

### Aplica también al front (Angular)
- Las mismas reglas: Claude **no escribe el proyecto**, solo guía y manda código **de a un archivo por mensaje**, nunca todo junto.

## Frontend (client/)

Angular 22 · Angular Material · proxy de `ng serve` hacia el server Go (`/api` → `:8080`)

- [x] Material 22: tema M3, inputs estilo "standard" (fill transparente vía `form-field-overrides`)
- [x] Proxy `/api` → `:8080` (la cookie HttpOnly funciona sin CORS)
- [x] `AuthenticationService`: signIn, signUp, signOut, restore (`/me`) con signal `user`
- [x] Restaurar sesión al arrancar con `provideAppInitializer`
- [x] Sign-in con Signal Forms (`form`, `[formRoot]`, `[formField]`) y errores del server
- [x] Guards `authenticatedGuard` / `guestGuard` (`canMatch`) y `home` con sign-out

### Pendiente front
- [x] Estructura: `AuthenticationLayout` (centrado + outlet) con hijas `sign-in` / `sign-up`; `title` en sign-in
- [ ] Sign-up en `pages/sign-up` (lo intento yo primero, Claude revisa)
- [ ] Interceptor: si una petición da 401, limpiar `user` y mandar al login
- [ ] Feed de posts en `home`
- [ ] Limpiar `app.ts` (`title` sin uso) y arreglar los specs (`app.spec.ts`; `authentication.spec.ts` fue borrado)

## Stack

Go 1.27 · `net/http` (sin frameworks) · pgx · sqlc · goose · Postgres 18 (BD `kin`, usuario `development`)

## Estado actual

- [x] `main.go` con `run()`, timeouts, graceful shutdown, `net.Listen` y log de rutas y URL
- [x] `config` (APP_ENV, ADDR, DATABASE_URL, ALLOWED_ORIGINS)
- [x] `database.Connect` (pgxpool + ping)
- [x] `httpx`: WriteJSON, DecodeJSON, Error/WriteError, HandlerFunc, Router con grupos y middlewares
- [x] Middlewares: Recover, Logger, CORS
- [x] Migraciones: users, sessions, posts
- [x] Auth: Argon2id, sesiones en BD (cookie `kin_session`), sign-up, sign-in, sign-out, me, RequireUser
- [x] Posts: create, get, list (cursor), delete (solo el autor)
- [x] Probado de punta a punta con curl, todo OK

## Pendiente (en orden)

### Rápidos
- [ ] Errores con contexto en `service.go`: cambiar `fmt.Errorf("%w", err)` por `fmt.Errorf("crear sesión: %w", err)`, etc.
- [ ] `SESSION_TTL` en config/.env/.env.example (hoy 30 días fijo en `main.go`)
- [ ] Limpiar sesiones vencidas: goroutine con ticker que llame `DeleteExpiredSessions`
- [ ] `Makefile`: run, migrate-up/down/status, migrate-new, generate, build, lint

### Calidad
- [ ] `/healthz` que haga ping a la BD (o `/readyz` aparte)
- [ ] Tipo propio `authentication.User` en vez de exponer `store.GetSessionUserRow`
- [ ] Rate limiting en sign-in
- [ ] Código de error estable en el JSON (`{"error": "...", "code": "email_taken"}`) para que el front no dependa del texto del mensaje en el 409 del sign-up
- [ ] Protección CSRF extra con `http.CrossOriginProtection` (Go 1.25+)
- [ ] Request ID middleware (header `X-Request-ID` + en logs)
- [ ] `updated_at`: actualizarlo al editar (trigger o en la query)

### Funcionalidades
- [ ] Editar post (PATCH)
- [ ] Comentarios (migración + queries + handler)
- [ ] Subida de archivos
- [ ] Login con Google (requiere Client ID en Google Cloud Console)

### Producción
- [ ] Tests (unitarios + handlers con `httptest`)
- [ ] Dockerfile
- [ ] CI con GitHub Actions (vet, golangci-lint, tests)
- [ ] Evaluar cambiar `github.com/google/uuid` por el `uuid` de la stdlib (Go 1.27) si pgx/sqlc lo soportan

## Reglas aprendidas

- **Nunca editar una migración ya aplicada.** Para cambiar algo, crear una migración nueva.
- Cargar el `.env` antes de correr: `set -a; source .env; set +a`
- Después de cambiar queries o migraciones: `go tool sqlc generate`
