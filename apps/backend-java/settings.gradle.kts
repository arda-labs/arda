pluginManagement {
    repositories {
        mavenCentral()
        gradlePluginPortal()
    }
}

plugins {
    id("org.gradle.toolchains.foojay-resolver-convention") version "0.8.0"
}

rootProject.name = "arda-backend-java"

// Apps
include("crm-service")
include("loan-service")
include("hrm-service")

// Libs - Flattened names for better compatibility
include("libs-java-common")
include("libs-java-database")
include("libs-java-grpc-client")
include("libs-java-security")
include("libs-java-messaging")

project(":libs-java-common").projectDir = file("../../libs/java/common")
project(":libs-java-database").projectDir = file("../../libs/java/database")
project(":libs-java-grpc-client").projectDir = file("../../libs/java/grpc-client")
project(":libs-java-security").projectDir = file("../../libs/java/security")
project(":libs-java-messaging").projectDir = file("../../libs/java/messaging")
