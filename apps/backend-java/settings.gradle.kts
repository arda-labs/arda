rootProject.name = "arda-backend-java"

// Apps
include("accounting")
include("crm-service")

// Libs
include(":libs:java:common")
include(":libs:java:database")
include(":libs:java:grpc-client")
include(":libs:java:security")
include(":libs:java:messaging")

// Map to physical paths
project(":accounting").projectDir = file("accounting")
project(":crm-service").projectDir = file("crm-service")
project(":libs:java:common").projectDir = file("../../libs/java/common")
project(":libs:java:database").projectDir = file("../../libs/java/database")
project(":libs:java:grpc-client").projectDir = file("../../libs/java/grpc-client")
project(":libs:java:security").projectDir = file("../../libs/java/security")
project(":libs:java:messaging").projectDir = file("../../libs/java/messaging")

