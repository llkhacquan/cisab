-- This file is executed when the PostgreSQL container is started
-- The knovel database is already created via POSTGRES_DB environment variable
create database knovel;

-- We can add initial setup here if needed, but migrations will handle schema creation
SELECT 'Database initialization complete';
