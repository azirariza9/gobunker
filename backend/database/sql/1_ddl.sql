-- DROP TABLE IF EXISTS samples CASCADE;
DROP TABLE IF EXISTS users CASCADE;



DROP TYPE IF EXISTS user_role CASCADE;

CREATE TYPE user_role AS ENUM ('admin', 'user');



-- CREATE TABLE samples (
--     id SERIAL PRIMARY KEY,
--     string VARCHAR(255) NOT NULL
-- );

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role user_role NOT NULL,
    created_at  TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

