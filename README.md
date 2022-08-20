# Introduction
I use Clean Architecture to structure the code because I think this way source code will cleaner
and easier maintenance when scale.

The technologies is used
- Go (Golang)
- Echo framwork
- MySQL
- Docker
- [migrate-sql](https://github.com/rubenv/sql-migrate)

# Prepare
- Install Go, Docker, MySQL client, make
- Install migrate-sql to run sql-migrate in order to apply migration to database

```bash
go get -v github.com/rubenv/sql-migrate/...
```

# Run
## Local devevelop
- Copy `config/env.template.yml` and rename to `config/env.yml`
- Start mysql docker server

```bash
docker compose -f ./docker-compose.local.yml up -d
```

- Wait for mysql docker server is started completely, run below command to apply migration and seed data

```bash
make seed
```

- Start api server local

```bash
go run cmd/main.go
```

# Trying some tests
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
