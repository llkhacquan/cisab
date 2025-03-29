-- This file is executed when the PostgreSQL container is started
-- The cisab database is already created via POSTGRES_DB environment variable
create database cisab;

-- We can add initial setup here if needed, but migrations will handle schema creation
SELECT 'Database initialization complete';
