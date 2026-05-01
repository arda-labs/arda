plugins {
    id("java-library")
}

java {
    toolchain {
        languageVersion.set(JavaLanguageVersion.of(25))
    }
}

dependencies {
    implementation("org.springframework.boot:spring-boot-starter-web")
    implementation("io.micrometer:context-propagation:1.1.0")
    implementation("io.micrometer:micrometer-tracing:1.2.0")
    implementation("com.fasterxml.jackson.core:jackson-databind")
    implementation("ch.qos.logback:logback-classic")
    implementation("net.logstash.logback:logstash-logback-encoder:7.4")
    implementation("io.grpc:grpc-stub:1.60.0")
    implementation("net.devh:grpc-server-spring-boot-starter:2.15.0.RELEASE")
}
