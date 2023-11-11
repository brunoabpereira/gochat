FROM maven:3.8.5-openjdk-17 as build

WORKDIR /src

# install dependencies
COPY pom.xml .
RUN mvn dependency:resolve

# compile
COPY . .
RUN mvn package

###

FROM openjdk:17-jdk-alpine

ENV SERVER_PORT="9999"
ENV POSTGRES_URL="localhost:5432/gochat"
ENV POSTGRES_USERNAME="gochat"
ENV POSTGRES_PASSWORD="gochat"

WORKDIR /app

COPY --from=build /src/target/gochat-0.0.1-SNAPSHOT.jar /app/gochat-0.0.1-SNAPSHOT.jar

EXPOSE $SERVER_PORT

CMD ["java","-jar","/app/gochat-0.0.1-SNAPSHOT.jar"]