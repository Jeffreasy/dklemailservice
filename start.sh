#!/bin/sh

# Print environment variables for debugging
echo "Starting container with environment variables:"
echo "DB_HOST: $DB_HOST"
echo "DB_PORT: $DB_PORT"
echo "DB_USER: $DB_USER"
echo "DB_NAME: $DB_NAME"
echo "DB_SSL_MODE: $DB_SSL_MODE"

# Run the Docker container with environment variables
docker run -p 8080:8080 \
  -e ADMIN_EMAIL="$ADMIN_EMAIL" \
  -e SMTP_HOST="$SMTP_HOST" \
  -e SMTP_PORT="$SMTP_PORT" \
  -e SMTP_USER="$SMTP_USER" \
  -e SMTP_PASSWORD="$SMTP_PASSWORD" \
  -e SMTP_FROM="$SMTP_FROM" \
  -e ALLOWED_ORIGINS="$ALLOWED_ORIGINS" \
  -e APP_ENV="$APP_ENV" \
  -e REGISTRATION_SMTP_HOST="$REGISTRATION_SMTP_HOST" \
  -e REGISTRATION_SMTP_PORT="$REGISTRATION_SMTP_PORT" \
  -e REGISTRATION_SMTP_USER="$REGISTRATION_SMTP_USER" \
  -e REGISTRATION_SMTP_PASSWORD="$REGISTRATION_SMTP_PASSWORD" \
  -e REGISTRATION_SMTP_FROM="$REGISTRATION_SMTP_FROM" \
  -e REGISTRATION_EMAIL="$REGISTRATION_EMAIL" \
  -e DB_HOST="$DB_HOST" \
  -e DB_PORT="$DB_PORT" \
  -e DB_USER="$DB_USER" \
  -e DB_PASSWORD="$DB_PASSWORD" \
  -e DB_NAME="$DB_NAME" \
  -e DB_SSL_MODE="$DB_SSL_MODE" \
  dklautomatie-backend 