FROM eclipse-temurin:17-jre-alpine

WORKDIR /app

# Install net-tools equivalent using apk
RUN apk update && apk add --no-cache busybox-extras

# Create data directory structure
RUN mkdir -p /data && chmod -R 775 /data

# Copy JAR
COPY pinewoodstore.jar pinewoodstore.jar

# Run with original filename
ENTRYPOINT ["java", "-jar", "pinewoodstore.jar"]
