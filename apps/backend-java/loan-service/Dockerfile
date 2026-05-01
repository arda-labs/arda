FROM eclipse-temurin:25-jdk AS builder
WORKDIR /app
COPY . .
# Chạy gradle từ thư mục chứa settings.gradle.kts
WORKDIR /app/apps/backend-java
# Install gradle in the jdk image if not using wrapper, but since we are missing wrapper,
# we use the pre-installed gradle if available or install it.
# Better yet, let's use a standard gradle image with jdk 25 if it exists in the future,
# but for now we follow the instruction.
RUN chmod +x gradlew && ./gradlew :loan-service:build -x test

FROM eclipse-temurin:25-jre-alpine
WORKDIR /app
COPY --from=builder /app/apps/backend-java/loan-service/build/libs/*.jar app.jar
ENTRYPOINT ["java", "-jar", "app.jar"]
