FROM gradle:8-jdk21 AS builder
WORKDIR /app
COPY . .
# Chạy gradle từ thư mục chứa settings.gradle.kts
WORKDIR /app/apps/backend-java
RUN gradle :crm-service:build -x test

FROM eclipse-temurin:21-jre-alpine
WORKDIR /app
COPY --from=builder /app/apps/backend-java/crm-service/build/libs/*.jar app.jar
ENTRYPOINT ["java", "-jar", "app.jar"]
