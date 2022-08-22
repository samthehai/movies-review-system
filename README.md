# Introduction
There some thoughts about system design as below:
- Use Clean Architecture to structure the code because I think this way source code will cleaner
and easier maintenance when scale.
- Separate code structure into modules such as movie, user for easier to separate to microservices when want to scale
- About full-text-search, I currently use MySQL FULLTEXT function for simplicity. In future if performance problem happend, we can refactor it use ElasticSearch, ...
- Use JWT to implement accesstoken

The technologies is used
- Go (Golang) version 18
- Echo framework
- MySQL
- Docker
- [migrate-sql](https://github.com/rubenv/sql-migrate)
- [go-sqlmock](https://github.com/DATA-DOG/go-sqlmock)
- [mockgen](https://github.com/golang/mock)

# Prepare
- Install Go, Docker, MySQL client, make
- Install mockgen, sql-migrate
- Install golangci-lint to check lint local: https://golangci-lint.run/usage/install/#local-installation

```bash
go get -v github.com/rubenv/sql-migrate/...
go install github.com/golang/mock/mockgen@v1.6.0
```

# Commands
Some convienience commands can be found at Makefile, so please refer to this file for details.

# Run
## Local devevelop
- Copy `config/env.template.yml` and rename to `config/env.yml`
- Start mysql docker server

```bash
docker compose -f ./docker-compose.yml up
```

- Wait for mysql docker server is started completely, open a new terminal and run below command to apply migration and seed data

```bash
make seed
```

- Start api server local

```bash
go run cmd/main.go
```

# Trying some tests
## Use curl command
- Login

```
curl -X POST http://localhost:5000/api/v1/users/login -H "Content-Type: application/json" -d '{"email":"testuser@gmail.com","password":"secret"}'
```

- Get list of top movies

```
curl -X GET http://localhost:5000/api/v1/movies
```

- Full text search

```
curl -X GET http://localhost:5000/api/v1/movies?search=gravida
```

- Favorite a movie

  - First login to get the accesstoken

```
curl -X POST http://localhost:5000/api/v1/users/login -H "Content-Type: application/json" -d '{"email":"testuser@gmail.com","password":"secret"}'
```
  - Second run below command (don't forget to fill the accesstoken)

```
curl -X POST http://localhost:5000/api/v1/favorites/1 \
         -H "Content-Type: application/json" \
         -H "Authorization: Bearer <accesstoken which is got from login api>"
```

- Get list previously marked favorite movies

```
curl -X GET http://localhost:5000/api/v1/favorites \
         -H "Content-Type: application/json" \
         -H "Authorization: Bearer <accesstoken which is got from login api>"
```
## Use Swagger
Access http://localhost:5000/swagger/index.html in order to access Swagger

## Run unittests
I have created some testcases for movie module, those testcases can be run by execute below command:
```
make test
```

or run directly: `go test ./...`
