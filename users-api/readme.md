## UsersAPI

#### Tools

* Language used: GoLang.
* Database Used: PostgreSQL.
* Design       : Domain driven design.
* Logging      : Structured log with [slog](https://pkg.go.dev/log/slog#section-documentation)
* 3rd Party Libraries Used:
  * Web Framework: [gin](https://github.com/gin-gonic/gin)
  * Postgres Driver : [pq](https://pkg.go.dev/github.com/lib/pq#section-readme)
  * DB migrations: [migrate](https://github.com/golang-migrate/migrate)
  * Mock Testing: [go-sqlmock](https://github.com/DATA-DOG/go-sqlmock)
  * Mock Testing: [testify](https://github.com/stretchr/testify/)
  * Fake data generate: [gofakeit](https://github.com/brianvoe/gofakeit)

#### Project Structure(Users-API)
```

├── cmd
│   └── app
│       └── app.go                        <-- Define routes, logger setup, wire up handler, start users-api server
│       ├── app_helpers.go                <-- Sanity check.
│       └── handlers.go                   <-- User handlers for app endpoints
│       └── handlers_test.go              <-- User handlers tests only for 200 OK.
├── database
│   └── postgres_db_conn.go               <-- Postgresql DB Connection config.
│   └── postgres_db_conn_test.go          <-- Connection string making test.
│   └── generate_users.go                 <-- Generate n(1000) users in database.
├── internal
│   └── domain
│       └── user.go                       <-- User struct based on database schema.
│       ├── user_dto.go                   <-- User level data with hiding sensitive fields.
│       ├── user_repository.go            <-- Includes core repository interface.
│       └── user_repository_db.go         <-- Repository interface implementation with db.
│       └── user_repository_db_test.go    <-- Mock tests for all repository db methods.
│       └── user_sql_queries.go           <-- SQL queries written seperately here.
│   └── service
│       └── user_service.go               <-- Generate salt,hash pass, covert dto to domain and vice versa.
│       └── service_helpers.go.go         <-- Included user input validation.
│       └── mock_user_service.go          <-- Mocked user services for handlers test.
├── pkg
│   └── constants
│       └── constants.go                  <-- Included constants for input validation regexes, database enum values.
│   └── hashpassword
│       └── hashpassword.go               <-- Generate random salt and hashpassword with it.
```

#### Design Decisions

* Handle password with salt mechanism
* Database transactions for multi table operations
* Mock testing: Used go-sql-mock for mock db setup and testify for require method.

#### Data Flow

    Incoming : Client --(JSON)-> REST Handlers --(DTO)-> Service --(Domain Object)-> RepositoryDB

    Outgoing : RepositoryDB --(Domain Object)-> Service --(DTO)-> REST Handlers --(JSON)-> Client


#### Example Requests

###### Create a user

```
curl --location 'localhost:8000/users' \
--header 'Content-Type: application/json' \
--data-raw '{"email":"keanu_reeves@gmail.com","full_name":"Keanu Reeves","password":"secrpsswrd","phone":"1234567890","sign_up_option":"general"}'
```

###### Update a user

```
curl --location --request PUT 'localhost:8000/users/{user_id}' \
--header 'Content-Type: application/json' \
--data-raw '{
	"email": "keanu_reeves@gmail.com",
	"full_name": "John Wick",
    "phone": "1234567890"
	
}

```

##### Get All Users (Paginated, Optional query params "status", "fromID", "pageSize", "timezone" and "signUpOption"

1. Populate users from users-api/cmd/app.go
   uncomment this line

```
// database.GenerateUsers(conn, l, 1000)

```
   it will generate 1000 users when app builds(so, consider comment out it again after first run)
2. Used cursor based pagination strategy with UserID(e.g: 0,1,5000) as cursor.
3. Timezone region example "UTC", "Asia/Dhaka".

Request examples with or without Params

```
curl --location 'localhost:8000/users'

curl --location 'localhost:8000/users?signUpOption=google&status=active&pageSize=20'

curl --location 'localhost:8000/users?signUpOption=google&status=active&pageSize=20&timezone=Asia%2FHarbin'

```
