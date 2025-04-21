FROM eclipse-temurin:17-jre-alpine
WORKDIR /app

# Create data directory structure
RUN mkdir -p /data && chmod -R 775 /data


# Copy JAR (adjust path if needed)
COPY pinewoodstore.jar pinewoodstore.jar

# Run with original filename
ENTRYPOINT ["java", "-jar", "pinewoodstore.jar"]