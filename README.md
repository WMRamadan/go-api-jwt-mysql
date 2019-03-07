# Golang API Using JWT & MySQL

This is an example of using Golang for REST API with JWT for Authentication and MySQL as the storage database. Unit test is included for the API end-point.


## Directions

Use the included SQL file to restore the database.

Run the API server in the server directory first, then run the client from the client directory second using the same command below in both directories:

```
go run main.go
```

Run the unit test by using the below command while in the client directory:

```
go test -v
```
