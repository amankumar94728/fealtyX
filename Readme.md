# FealtyX Student API

FealtyX Student API is a RESTful API built with Go that manages student information and generates summaries using Ollama.

## Features

- CRUD operations for student records
- In-memory data storage
- Student summary generation using Ollama
- Concurrent-safe operations
- Input validation

## Prerequisites

- Go 1.16 or higher
- Ollama installed and running locally

## Installation

1. Clone the repository:

```bash
git clone https://github.com/amankumar94728/fealtyx-student-api.git
cd fealtyx-student-api
```

2. Install dependencies:

```bash
go mod tidy
```

## Usage

1. Start the server:

```bash
go run main.go
```

The server will start on `http://localhost:8080`.

2. Use the following endpoints:

- Create a student: `POST /students`
- Get all students: `GET /students`
- Get a student by ID: `GET /students/{id}`
- Update a student: `PUT /students/{id}`
- Delete a student: `DELETE /students/{id}`
- Generate a student summary: `GET /students/{id}/summary`

## API Examples

### Create a student

```bash
curl -X POST -H "Content-Type: application/json" -d '{"name":"John Doe","age":20,"email":"john@example.com"}' http://localhost:8080/students
```

### Get all students

```bash
curl http://localhost:8080/students
```

### Get a student by ID

```bash
curl http://localhost:8080/students/1
```

### Update a student

```bash
curl -X PUT -H "Content-Type: application/json" -d '{"name":"John Doe","age":21,"email":"john.doe@example.com"}' http://localhost:8080/students/1
```

### Delete a student

```bash
curl -X DELETE http://localhost:8080/students/1
```

### Generate a summary for a student

```bash
curl http://localhost:8080/students/1/summary