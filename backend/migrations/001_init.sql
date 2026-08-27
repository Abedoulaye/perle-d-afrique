CREATE TABLE users( 
    id SERIAL PRIMARY KEY, 
    email TEXT NOT NULL UNIQUE, 
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE products(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    description TEXT,
    price_cents INT NOT NULL,
    stock INT NOT NULL DEFAULT 0,
    image TEXT
);

CREATE TABLE carts(
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE cart_items(
    id SERIAL PRIMARY KEY,
    cart_id INT REFERENCES carts(id) ON DELETE CASCADE,
    product_id INT REFERENCES products(id) ON DELETE CASCADE,
    quantity INT NOT NULL CHECK (quantity > 0)
);

CREATE TABLE orders(
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE SET NULL,
    status VARCHAR(50) DEFAULT 'pending',
    total_cents INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE order_items(
    id SERIAL PRIMARY KEY,
    order_id INT REFERENCES orders(id) ON DELETE CASCADE,
    product_id INT REFERENCES products(id) ON DELETE SET NULL,
    quantity INT NOT NULL CHECK (quantity > 0),
    price_cents_at_purchase INT NOT NULL
);

-- product_id doesn't get ON DELETE CASCADE even if the quick cleanup is useful to keep order history
-- instead we set it to null when it is deleted, same for when a user is deleted we keep order history by setting null