-- Import schema.sql first to create the database structure
\i /docker-entrypoint-initdb.d/schema.sql

-- Import seed.sql to populate the database with sample data
\i /docker-entrypoint-initdb.d/seed.sql