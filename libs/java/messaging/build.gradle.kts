plugins {
    id("java-library")
}

java {
    toolchain {
        languageVersion.set(JavaLanguageVersion.of(25))
    }
}

dependencies {
    implementation(project(":libs-java-common"))
    implementation("org.springframework.kafka:spring-kafka")
    implementation("com.fasterxml.jackson.core:jackson-databind")
    implementation("io.cloudevents:cloudevents-kafka:2.5.0")
    implementation("io.cloudevents:cloudevents-json-jackson:2.5.0")
}
