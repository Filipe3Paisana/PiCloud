CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL
);

CREATE TABLE files (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    size INT NOT NULL,
    date TIMESTAMP NOT NULL DEFAULT NOW(),
    user_id INT REFERENCES users(id)
);

-- Adicione outras tabelas conforme necessário
