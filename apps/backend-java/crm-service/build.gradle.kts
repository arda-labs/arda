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

    // Spring Boot Starters
    implementation("org.springframework.boot:spring-boot-starter-webflux")
    implementation("org.springframework.boot:spring-boot-starter-json")
    implementation("org.springframework.boot:spring-boot-starter-data-r2dbc")
    implementation("org.postgresql:r2dbc-postgresql")

    // Camunda / Zeebe (Dùng bản Java client)
    implementation("io.camunda:spring-zeebe-starter:8.4.0")

    // Test
    testImplementation("org.springframework.boot:spring-boot-starter-test")
    testImplementation("io.projectreactor:reactor-test")
}
