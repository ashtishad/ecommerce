## Ecommerce Microservice

#### Environment Setup

###### Clone using ssh protocol `git clone git@github.com:ashtishad/ecommerce.git`

To run the application, you have to define the environment variables, default values of the variables are defined inside `start.sh`

- SERVER_ADDRESS    `[IP Address of the machine]` : `localhost`
- SERVER_PORT       `[Port of the machine]` : `8000` 
- DB_USER           `[Database username]` : `root`
- DB_PASSWD         `[Database password]`: `root`
- DB_ADDR           `[IP address of the database]` : `localhost`
- DB_PORT           `[Port of the database]` : `3306`
- DB_NAME           `[Name of the database]` : `users`

###### MySQL Database Setup
* Make the changes to your `start.sh` file for modifying default db configurations.
* `docker-compose.yml` file. This contains the database migrations scripts. You just need to bring the container up.
* `docker-compose down
  docker volume rm ecommerce_mysqldata` to wipe up a database and remove applied migrations.
  To start the docker container, run the `docker-compose up`.

###### Run the application
* Run the application with `./start.sh` command from project root. or, if you want to run it from IDE, please set
  environment variables by executing command from start.sh on your terminal.
* (optional) Make the Script Executable: You must give the script execute permissions before you can run it. Use the following command:
    `chmod +x start.sh`

#### Tools Used

* Language used: GoLang
* Database Used: MySQL
* Design       : Domain driven design
  * Libraries Used:
    * [Gin](https://github.com/gin-gonic/gin)
    * [Go SQL Driver](https://github.com/go-sql-driver/mysql)
    * [Google OAuth Library](golang.org/x/oauth2/google)

#### Project Structure
```

├── users-api                       <-- Users API for this microservice - [Readme](https://github.com/ashtishad/ecommerce/users-api/readme.md)
├── docker-compose.yml              <-- Docker setup
├── start.sh                        <-- Builds the whole app with exporting environment variables
├── readme.md                       <-- Ecommerce Project Central Readme
├── main.go                         <-- Responsible to start all server of this microservice

```

##### USERS_API

#### Data Flow (Hexagonal architecture)

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


###### Hexagonal Architecture

![hexagonal_architecture.png](assets%2Fimages%2Fhexagonal_architecture.png)
