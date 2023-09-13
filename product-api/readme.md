## Product-API

#### Tools

* Language used: GoLang.
* Database Used: PostgreSQL.
* Design       : Domain driven design.
* Logging      : Structured log with [slog](https://pkg.go.dev/log/slog#section-documentation)
* Web Framework: [gin](https://github.com/gin-gonic/gin)

#### Project Structure(Users-API)

```

├── cmd
│   └── app
│       └── app.go                          <-- wire up the handlers, route setup, start product-api server
│       └── category_handlers.go            <-- Handlers for categories and sub-categories endpoints.
│       └── product_handlers.go             <-- Handlers for product endpoints
├── internal
│   └── domain
│       └── category.go                     <-- Category struct based on database schema.
│       ├── category_dto.go                 <-- Hiding sensitive fields here.
│       ├── category_repo_queries.go        <-- Includes sql queries.
│       └── category_repository.go          <-- Includes core repository interface
│       └── category_repository_db.go       <-- Repository interface implementation with db.
│   └── service
│       └── category_service.go             <-- Validate request, convert dto to domain and vice versa.
│       └── service_helpers.go              <-- Included user input validation.
│       └── service_helpers_test.go         <-- Tests for validation methods.

```

#### Data Flow

    Incoming : Client --(JSON)-> REST Handlers --(DTO)-> Service --(Domain Object)-> RepositoryDB

    Outgoing : RepositoryDB --(Domain Object)-> Service --(DTO)-> REST Handlers --(JSON)-> Client

#### Example Requests

###### Create a Category

POST: /categories

1. DB transaction
2. check category name uniqueness
3. create category

```

curl --location 'localhost:8001/categories' \
--header 'Content-Type: application/json' \
--data '{
	"name": "Wearables",
    "description": "All kind of sound Wearables"
}'

```

###### Create a sub-category

POST: /categories//:category_id/subcategories
Note: category_id is the uuid of parent category

1. DB transaction
2. check category name uniqueness
3. validate uuid, and get's id
4. calculate subcategory level(depth) and insert ancestor descendant relationships

```

curl --location 'localhost:8001/categories/bd11d903-7549-42b2-bea6-dd8a7cb8821e/subcategories' \
--data '{
	"name": "Smart Watch",
    "description": "Any descriptions"
}'

```

##### Get All categories with hierarchy(level by level with all sub-categories)

1. RECURSIVE CategoryTree query
2. Then scan rows and build hierarchy
3. Return all hierarchy at once

```

curl --location 'localhost:8001/categories'

```
