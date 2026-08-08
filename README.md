# Tasks API

A simple RESTful task management API built with Go, Gin, PostgreSQL, and JWT authentication.

Each user can register, log in, and manage only their own tasks.

## Features

* User registration
* User login with JWT authentication
* Password hashing with bcrypt
* Create, read, update, and delete tasks
* User-specific task access
* PostgreSQL database
* Layered architecture
* Unit tests with mocks
* Repository tests with `sqlmock`

## Technologies

* Go
* Gin
* PostgreSQL
* JWT
* bcrypt
* `go-sqlmock`

## Project Structure

```text
tasks/
├── internal/
│   ├── domain/
│   ├── handler/
│   ├── repository/
│   │   └── postgres/
│   └── service/
├── pkg/
│   └── db/
├── main.go
├── go.mod
└── go.sum
```

The project follows this flow:

```text
Handler → Service → Repository → PostgreSQL
```

## API Endpoints

### Authentication

| Method | Endpoint       | Description              |
| ------ | -------------- | ------------------------ |
| POST   | `/user/signup` | Register a new user      |
| POST   | `/user/login`  | Log in and receive a JWT |

### Tasks

The following endpoints require a Bearer token.

| Method | Endpoint     | Description                                 |
| ------ | ------------ | ------------------------------------------- |
| GET    | `/tasks`     | Get all tasks belonging to the current user |
| GET    | `/tasks/:id` | Get one task                                |
| POST   | `/tasks`     | Create a task                               |
| PUT    | `/tasks/:id` | Update a task                               |
| DELETE | `/tasks/:id` | Delete a task                               |

## Authentication

Send the JWT in the `Authorization` header:

```http
Authorization: Bearer YOUR_TOKEN
```

## Example Requests

### Sign Up

```json
{
  "username": "amin",
  "password": "123456"
}
```

### Login

```json
{
  "username": "amin",
  "password": "123456"
}
```

### Create Task

```json
{
  "text": "Learn Go",
  "desc": "Practice building REST APIs",
  "status": "pending"
}
```

## Running the Project

Clone the repository:

```bash
git clone https://github.com/AminOp43/tasks.git
cd tasks
```

Install dependencies:

```bash
go mod download
```

Create a PostgreSQL database and configure the database connection used in:

```text
pkg/db/db.go
```

Set a JWT secret in your environment:

### PowerShell

```powershell
$env:JWT_SECRET="your-secure-secret"
```

Run the application:

```bash
go run main.go
```

## Running Tests

Run all tests:

```bash
go test ./... -count=1
```

## What I Practiced

This project was created to practise:

* Building REST APIs with Go and Gin
* Working with PostgreSQL
* Authentication with JWT
* Password hashing
* Layered backend architecture
* Interfaces and dependency injection
* Mocking and unit testing
* Git and GitHub workflow
