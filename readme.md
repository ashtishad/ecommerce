## Ecommerce Microservice

### Microservice List

| Microservices     | Design Decisions      | Status    | Readme Link                                                                 |
|-------------------|-----------------------|-----------|-----------------------------------------------------------------------------|
| Users-API         | Go, RDBMS(PostgreSQL) | Completed | [Link](https://github.com/ashtishad/ecommerce/tree/main/users-api#readme)   |
| Product-API       | Go, RDBMS(PostgreSQL) | Ongoing   | [Link](https://github.com/ashtishad/ecommerce/tree/main/product-api#readme) |
| Auth-API          |                       | Pending   |                                                                             |
| Order-API         |                       | Pending   |                                                                             |
| Cart-API          |                       | Pending   |                                                                             |
| Payment-API       |                       | Pending   |                                                                             |
| Review-API        |                       | Pending   |                                                                             |
| Customer-Care-API |                       | Pending   |                                                                             |

### Design Decisions(V1)

###### Frameworks/Concepts

* Software Architecture: Hexagonal(ports and adapters).
* Api Architecture Style: Restful API.
* Design Pattern: Domain Driven Design.
* Web Framework: Gin.
* Cloud: AWS.
* Containerization: Docker.
* CI/CD: GitHub Actions.
* Event Bus: Apache Kafka.
* Relational DB Preference: PostgreSQL.
* Document Based/NoSQL DB: MongoDB.
* Cache preference: Redis.


#### Environment Setup

###### Clone using ssh protocol `git clone git@github.com:ashtishad/ecommerce.git`

Change environment variables in Makefile: Set values in Makefile stored in project root.

- SERVER_ADDRESS    `[IP Address of the machine]` : `localhost`
- USER_API_PORT      `[Port of the user api]` : `8000`
- PRODUCT_API_PORT   `[Port of the product api]` : `8001`
- DB_USER           `[Database username]` : `postgres`
- DB_PASSWD         `[Database password]`: `potgres`
- DB_ADDR           `[IP address of the database]` : `localhost`
- DB_PORT           `[Port of the database]` : `5432`
- DB_NAME           `[Name of the database]` : `ecommerce`
- POSTGRESQL_URL    `[Postgres DB Connection URL for golang-migrate cli]`

###### Postgres Database Setup

* Run docker compose: Bring the container up with `docker compose up`. Configurations are in `compose.yaml` file.
* (optional) Remove databases and volumes:
  ```
  docker compose down
  docker volume rm ecommerce_postgresdata
  ```

###### Run the application

* Run the application with `make run` command from project root. or, if you want to run it from IDE, please set
  environment variables by executing commands mentioned in Makefile on your terminal.


#### Project Structure
```
├── .github/workflows        <-- Github CI workflows(Build, Test, Lint).
├── assets                   <-- For project root specific static assets.
├── config                   <-- Database initialization on docker compose.
├── db/migrations            <-- Postgres DB migrations scripts for golang-migrate.
├── users-api                <-- Users API microservice.
├── product-api              <-- Auth API microservice.
├── lib                      <-- Common setup, configs used across all services.
├── compose.yaml             <-- Docker services setup(databases)
├── golangci.yml             <-- Config for golangci-lint. 
├── Makefile                 <-- Builds the whole app with exporting environment variables.
├── main.go                  <-- Start all server concurrently, init logger, init db, env port check, graceful shutdown.
├── readme.md                <-- Ecommerce Project Central Readme.

```

#### Data Flow (Hexagonal architecture)

    Incoming : Client --(JSON)-> REST Handlers --(DTO)-> Service --(Domain Object)-> RepositoryDB

    Outgoing : RepositoryDB --(Domain Object)-> Service --(DTO)-> REST Handlers --(JSON)-> Client



###### Hexagonal Architecture

![hexagonal_architecture.png](assets%2Fimages%2Fhexagonal_architecture.png)

### Design Decisions(V2)

1. Database per service pattern

* Separate database for each service, one service won't communicate directly to another service's database. Why?
  * Separation of concerns (Each service to run independently).
  * Database schema/structure of another service that might change unexpectedly won't affect another.
  * There won't be a single point of failure would increase Site Reliability.
  * Some services might function more efficiently with different types of DB's (sql vs nosql).
  * Easy to scale, test, manage, maintain and audit.

2.How to exchange data between services?

* Asynchronous Data Communication (Event Driven). Use Event Bus to exchange data(eg: Pache Kafka/RabbitMQ/NATS). Why
  Async Communication?
  * Zero dependency on other services.
  * No need to wait for other services to be ready.
  * Addition of new services is easy and service operations will be extremely fast.
* Downside? - Data duplication.

