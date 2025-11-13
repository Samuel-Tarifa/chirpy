# Chirpy

Chirpy es una aplicación de red social tipo Twitter construida con **Go** en el backend y **React** en el frontend. Permite a los usuarios crear, leer y gestionar "chirps" (tweets cortos), autenticarse con JWT y gestionar su perfil.

## Tecnologías

- **Backend:** Go 1.22+
- **Base de datos:** PostgreSQL
- **Autenticación:** JWT (JSON Web Tokens)
- **Frontend:** React
- **Migraciones:** Goose

---

## Instalación y Setup

### Requisitos previos

- Go 1.22 o superior
- PostgreSQL
- Node.js y npm (para el frontend)

### 1. Clonar el repositorio

```sh
git clone https://github.com/Samuel-Tarifa/chirpy.git
cd chirpy
```

### 2. Configurar variables de entorno

Crea un archivo `.env` en la raíz del proyecto:

```env
DB_URL=postgres://postgres:postgres@localhost:5432/chirpy?sslmode=disable
JWT_SECRET=tu_clave_secreta_muy_segura_aqui
PLATFORM=dev
POLKA_KEY=tu_polka_key_aqui
```

### 3. Instalar dependencias de Go

```sh
go mod download
```

### 4. Ejecutar migraciones

```sh
chmod +x ./up.sh
./up.sh
```

O manualmente:

```sh
goose -dir sql/schema postgres "$DB_URL" up
```

### 5. Ejecutar el servidor

```sh
go run .
```

El servidor estará disponible en `http://localhost:8080`.

---

## Endpoints de la API

### Health Check

#### `GET /api/healthz`

Verifica que el servidor está funcionando.

**Respuesta:**

```
OK
```

---

### Chirps

#### `POST /api/chirps`

Crear un nuevo chirp.

**Headers:**

```
Authorization: Bearer <token>
Content-Type: application/json
```

**Body:**

```json
{
  "body": "Este es mi primer chirp!"
}
```

**Respuesta (201):**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "created_at": "2025-11-13T10:30:00Z",
  "updated_at": "2025-11-13T10:30:00Z",
  "body": "Este es mi primer chirp!",
  "user_id": "550e8400-e29b-41d4-a716-446655440001"
}
```

---

#### `GET /api/chirps`

Obtener todos los chirps (con filtro opcional por autor).

**Query Parameters:**

- `author_id` (opcional): UUID del autor para filtrar chirps

**Respuesta (200):**

```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "created_at": "2025-11-13T10:30:00Z",
    "updated_at": "2025-11-13T10:30:00Z",
    "body": "Este es mi primer chirp!",
    "user_id": "550e8400-e29b-41d4-a716-446655440001"
  }
]
```

---

#### `GET /api/chirps/{id}`

Obtener un chirp específico por ID.

**Respuesta (200):**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "created_at": "2025-11-13T10:30:00Z",
  "updated_at": "2025-11-13T10:30:00Z",
  "body": "Este es mi primer chirp!",
  "user_id": "550e8400-e29b-41d4-a716-446655440001"
}
```

---

### Usuarios

#### `POST /api/users`

Crear un nuevo usuario.

**Body:**

```json
{
  "email": "usuario@example.com",
  "password": "micontraseña123"
}
```

**Respuesta (201):**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440001",
  "created_at": "2025-11-13T10:00:00Z",
  "updated_at": "2025-11-13T10:00:00Z",
  "email": "usuario@example.com",
  "is_chirpy_red": false
}
```

---

#### `POST /api/login`

Autenticar un usuario y obtener tokens JWT.

**Body:**

```json
{
  "email": "usuario@example.com",
  "password": "micontraseña123"
}
```

**Respuesta (200):**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440001",
  "created_at": "2025-11-13T10:00:00Z",
  "updated_at": "2025-11-13T10:00:00Z",
  "email": "usuario@example.com",
  "is_chirpy_red": false,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

---

#### `PUT /api/users`

Actualizar datos del usuario (requiere autenticación).

**Headers:**

```
Authorization: Bearer <token>
Content-Type: application/json
```

**Body:**

```json
{
  "email": "nuevomail@example.com",
  "password": "nuevacontraseña123"
}
```

**Respuesta (200):**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440001",
  "created_at": "2025-11-13T10:00:00Z",
  "updated_at": "2025-11-13T11:00:00Z",
  "email": "nuevomail@example.com",
  "is_chirpy_red": false
}
```

---

### Autenticación

#### `POST /api/refresh`

Refrescar el token de acceso usando el refresh token.

**Body:**

```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Respuesta (200):**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

---

#### `POST /api/revoke`

Revocar un refresh token (logout).

**Headers:**

```
Authorization: Bearer <refresh_token>
```

**Respuesta (204):**
Sin contenido.

---

### Admin

#### `GET /admin/metrics`

Obtener métricas de uso (solo admin).

**Respuesta (200):**

```html
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited 42 times!</p>
  </body>
</html>
```

---

#### `POST /admin/reset`

Resetear contadores de métricas (solo admin).

**Respuesta (200):**

```
Number of hits reseted.
```

---

## Estructura del Proyecto

```
chirpy/
├── main.go                 # Punto de entrada
├── middleware.go           # Middlewares
├── handlers.go             # Handlers de los endpoints
├── fileserver.go           # Servidor de archivos estáticos
├── types.go                # Definiciones de structs
├── internal/
│   ├── auth/               # Funciones de autenticación (JWT, hash)
│   └── database/           # Queries generadas por sqlc
├── sql/
│   ├── schema/             # Migraciones de Goose
│   └── queries/            # Queries SQL
├── frontend/               # Aplicación React (opcional)
│   ├── src/
│   ├── public/
│   └── build/              # Archivos compilados
└── .env                    # Variables de entorno (no commitear)
```

---

## Desarrollo

### Generar código desde queries SQL

```sh
sqlc generate
```

### Crear una nueva migración

```sh
goose -dir sql/schema create nombre_migracion sql
```

### Revertir migraciones

```sh
goose -dir sql/schema postgres "$DB_URL" down
```

---

## Errores comunes

| Error                    | Solución                                                            |
| ------------------------ | ------------------------------------------------------------------- |
| `error opening database` | Verifica que PostgreSQL está corriendo y el `DB_URL` es correcto    |
| `Method Not Allowed`     | Verifica que estás usando el método HTTP correcto (GET, POST, etc.) |
| `invalid token`          | El JWT ha expirado o fue firmado con otra clave secreta             |
| `invalid credentials`    | Email o contraseña incorrectos                                      |

---

## Licencia

MIT

---

## Autor

Samuel Tarifa - [GitHub](https://github.com/Samuel-Tarifa)
