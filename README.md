# iLang Backend
**STATUS**: NOT RELEASED

The iLang backend is a RESTful API designed to support a language-learning platform. It provides endpoints for managing games, fetching data, and integrating with a Supabase database.

## Features

- **Health Check Endpoint**: Simple endpoint to check if the server is running.
- **CRUD Operations for Games**:
  - Fetch all games.
  - Fetch a specific game by ID.
  - Create a new game.
- **Database Integration**: Uses Supabase as the backend database for managing data.
- **Environment Configurations**: Securely loads configurations such as API keys and database URLs from environment variables.

---

## Endpoints

### Summary of CRUD Endpoints
Current CRUD operations for `games`:

| **Method** | **Endpoint**        | **Description**         |
|------------|---------------------|-------------------------|
| `GET`      | `/games`            | Fetch all games         |
| `GET`      | `/games/{id}`       | Fetch a game by ID      |
| `POST`     | `/games`            | Create a new game       |
| `PUT`      | `/games/{id}`       | Update a game by ID     |
| `DELETE`   | `/games/{id}`       | Delete a game by ID     |


### Health Check
- **GET `/health`**
  - Responds with the server's health status.
  - **Response**:
    ```json
    {
        "status": "OK",
        "message": "Server is healthy and running"
    }
    ```

### Games
- **GET `/games`**
  - Fetches a list of all games.
  - **Response**:
    ```json
    [
        {
            "id": "12345",
            "title": "Game Title",
            "description": "Game description",
            "subject_id": "550e8400-e29b-41d4-a716-446655440000",
            "difficulty_level": 2,
            "created_at": "2024-01-01T00:00:00Z"
        }
    ]
    ```

- **POST `/games`**
  - Creates a new game.
  - **Request**:
    ```json
    {
        "title": "New Game",
        "description": "A fun new game",
        "subject_id": "550e8400-e29b-41d4-a716-446655440000",
        "difficulty_level": 2
    }
    ```
  - **Response**:
    ```json
    {
        "id": "67890",
        "title": "New Game",
        "description": "A fun new game",
        "subject_id": "550e8400-e29b-41d4-a716-446655440000",
        "difficulty_level": 2,
        "created_at": "2024-01-01T00:00:00Z"
    }
    ```

- **GET `/games/{id}`**
  - Fetches a specific game by its ID.
  - **Response**:
    ```json
    {
        "id": "12345",
        "title": "Game Title",
        "description": "Game description",
        "subject_id": "550e8400-e29b-41d4-a716-446655440000",
        "difficulty_level": 2,
        "created_at": "2024-01-01T00:00:00Z"
    }
    ```

  - **DELETE `/games/{id}`**
    - Deletes a specific game by its ID.
    - **Response**: NONE

---

## Installation

### Prerequisites
- [Go](https://golang.org/) 1.22 or higher
- Supabase account and project set up
- Environment variables configured in a `.env` file:
  ```env
  SUPABASE_URL=https://your-supabase-url.supabase.co
  SUPABASE_KEY=your-supabase-key
```

### Steps

1.  Clone the repository:

    ```bash
    git clone https://github.com/yourusername/ilang-backend.git
    cd ilang-backend
    ```

2.  Install dependencies:

    ```bash
    go mod tidy
    ```

3.  Run the server:

    ```bash
    go run main.go
    ```


* * *

Project Structure
-----------------

```bash
.
├── config/                # Configuration handling
│   └── config.go
├── handlers/              # HTTP route handlers
│   ├── GamesHandler.go
│   └── HealthHandler.go
├── models/                # Data models
│   └── Game.go
├── services/              # Business logic and Supabase integration
│   └── GameService.go
├── main.go                # Entry point for the application
├── .env                   # Environment variables (not included in Git)
└── go.mod                 # Go module file
```

* * *

Tech Stack
----------

*   **Language**: Go
*   **Database**: Supabase (PostgreSQL)
*   **HTTP Framework**: Built-in `net/http` package

* * *

Future Improvements
-------------------

*   Implement `PUT /games/{id}` and `DELETE /games/{id}` endpoints to complete CRUD operations.
*   Add authentication middleware for secured access.
*   Enhance error handling for better debugging and user feedback.
*   Expand the API to support additional features like game analytics and recommendations.

* * *

License
-------

This project is licensed under the MIT License. See the LICENSE file for details.

### **What’s Next?**
- Build more end points
- Build Frontend components
- Integrate frontend with backend and Supabase.
- Develop and test games.