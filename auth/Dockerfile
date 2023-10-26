FROM openjdk:17-jdk-alpine

ENV SERVER_PORT="9999"
ENV POSTGRES_URL="localhost:5432/gochat"
ENV POSTGRES_USERNAME="gochat"
ENV POSTGRES_PASSWORD="gochat"

COPY target/gochat-0.0.1-SNAPSHOT.jar /app/gochat-0.0.1-SNAPSHOT.jar

EXPOSE $SERVER_PORT

ENTRYPOINT ["java","-jar","/app/gochat-0.0.1-SNAPSHOT.jar"]