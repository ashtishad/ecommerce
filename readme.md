### Auth Challenge - Backend

##### Run `./start.sh` to download the dependencies and run the application.
* Make the Script Executable: You must give the script execute permissions before you can run it. Use the following command:
  `chmod +x start.sh`


To run the application, you have to define the environment variables, default values of the variables are defined inside `start.sh`

- SERVER_ADDRESS    `[IP Address of the machine]` : `localhost`
- SERVER_PORT       `[Port of the machine]` : `8000` `only 5000,5001,8000,80001 are allowed for google auth callback`
- DB_USER           `[Database username]` : `root`
- DB_PASSWD         `[Database password]`: `root`
- DB_ADDR           `[IP address of the database]` : `localhost`
- DB_PORT           `[Port of the database]` : `3306`
- DB_NAME           `[Name of the database]` : `users`

#### MySQL Database
Make the changes to your `start.sh` file for modifying default db configurations.
* `docker-compose.yml` file. This contains the database migrations scripts. You just need to bring the container up.
* `docker-compose down
  docker volume rm as_ti_mysqldata` to wipe up database and remove applied migrations.
   To start the docker container, run the `docker-compose up`.
* Run the application with `./start.sh` command from project root. or, if you want to run it from IDE, please set
  environment variables by executing command from start.sh on your terminal

#### Tools
  * Language used: GoLang
  * Database Used: MySQL
* Design       : Domain driven design
  * Libraries Used:
    * [Gorilla/Mux](https://github.com/gorilla/mux)
    * [Go SQL Driver](https://github.com/go-sql-driver/mysql)
    * [Google OAuth Library](golang.org/x/oauth2/google)

#### Project Structure
```

├── cmd
│   └── app
│       └── app.go                  <-- Define routes, logger setup, wire up handler, start server
│       ├── app_helpers.go          <-- Sanity check, validate input, writeResponse
│       ├── app_helpers_test.go     <-- Unit tests of Sanity check, validate input, writeResponse
│       └── db_connection.go        <-- MySQL db connection with dsn
│       └── handlers.go             <-- User handlers for app endpoints
│       └── google_auth_handlers.go <-- Google auth handlers for app endpoints
├── domain
│     └── user.go               <-- User struct based on database schema
│     ├── user_dto.go           <-- User level data with hiding sensitive fields
│     ├── user_repository.go    <-- Includes core repository interface
│     └── user_repository_db.go <-- Repository interface implementation with db
│     └── user_sql_queries.go   <-- SQL queries written seperately here
├── migrations                  <-- Database schema migrations scripts
├── docker-compose.yml          <-- Docker setup
├── start.sh                    <-- Builds app with exporting environment variables
├── readme.md                   <-- Self explanetory
├── main.go                     <-- Self explanetory

```


#### Data Flow (Hexagonal architecture)

    Incoming : Client --(JSON)-> REST Handlers --(DTO)-> Service --(Domain Object)-> RepositoryDB

    Outgoing : RepositoryDB --(Domain Object)-> Service --(DTO)-> REST Handlers --(JSON)-> Client

#### Testing Google Auth

```
1. Go to localhost:port/login page, login with your google account, select/enter your gmail account.
2. After redirecting to callback url, you will see a user created with sign_up_option= google.
3. Please use only 8000,8001,5000,5001 ports as SERVER_PORT in environemnt varibales..


As only backend is implemented for it, so I have skipped usual steps:
Frontend obtains Google OAuth token
Frontend sends OAuth token to backend
Backend verifies the token
```

#### Example Requests

###### Create a user

```
curl --location 'localhost:8000/user' \
--header 'Content-Type: application/json' \
--data-raw '{"email":"keanu_reeves@gmail.com","full_name":"Keanu Reeves","password":"secrpsswrd","phone":"1234567890","sign_up_option":"general"}'
```

###### Update a user

```
curl --location 'localhost:8000/user' \
--header 'Content-Type: application/json' \
--data-raw '{"email":"keanu_reeves@gmail.com","full_name":"John Wick","password":"secrpsswrd","phone":"1234567890","sign_up_option":"general"}'
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

#### Handled password with salt mechanism on [feature/pasword_hashing_salt](https://bitbucket.org/ashtishad/as_ti/src/c52d550196fc6169124340defb71af37e2a00e19/?at=feature%2Fpasword_hashing_salt) branch

* generates new salt + hashed-password on user create
* on updates it also update the salt value
* used database transactions for multi table update, insert.
