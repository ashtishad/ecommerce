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

#### Example Response

##### Get All categories with hierarchy(level by level with all sub-categories)

```

[
    {
        "categoryUuid": "bd11d903-7549-42b2-bea6-dd8a7cb8821e",
        "parentCategoryUuid": "",
        "name": "Sound Equipment",
        "description": "All kind of sound equipments",
        "status": "active",
        "createdAt": "2023-09-13T04:55:41.174441Z",
        "updatedAt": "2023-09-13T04:55:41.174441Z",
        "level": 1,
        "subcategories": [
            {
                "categoryUuid": "969517ac-dba1-4e07-b3c1-3bf4233b84de",
                "parentCategoryUuid": "bd11d903-7549-42b2-bea6-dd8a7cb8821e",
                "name": "TWS",
                "description": "Budget TWS",
                "status": "active",
                "createdAt": "2023-09-13T04:56:44.892263Z",
                "updatedAt": "2023-09-13T04:56:44.892263Z",
                "level": 2,
                "subcategories": [
                    {
                        "categoryUuid": "fb3473ac-6b51-4671-b3da-e38379a21768",
                        "parentCategoryUuid": "969517ac-dba1-4e07-b3c1-3bf4233b84de",
                        "name": "Budget TWS",
                        "description": "Budget TWS",
                        "status": "active",
                        "createdAt": "2023-09-13T04:57:33.617043Z",
                        "updatedAt": "2023-09-13T04:57:33.617043Z",
                        "level": 3,
                        "subcategories": [
                            {
                                "categoryUuid": "0db872ce-1be4-4c21-87e2-7f9194e27089",
                                "parentCategoryUuid": "fb3473ac-6b51-4671-b3da-e38379a21768",
                                "name": "L3",
                                "description": "Budget TWS",
                                "status": "active",
                                "createdAt": "2023-09-13T05:11:41.611501Z",
                                "updatedAt": "2023-09-13T05:11:41.611501Z",
                                "level": 4,
                                "subcategories": []
                            }
                        ]
                    }
                ]
            }
        ]
    },
    {
        "categoryUuid": "e085c298-35b0-4b05-bcc1-a24d4fff4794",
        "parentCategoryUuid": "",
        "name": "Phone",
        "description": "All kind of sound equipments",
        "status": "active",
        "createdAt": "2023-09-13T05:25:12.766975Z",
        "updatedAt": "2023-09-13T05:25:12.766975Z",
        "level": 1,
        "subcategories": [
            {
                "categoryUuid": "821c135a-453f-4f8c-af37-6e78d5a7bdd6",
                "parentCategoryUuid": "e085c298-35b0-4b05-bcc1-a24d4fff4794",
                "name": "Apple",
                "description": "Budget TWS",
                "status": "active",
                "createdAt": "2023-09-13T05:25:43.514098Z",
                "updatedAt": "2023-09-13T05:25:43.514098Z",
                "level": 2,
                "subcategories": [
                    {
                        "categoryUuid": "7e37c688-effa-40ca-be59-e61a54d4abbd",
                        "parentCategoryUuid": "821c135a-453f-4f8c-af37-6e78d5a7bdd6",
                        "name": "Refurbished",
                        "description": "Budget TWS",
                        "status": "active",
                        "createdAt": "2023-09-13T05:26:07.057982Z",
                        "updatedAt": "2023-09-13T05:26:07.057982Z",
                        "level": 3,
                        "subcategories": []
                    },
                    {
                        "categoryUuid": "b88ecd5c-d58f-4077-84a0-cfa8fada134d",
                        "parentCategoryUuid": "821c135a-453f-4f8c-af37-6e78d5a7bdd6",
                        "name": "Original",
                        "description": "Budget TWS",
                        "status": "active",
                        "createdAt": "2023-09-13T05:26:11.955702Z",
                        "updatedAt": "2023-09-13T05:26:11.955702Z",
                        "level": 3,
                        "subcategories": [
                            {
                                "categoryUuid": "dc689a06-c859-41f4-8b44-d29787351c67",
                                "parentCategoryUuid": "b88ecd5c-d58f-4077-84a0-cfa8fada134d",
                                "name": "Before X",
                                "description": "Budget TWS",
                                "status": "active",
                                "createdAt": "2023-09-13T05:26:29.213884Z",
                                "updatedAt": "2023-09-13T05:26:29.213884Z",
                                "level": 4,
                                "subcategories": []
                            },
                            {
                                "categoryUuid": "bc15f1a4-159e-4e78-a2fb-fe8ab8aa6faa",
                                "parentCategoryUuid": "b88ecd5c-d58f-4077-84a0-cfa8fada134d",
                                "name": "After X",
                                "description": "Budget TWS",
                                "status": "active",
                                "createdAt": "2023-09-13T05:26:33.654446Z",
                                "updatedAt": "2023-09-13T05:26:33.654446Z",
                                "level": 4,
                                "subcategories": []
                            }
                        ]
                    }
                ]
            }
        ]
    },
    {
        "categoryUuid": "eef83251-6e74-4039-8399-2337257a99bf",
        "parentCategoryUuid": "",
        "name": "Wearables",
        "description": "All kind of sound Wearables",
        "status": "active",
        "createdAt": "2023-09-13T05:31:57.571796Z",
        "updatedAt": "2023-09-13T05:31:57.571796Z",
        "level": 1,
        "subcategories": [
            {
                "categoryUuid": "51d36a1c-3ef5-4fda-8a5f-90a9054de0e2",
                "parentCategoryUuid": "eef83251-6e74-4039-8399-2337257a99bf",
                "name": "Smart Watch",
                "description": "Budget TWS",
                "status": "active",
                "createdAt": "2023-09-13T05:32:11.710595Z",
                "updatedAt": "2023-09-13T05:32:11.710595Z",
                "level": 2,
                "subcategories": [
                    {
                        "categoryUuid": "09e8f5d0-9b6d-471d-b680-22ded152ac0f",
                        "parentCategoryUuid": "51d36a1c-3ef5-4fda-8a5f-90a9054de0e2",
                        "name": "Level 2",
                        "description": "Budget TWS",
                        "status": "active",
                        "createdAt": "2023-09-13T05:32:51.438349Z",
                        "updatedAt": "2023-09-13T05:32:51.438349Z",
                        "level": 3,
                        "subcategories": [
                            {
                                "categoryUuid": "72d40b1d-27e0-425d-b988-14db91e4a61b",
                                "parentCategoryUuid": "09e8f5d0-9b6d-471d-b680-22ded152ac0f",
                                "name": "Level 3-1",
                                "description": "Budget TWS",
                                "status": "active",
                                "createdAt": "2023-09-13T05:33:08.933988Z",
                                "updatedAt": "2023-09-13T05:33:08.933988Z",
                                "level": 4,
                                "subcategories": []
                            },
                            {
                                "categoryUuid": "8ad5cf6a-e4c7-4f7d-8ad7-7983316ad2da",
                                "parentCategoryUuid": "09e8f5d0-9b6d-471d-b680-22ded152ac0f",
                                "name": "Level 3-2",
                                "description": "Budget TWS",
                                "status": "active",
                                "createdAt": "2023-09-13T05:33:12.376908Z",
                                "updatedAt": "2023-09-13T05:33:12.376908Z",
                                "level": 4,
                                "subcategories": [
                                    {
                                        "categoryUuid": "6afaee0a-0d25-4360-bf14-aa5351e561d1",
                                        "parentCategoryUuid": "8ad5cf6a-e4c7-4f7d-8ad7-7983316ad2da",
                                        "name": "Level 4-1",
                                        "description": "Budget TWS",
                                        "status": "active",
                                        "createdAt": "2023-09-13T05:33:29.127599Z",
                                        "updatedAt": "2023-09-13T05:33:29.127599Z",
                                        "level": 5,
                                        "subcategories": []
                                    }
                                ]
                            }
                        ]
                    }
                ]
            }
        ]
    }
]

```
