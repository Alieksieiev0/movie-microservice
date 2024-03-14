# Movie Microservice

[![GoDoc](https://pkg.go.dev/badge/github.com/Alieksieiev0/movie-microservice?status.svg)](https://godoc.org/github.com/Alieksieiev0/movie-microservice)

Movie Microservice is a small project with microservice architecture, that allows user to perform CRUD operations. It is possible to:

 - Get Movie (GET /movies/{id})
 - Get All Movies (GET /movies)
 - Create Movie (POST /movies)
 - Update Movie (PUT /movies)
 - Delete Movie (DELETE /movies)

Postresql(with pgx) was used to store all data, and all of the operations includes calls to the database. Also, the project has 2 loggers:
 - pgx tracer, that logs information about all DB-related operations
 - logger implementing Service interface, which can be wrapped aroung the actual struct implementing the Service interface, what allows to debug any implementation of the Service

The project has also Docker setup with air package for live reloads. 

# Usage
The only step needed to run it, is to add the following environment variables, so it would be possible to connect to the database:
- PGUSER
- PGPASSWORD
- PGHOST
- PGDATABASE

If all the variables were added correctly, the project should start successfully by using docker-compose up --build. Also, this project should work with local db as well, if you have specified the mentioned above variables.

## Examples

### GET /movies/{id}
```cURL
curl http://localhost:3000/movies/fcd05f15-216c-4fef-b88f-1a7c90aa43ee
``` 
Response:
```json
{"id":"fcd05f15-216c-4fef-b88f-1a7c90aa43ee","name":"Dune2","release_year":2024,"rating":"8.9","genres":["Action", "Adventure", "Drama"],"director":"Denis Villeneuve"}
```

### GET /movies
```cURL
curl http//localhost:3000/movies
```
Response:
```json
[
  {"id":"fcd05f15-216c-4fef-b88f-1a7c90aa43ee","name":"Dune2","release_year":2024,"rating":"8.9","genres":["Action", "Adventure", "Drama"],"director":"Denis Villeneuve"},
  {"id":"376c60ef-05e4-45af-806c-d4207c9ea43b","name":"Dune","release_year":2021,"rating":"8","genres":["Action","Adventure","Drama"],"director":"Denis Villeneuve"}
]
```

### POST /movies
```cURL
curl -X POST \
-d '{"name": "Dune", "release_year": 2021, "rating": 8.0, "genres": ["Action", "Adventure", "Drama"], "director": "Denis Villeneuve"}' \
http://localhost:3000/movies
```
Response:
```json
{"id":"376c60ef-05e4-45af-806c-d4207c9ea43b"}
```

### PUT /movies/{id}
```cURL
curl -X PUT \
-d '{"rating": 8.2}' \
http://localhost:3000/movies/376c60ef-05e4-45af-806c-d4207c9ea43b
```
No Response

### DELETE /movies/{id}
```cURL
curl -X DELETE \
http://localhost:3000/movies/376c60ef-05e4-45af-806c-d4207c9ea43b
```
No Response
