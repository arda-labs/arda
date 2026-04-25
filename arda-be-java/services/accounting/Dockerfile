# Build stage
FROM gradle:8-jdk21 AS build
WORKDIR /app

# Copy all files for context
COPY . .

# Build the application
RUN gradle :services:accounting:build -x test

# Runtime stage
FROM eclipse-temurin:21-jre-jammy
WORKDIR /app
COPY --from=build /app/services/accounting/build/libs/*.jar app.jar

EXPOSE 8080
ENTRYPOINT ["java", "-jar", "app.jar"]
