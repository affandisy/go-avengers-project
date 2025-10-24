CREATE TABLE IF NOT EXISTS inventories {
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) NOT NULL,
    stock int NOT NULL CHECK (stock >= 0)
    description VARCHAR(255),
    status VARCHAR(255) NOT NULL CHECK (status IN ('active', 'broken')),
};

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    full_name VARCHAR(50) NOT NULL,
    age INT NOT NULL CHECK (age >= 17),
    occupation VARCHAR(100) NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('admin', 'superadmin'))
);

CREATE TABLE IF NOT EXISTS recipes (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    cook_time INT NOT NULL,
    rating FLOAT NOT NULL
);
