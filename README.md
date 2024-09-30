# recipls

[![license](https://img.shields.io/github/license/s-hammon/recipls?style=for-the-badge)](https://github.com/s-hammon/recipls/blob/master/LICENSE)
[![report](https://goreportcard.com/badge/github.com/s-hammon/recipls?style=for-the-badge)](https://goreportcard.com/report/github.com/s-hammon/recipls)

## *Just the recipe, please*

Aren't you sick of searching online for a recipe, only to find that your top 10 results are webpages that lead with gigantic prose and contain endless ads which leave you starving to death? I sure am. This project is a back-end web service for storing and sharing not a bit more than just the recipe.

## Installation

This project requires an installation of [Go](https://go.dev/) 1.23+ and [PostgreSQL](https://www.postgresql.org/).

### Database

1. *(Optional)* Create a database--for example, on Linux:
    ```bash
    psql -U postgres -c "create database <dbname>;"
    ```
1. Create a username/password for the application to CRUD on the database.

    ```bash
    psql -U postgres -c "create user <username> with password '<password>';"
    psql -U postgres -c "grant all privileges on database <dbname> to <user>";
    ```

### Setting up the environment

1. Clone the repository

    ```bash
    git clone https://github.com/s-hammon/recipls.git
    cd recipls
    ```

1. Run `scripts/setup.sh`:

    2a. Create an environment variable containing the connection string that **recipls** will use to connect to your database. This must be passed to `setup.sh` as an arg:

    ```bash
    DATABASE_URL="host=<host> port=<port> user=<username> password=<password> dbname=<dbname> sslmode=disable"

    scripts/setup.sh $DATABASE_URL
    ```

    alternatively, you can pass the conn string directly into the script:

    ```bash
    scripts/setup.sh "host=<host> port=<port> user=<username> password=<password> dbname=<dbname> sslmode=disable"
    ```

    2b. This will create a starter `.env` file, as well as install `goose` and `sqlc`. The latter is only important for development, but `goose` is necessary to migrate the database. 

    2c. Verify that both are installed:

    ```bash
    goose --version     # goose version: v3.22.1
    sqlc version        # v1.27.0
    ```

1. Run `scripts/db.sh`:

    3a. Alternatively, on Linux you can run `make up` (reference the Makefile) to migrate your database.

1. Run `go run main.go` or `make run` (which will build a binary) to run the program with the default host (`localhost`) and port (`8080`). You can change either by using the following flags:

    ```bash
    --port <port>
    --host <host>
    
    # e.g.
    go run main.go --port 80 --host hunters-pub
    ```

## Services

### Authentication

|HTTP Method|URL|Parameter|Summary|
|:---|:---:|:---|:---|
|POST|`/v1/login`|`email`,`password`|Authentication with email and password. Issues a `refresh_token` and `access_token` (JWT)|
|POST|`/v1/refresh`|`refresh_token`|Use to issue a new JWT using the `refresh_token`j. If `refresh_token` is expires, will return a `401 Unauthorized` code|
|POST|`/v1/revoke`|`refresh_token`|Use to log a user out--or lock them out, if they're being naughty...|

### Users

|HTTP Method|URL|Parameter|Summary|
|:---|:---:|:---|:---|
|POST|`/v1/users`|`name`,`email`,`password`|Create a new user. Issues an `api_key` to use for `metrics` endpoints.|
|GET|`/v1/users`|`access_token`|Use to check user's session status. (may change this)|
|GET|`/v1/users/{id}`|`id`|Fetch a user by `id`. Data type: `uuid` (string)|

### Recipes
|HTTP Method|URL|Parameter|Summary|
|:---|:---:|:---|:---|
|GET|`/v1/recipes?user_id={}`|`user_id` (optional)|Fetch the list of recipes published by a `user_id`. If `user_id` is not specified, will return all recipes (limit 100, for now). TODO: add more query params.|
|GET|`/v1/recipes/{id}`|`id`|Fetch a recipe by `id`. Data type: `uuid` (string)|
|POST|`/v1/recipes`|`access_token`, `title`, `description`, `difficulty`, `ingredients`, `instructions`, `category`|Create a new recipe using JWT for session. Category must be in available list of categories (see `/v1/categories`)|
|PUT|`/v1/recipes/{id}`|`access_token`, `id`, `title`, `description`, `difficulty`, `ingredients`, `instructions`, `category`|Update a recipe by its `id`, using JWT for session. Again, category must be in available list of categories.|
|DELETE|`v1/recipes/{id}`|`access_token`, `id`|Delete a recipy by its `id`. Presently, only the author of the recipe can delete it (or the DBA lul).|

### Misc
|HTTP Method|URL|Parameter|Summary|
|:---|:---:|:---|:---|
|GET|`/v1/healthz`|none|Server health check--`200 StatusOK` if a-okay|
|GET|`/v1/categories`|none|Gets a list of available categories with which to categorize recipes.|
|GET|`/v1/metrics`|`api_token`|Presently available to a user with their `api_token`, returns a list of all users and the number of recipes they have published, as well as a list of all recipes and the number of steps in their instructions. (kind of boring, definitely want to expand this endpoint)|

## TODO

### Features

* Include photo of recipe results
* Add icons for required/special equipment
* Allow options to include common "gotchas" for certain instruction steps
* Enable user-specific RSS feeds

### Functionality

* Abstract service layer from handlers and repository (to enable testing of former two)
    * This in turn will allow functionality for other database engines

### CI/CD

* Add Github integration workflows
    * Test (after abstraction/more tests are made)
    * Gosec
    * Style/format/lint
* Streamline Docker image build