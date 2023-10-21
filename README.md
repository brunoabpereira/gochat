
### auth service

compile
```
mvn clean package -DskipTests
```
build image
```
sudo docker build --tag=spring-auth:latest .
```
run
```
sudo docker compose up auth
```