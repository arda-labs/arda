rootProject.name = "arda-backend-java"

// Apps
include("accounting")

// Libs
include("libs-common")
include("libs-database")
include("libs-grpc-client")
include("libs-security")
include("libs-messaging")

// Map to physical paths
project(":accounting").projectDir = file("accounting")
project(":libs-common").projectDir = file("../../libs/java/common")
project(":libs-database").projectDir = file("../../libs/java/database")
project(":libs-grpc-client").projectDir = file("../../libs/java/grpc-client")
project(":libs-security").projectDir = file("../../libs/java/security")
project(":libs-messaging").projectDir = file("../../libs/java/messaging")

