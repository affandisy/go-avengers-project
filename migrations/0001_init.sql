-- Inventories Table
CREATE TABLE IF NOT EXISTS inventories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(50) UNIQUE NOT NULL,
    stock INTEGER NOT NULL CHECK (stock >= 0),
    description TEXT,
    status VARCHAR(10) NOT NULL CHECK (status IN ('active', 'broken')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_inventories_code ON inventories(code);
CREATE INDEX IF NOT EXISTS idx_inventories_status ON inventories(status);

-- Users Table (GORM will auto-create, but here's the manual version)
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    full_name VARCHAR(50) NOT NULL,
    age INTEGER NOT NULL CHECK (age >= 17),
    occupation VARCHAR(100) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'admin' CHECK (role IN ('admin', 'superadmin')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Create indexes for users
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

-- Recipes Table (GORM will auto-create, but here's the manual version)
CREATE TABLE IF NOT EXISTS recipes (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    cook_time INTEGER NOT NULL CHECK (cook_time > 0),
    rating DECIMAL(3,2) NOT NULL CHECK (rating >= 0 AND rating <= 5),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Create indexes for recipes
CREATE INDEX IF NOT EXISTS idx_recipes_deleted_at ON recipes(deleted_at);
CREATE INDEX IF NOT EXISTS idx_recipes_rating ON recipes(rating);

-- Insert sample data (optional)
INSERT INTO inventories (name, code, stock, description, status) 
VALUES 
    ('Laptop Dell', 'LPT001', 10, 'Dell Inspiron 15', 'active'),
    ('Mouse Logitech', 'MSE001', 50, 'Logitech Wireless Mouse', 'active'),
    ('Keyboard Mechanical', 'KBD001', 0, 'Broken keyboard', 'broken')
ON CONFLICT (code) DO NOTHING;