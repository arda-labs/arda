plugins {
    id("org.springframework.boot") version "3.2.0" apply false
    id("io.spring.dependency-management") version "1.1.4" apply false
    id("org.graalvm.buildtools.native") version "0.10.0" apply false
    kotlin("jvm") version "1.9.22" apply false
    kotlin("plugin.spring") version "1.9.22" apply false
}

allprojects {
    group = "com.arda.labs"
    version = "0.0.1-SNAPSHOT"

    repositories {
        mavenCentral()
    }
}
