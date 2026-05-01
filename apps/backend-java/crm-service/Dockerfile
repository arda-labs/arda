FROM eclipse-temurin:21-jdk-alpine AS builder
WORKDIR /app
COPY . .
# We build the specific project. Note: Gradle wrapper is at the root.
RUN ./gradlew :apps:backend-java:crm-service:build -x test

FROM eclipse-temurin:21-jre-alpine
WORKDIR /app
COPY --from=builder /app/apps/backend-java/crm-service/build/libs/*.jar app.jar
ENTRYPOINT ["java", "-jar", "app.jar"]
