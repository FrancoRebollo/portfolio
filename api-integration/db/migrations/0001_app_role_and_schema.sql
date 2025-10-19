-- 0) App role (owner of the schema/objects)
DO $$
BEGIN
  IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'api_integration') THEN
    CREATE ROLE api_integration
      LOGIN PASSWORD 'api_integration'
      NOSUPERUSER NOCREATEDB NOCREATEROLE NOINHERIT;
  ELSE
    ALTER ROLE api_integration WITH PASSWORD 'api_integration';
  END IF;
END$$;

-- CREATE DATABASE api_integration_db;

-- 1) Extensions you need
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- 2) Create a dedicated schema OWNED by your app role
CREATE SCHEMA IF NOT EXISTS api_int AUTHORIZATION api_integration;

-- (Optional) keep public clean so only superuser can create there
REVOKE CREATE ON SCHEMA public FROM PUBLIC;
GRANT USAGE ON SCHEMA public TO PUBLIC;

-- 3) Make your app role default to that schema
ALTER ROLE api_integration IN DATABASE api_integration_db SET search_path = api_int, public;

-- 4) Build everything as the app role so it OWNS the objects
SET ROLE api_integration;

-- Default privileges for NEW objects the *app role* creates later
ALTER DEFAULT PRIVILEGES IN SCHEMA api_int
  GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO api_integration;
ALTER DEFAULT PRIVILEGES IN SCHEMA api_int
  GRANT USAGE, SELECT, UPDATE ON SEQUENCES TO api_integration;

CREATE TABLE api_int.event_example (
  id_event_example  BIGINT PRIMARY KEY,       -- llega desde afuera
  event_type        varchar(60)  NOT NULL,
  event_content     varchar(60)  NOT NULL,
  fecha_last_update timestamp    NOT NULL DEFAULT now(),
  actualizado_por   varchar(30)  NOT NULL DEFAULT current_user,
  CHECK (id_event_example > 0)                -- opcional, si siempre es positivo
);
