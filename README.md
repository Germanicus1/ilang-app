# iLang Backend API

A Golang-based backend for managing users and authentication with Supabase. This API supports CRUD operations, JWT-based authentication, and token management. No Supabase tokens are exposed to the client.

## Features

- **User Management**: Create, update, retrieve, and delete user accounts.
- **Authentication**: Supabase-based user authentication with JWT.
- **Role-Based Authorization**: Secure routes using roles (e.g., admin, user).
- **Middleware**: Validates JWTs for secured endpoints.
- **Configuration**: Centralized `.env`-based configuration.

---

## Project Structure

```TXT

├── main.go               # Entry point of the application
├── handlers/             # Contains HTTP handler functions
│   ├── UserHandlers.go   # Handlers for user-related CRUD operations
├── middleware/           # Middleware for request validation
│   ├── AuthMiddleware.go # Middleware to validate JWT tokens
├── routes/               # Route registration
│   ├── publicRoutes.go   # Routes accessible without authentication
│   ├── securedRoutes.go  # Secured routes requiring JWT
├── config/               # Configuration loading from .env
│   ├── Config.go         # Handles environment variable loading

```

---

## Setup Instructions

### Prerequisites

- [Golang](https://go.dev/) (v1.22 or higher)
- [Supabase](https://supabase.com/) account with a configured project
- `.env` file with the following keys:

    `SUPABASE_URL=<Your Supabase URL> SUPABASE_KEY=<Your Supabase Public Key> JWT_SECRET=<Your Supabase JWT Secret> SERVICE_ROLE_KEY=<Your Supabase Service Role Key>`


### Installation

1. Clone the repository:

    `git clone https://github.com/<your-repo>.git cd ilang-backend`

2. Install dependencies:

    `go mod tidy`

3. Run the server:

    `go run main.go`

4. Server runs at `http://localhost:8080`.


---

## API Endpoints

### Public Routes

|Method|Endpoint|Description|
|---|---|---|
|POST|`/users`|Create a user. Creates on public.users and auth.users|
|POST|`/login`|User login|
|POST|`/logout`|User logout (optional)|

### Secured Routes

|Method|Endpoint|Description|
|---|---|---|
|GET|`/users/{id}`|Get user by ID|
|PATCH|`/users/{id}`|Update user by ID. Patches on public.users and auth.users|
|DELETE|`/users/{id}`|Delete user by ID. Deletes on public.users and auth.users|

---

## Development Notes

### Middleware

- **AuthMiddleware**: Validates incoming JWTs using the Supabase secret. Adds `userID` and `role` to the request context.

### Token Refresh

- Needs to be handled by the frontend

---

## Future Enhancements

- Add unit tests for handlers and middleware.
- Implement more granular roles and permissions.
- Enhance error handling and logging.

---

Feel free to customize this README as per your project requirements! Let me know if you need further refinements. ​​