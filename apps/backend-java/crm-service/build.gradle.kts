import com.google.protobuf.gradle.*

plugins {
    id("java")
    id("org.springframework.boot")
    id("io.spring.dependency-management")
    id("com.google.protobuf") version "0.9.4"
}

java {
    toolchain {
        languageVersion.set(JavaLanguageVersion.of(25))
    }
}

dependencies {
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

    // gRPC Server (Standard Blocking)
    implementation("net.devh:grpc-server-spring-boot-starter:2.15.0.RELEASE")
    implementation("io.grpc:grpc-netty-shaded:1.60.0")
    implementation("io.grpc:grpc-protobuf:1.60.0")
    implementation("io.grpc:grpc-stub:1.60.0")
    implementation("javax.annotation:javax.annotation-api:1.3.2")

    implementation("io.camunda:spring-zeebe-starter:8.4.0")

    testImplementation("org.springframework.boot:spring-boot-starter-test")
    testImplementation("io.projectreactor:reactor-test")
}

protobuf {
    protoc {
        artifact = "com.google.protobuf:protoc:3.25.1"
    }
    plugins {
        create("grpc") {
            artifact = "io.grpc:protoc-gen-grpc-java:1.60.0"
        }
    }
    generateProtoTasks {
        all().forEach {
            it.plugins {
                create("grpc")
            }
        }
    }
}
