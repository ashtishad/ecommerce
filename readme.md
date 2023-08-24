## Ecommerce Microservice

### Microservice List



| Microservices      | Design Decisions            | Status             | Readme Link                                                               |
|--------------------|-----------------------------|--------------------|---------------------------------------------------------------------------|
| Users-API          | Go, RDBMS(PostgreSQL/MySQL) | Completed          | [Link](https://github.com/ashtishad/ecommerce/tree/main/users-api#readme) |
| Auth-API           | JWT, Google Auth            | Ongoing            |                                                                           |
| Product-API        |                             | Pending            |                                                                           |
| Order-API          |                             | Pending            |                                                                           |
| Cart-API           |                             | Pending            |                                                                           |
| Payment-API        |                             | Pending            |                                                                           |
| Review-API         |                             | Pending            |                                                                           |
| Customer-Care-API  |                             | Pending            |                                                                           |

### Design Decisions

1. Database per service pattern
  * Separate database for each service, one service won't communicate directly to another service's database. Why?
    * Separation of concerns (Each service to run independently).
    * Database schema/structure of another service that might change unexpectedly won't affect another.
    * There won't be a single point of failure would increase Site Reliability.
    * Some services might function more efficiently with different types of DB's (sql vs nosql).
    * Easy to scale, test, manage, maintain and audit.
    
2.How to exchange data between services?
   * Asynchronous Data Communication (Event Driven). Use Event Bus to exchange data(eg: Pache Kafka/RabbitMQ/NATS). Why Async Communication?
     * Zero dependency on other services.
     * No need to wait for other services to be ready.
     * Addition of new services is easy and service operations will be extremely fast.
  * Downside? - Data duplication.



###### General

* Software Architecture: Hexagonal(ports and adapters).
* Api Architecture Style: Restful API.
* Design Pattern: Domain Driven Design.
* Web Framework: Gin.
* Cloud: AWS.
* Containerization: Docker.
* CI/CD: GitHub Actions.
* Event Bus: Apache Kafka.
* Relational DB Preference: PostgreSQL/MySQL.
* Document Based/NoSQL DB: MongoDB.
* Cache preference: Redis.

###### Users-API

* RDBMS: MySQL/PostgreSQL (for structured data)
* Database Name: Users (will change it to Ecommerce)
* Cache: Redis
* Password Hashing: Salt

###### Auth-api

* Auth System: Oauth2 (JWT and Google Auth)
* DBMS: PostgreSQL (handling secure transactions)
* Database Name: Auth
* Cache: Redis (to quickly retrieve tokens and session information)

###### Product-API

* DBMS: PostgreSQL (for relational product attributes and categories).
* Database Name: Products.
* Cache: Redis (for caching popular products).

###### Order-API

* DBMS: PostgreSQL (to store complex order relationships).
* Database Name: Orders.
* Cache: Redis (for caching user cart details).

###### Cart-API

* DBMS: MongoDB (for flexible cart structures).
* Database Name: Carts.
* Cache: Redis (for quick access to cart data).

###### Payment-API

* DBMS: PostgreSQL (secure transaction handling).
* Database Name: Payments.
* Cache: Redis (for caching transaction data).

###### Review-API

* DBMS: MongoDB (to store varied review formats)
* Database Name: Reviews
* Cache: Redis (for caching popular reviews)

###### Customer-Care-API

* DBMS: PostgreSQL (structured customer care tickets)
* Database Name: CustomerCare
* Cache: Redis (for quickly accessing open tickets)

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


#### Project Structure
```
├── assets                          <-- For project root specific static assets.
├── lib                             <-- Common files shared between services(error library, logging, ginconfig etc)
├── users-api                       <-- Users API microservice.
├── auth-api                        <-- Auth API microservice.
├── docker-compose.yml              <-- Docker setup golangci.yml
├── golangci.yml                    <-- Config for golangci-lint. 
├── start.sh                        <-- Builds the whole app with exporting environment variables.
├── main.go                         <-- Responsible to start all server of this microservice
├── readme.md                       <-- Ecommerce Project Central Readme.

```

#### Data Flow (Hexagonal architecture)

    Incoming : Client --(JSON)-> REST Handlers --(DTO)-> Service --(Domain Object)-> RepositoryDB

    Outgoing : RepositoryDB --(Domain Object)-> Service --(DTO)-> REST Handlers --(JSON)-> Client


#### Example Requests(Users-API routes)

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




###### Hexagonal Architecture

![hexagonal_architecture.png](assets%2Fimages%2Fhexagonal_architecture.png)
