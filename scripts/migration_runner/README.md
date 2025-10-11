# Migration Runner

A standalone Go application to run database migrations for the DKL Email Service.

## Usage

1. Ensure you have a `.env` file in the project root with database credentials:
   ```
   DB_HOST=your_host
   DB_PORT=5432
   DB_USER=your_user
   DB_PASSWORD=your_password
   DB_NAME=your_database
   DB_SSL_MODE=require
   ```

2. Build the migration runner:
   ```bash
   cd scripts/migration_runner
   go build -o migration_runner .
   ```

3. Run the migration:
   ```bash
   ./migration_runner
   ```

The script will automatically load environment variables from the `.env` file and apply the migration specified in the code.

## Current Migration

This runner applies: `V1_32__migrate_partners_and_radio_recordings.sql`

To apply different migrations, modify the `migrationFile` variable in `apply_migration.go`.