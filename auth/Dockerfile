FROM openjdk:17-jdk-alpine

ENV SERVER_PORT="9999"
ENV POSTGRESQL_URL="localhost:5432/gochat"
ENV POSTGRESQL_USERNAME="gochat"
ENV POSTGRESQL_PASSWORD="gochat"

COPY target/gochat-0.0.1-SNAPSHOT.jar /app/gochat-0.0.1-SNAPSHOT.jar

EXPOSE $SERVER_PORT

ENTRYPOINT ["java","-jar","/app/gochat-0.0.1-SNAPSHOT.jar"]