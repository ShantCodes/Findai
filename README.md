# Findai

this is the backend service of an ai idea, this project's goal is to only gather data from 
users for ai training goals, feel free to collabrate, reach via email for more details about the project

## Getting Started

To get this project up and running, follow these steps:

### 1. Configuration

Before running the application, you need to set up your configuration.
Copy the `config.example.yml` (you will need to create this based on `config.yml` structure) to `config.yml` and fill in the necessary details.

### 2. Run Database Migrations (Optional)

If your project requires a database, you can run migrations to set up your schema.

To apply all pending migrations:
```bash
go run cmd/migrate/main.go up
```

To create a new migration file:
```bash
go run cmd/migrate/main.go new <migration_name>
```

### 3. Run the Main Application

To start the main API server:
```bash
go run cmd/app/main.go
```
The application will usually run on the port specified in your `config.yml` (defaulting to 8080). You can test it by navigating to `/ping` (e.g., `http://localhost:8080/ping`) in your browser or with a tool like `curl`.

### 4. Run the Worker (Optional)

If your application includes background workers, you can start the worker process:
```bash
go run cmd/worker/main.go
```


### FOR FIRST TIME RUN:
```
$ cd findai
$ cp config.example.yml config.yml
$ sudo docker-compose up -d
$ go get
$ go run cmd/migration/main.go up
$ go run cmd/app/main

```


