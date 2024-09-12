SELECT 'CREATE DATABASE boodb'
    WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'boodb');