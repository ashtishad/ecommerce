## Users-API

#### Tools Used

* Language used: GoLang
* Database Used: MySQL
* Design       : Domain driven design
  * Libraries Used:
    * [Gin](https://github.com/gin-gonic/gin)
    * [Go SQL Driver](https://github.com/go-sql-driver/mysql)

#### Project Structure(Users-API)
```

├── cmd
│   └── app
│       └── app.go                    <-- Define routes, logger setup, wire up handler, start users-api server
│       ├── app_helpers.go            <-- Sanity check.
│       ├── app_helpers_test.go       <-- Unit tests of Sanity check.
│       └── handlers.go               <-- User handlers for app endpoints
├── database
│   └── migrations                    <-- Database schema migrations scripts of RDBMS.
│   └── mysql_db_connection.go        <-- MySQL DB Connection config.
├── internal
│   └── domain
│       └── user.go                   <-- User struct based on database schema.
│       ├── user_dto.go               <-- User level data with hiding sensitive fields.
│       ├── user_repository.go        <-- Includes core user repository interface.
│       └── user_repository_db.go     <-- Repository interface implementation with db.
│       └── user_sql_queries.go       <-- SQL queries written seperately here.
│   └── service
│       └── user_service.go           <-- Generate salt,hash pass, convert dto to domain and vice versa.
│       └── service_helpers.go.go     <-- User input validation.
├── pkg
│   └── hashpassword
│       └── hashpassword.go           <-- Generate random salt and hashpassword with it.
```


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
'
```

###### Check Existing User by Email and Password

```

curl --location 'localhost:8000/existing-user' \
--header 'Content-Type: application/json' \
--data-raw '{
"email": "keanu_reeves@gmail.com",
"password": "seepasword"
}'

```

#### Design Decisions

###### 1. Handle password with salt mechanism

* generates new salt + hashed-password on user creation.
* on update, it also updates the salt value corresponding to user_id.
* used database transactions for multi table update, insert.

###### 2. Database transactions for multi table operations
