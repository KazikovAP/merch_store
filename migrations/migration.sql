CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    coins INT NOT NULL
);

CREATE TABLE IF NOT EXISTS inventory (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    item_type TEXT NOT NULL,
    quantity INT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_inventory_user_item ON inventory (user_id, item_type);

CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    type TEXT NOT NULL,
    other_user TEXT NOT NULL,
    amount INT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_transactions_user ON transactions (user_id);

CREATE TABLE IF NOT EXISTS merchandise (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    price INT NOT NULL
);
