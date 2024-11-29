# movie-festival-be

## instalation

## Intiate Application
Run go mod tidy to download dependency
```
go mod tidy
```

### Confiration
Go to directory /cmd/config/
Update config file config.yaml
```
database:
  main:
    host: "localhost"
    port: "5432"
    user: "postgres"
    password: ""
    dbname: "movie_festival"
    sslmode: "disable"
    timezone: "Asia/Jakarta"
    encoding: "UTF8"
    debug: true

appport: :3000
```

### Migration Database Schema
Go to directory /cmd/migration/
Run Migration
```
go run .\migrate.go --
```

### Runing Application
Go to directory /cmd/
Run Application
```
go run .\main.go
```

### Build Application
Go to directory /cmd/
Build Application
```
go build -v -o movie-festival-be
```