### Auth Challenge

##### Run `./start.sh` to download the dependencies and run the application.
* Make the Script Executable: You must give the script execute permissions before you can run it. Use the following command:
  `chmod +x start.sh`


To run the application, you have to define the environment variables, default values of the variables are defined inside `start.sh`

- SERVER_ADDRESS    `[IP Address of the machine]` : `localhost`
- SERVER_PORT       `[Port of the machine]` : `8000`
- DB_USER           `[Database username]` : `root`
- DB_PASSWD         `[Database password]`: `root`
- DB_ADDR           `[IP address of the database]` : `localhost`
- DB_PORT           `[Port of the database]` : `3306`
- DB_NAME           `[Name of the database]` : `users`

#### MySQL Database
Make the changes to your `start.sh` file for modifying default db configurations.
* `docker-compose.yml` file. This contains the database migrations scripts. You just need to bring the container up.
* `ocker-compose down
  docker volume rm as_ti_mysqldata` to wipe up database and remove applied migrations.
   To start the docker container, run the `docker-compose up`.


#### Technical Requirements
  * Language used: GoLang
  * Database Used: MySQL
  * Libraries Used:
    * [Gorilla/Mux](https://github.com/gorilla/mux)
    * [Go SQL Driver](https://github.com/go-sql-driver/mysql)

#### Project Structure
```

├── cmd
│   └── app
│       └── app.go              <-- Define routes, logger setup, wire up handler, start server
│       ├── app_helpers.go      <-- Sanity check, validate input, writeResponse
│       ├── app_helpers_test.go <-- Unit tests of Sanity check, validate input, writeResponse
│       └── db_connection.go    <-- MySQL db connection with dsn
│       └── handlers.go         <-- Handlers for app routes
├── domain
│     └── user.go               <-- User struct based on database schema
│     ├── user_dto.go           <-- User level data with hiding sensitive fields
│     ├── user_repository.go    <-- Includes core repository interface
│     └── user_repository_db.go <-- Repository interface implementation with db
├── migrations                  <-- Database schema migrations scripts
├── docker-compose.yml          <-- Docker setup
├── start.sh                    <-- Builds app with exporting environment variables
├── readme.md                   <-- Self explanetory
├── main.go                     <-- Self explanetory

```


#### Data Flow (Hexagonal architecture)

    Incoming : Client --(JSON)-> REST Handlers --(DTO)-> Service --(Domain Object)-> RepositoryDB

    Outgoing : RepositoryDB --(Domain Object)-> Service --(DTO)-> REST Handlers --(JSON)-> Client

