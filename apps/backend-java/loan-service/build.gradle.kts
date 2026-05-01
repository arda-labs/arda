plugins {
    id("java")
    id("org.springframework.boot")
    id("io.spring.dependency-management")
}

java {
    toolchain {
        languageVersion.set(JavaLanguageVersion.of(25))
    }
}

dependencies {
    // Project dependencies
    implementation(project(":libs-java-common"))
    implementation(project(":libs-java-database"))
    implementation(project(":libs-java-security"))
    implementation(project(":libs-java-messaging"))

    // Spring Boot Starters (Imperative & Modern)
    implementation("org.springframework.boot:spring-boot-starter-web")
    implementation("org.springframework.boot:spring-boot-starter-json")
    implementation("com.fasterxml.jackson.datatype:jackson-datatype-jsr310")
    implementation("org.springframework.boot:spring-boot-starter-data-jpa")
    implementation("org.springframework.boot:spring-boot-starter-kafka")
    implementation("org.postgresql:postgresql")

    // Camunda / Zeebe (Dùng bản Java client)
    implementation("io.camunda:spring-zeebe-starter:8.4.0")

    // Test
    testImplementation("org.springframework.boot:spring-boot-starter-test")
}
