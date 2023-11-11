
# gochat
Web chat application.

# table of Contents

- [compile, test and build images](#compile-test-and-build-images)
- [intialize db](#intialize-db)
- [run](#run)
- [erase db data](#erase-db-data)
- [description](#description)
- [tests](#tests)

## compile, test and build images
```
sudo docker compose build auth chatserver webapp
```

## intialize db
Create tables and add mock data
```
sudo docker compose --profile initdb up --abort-on-container-exit && \
sudo docker compose --profile initdb down
```

## run
Execute command and access "localhost:80"
```
sudo docker compose --profile all up
```

## erase db data
```
sudo docker volume rm gochat_db
```

## description

Components:
- web server using Go, Gin, GORM and JWT
- chat service using Go, websockets, GORM and JWT
- chat client using React and websockets
- auth service using Java Spring Boot and JWT
- nginx proxy
- postgresql db

Architecture diagram:

<p style="text-align: center">
  <img src="diagram.png" />
</p>

## tests
Tests run automatically when building each service's docker image. 

The following services currently implement tests:

- auth (unit tests)