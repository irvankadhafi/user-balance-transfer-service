-- +migrate Up notransaction
CREATE TYPE transaction_type AS ENUM (
    'DEBIT',
    'CREDIT'
);

CREATE TABLE IF NOT EXISTS "users" (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL,
    password TEXT NOT NULL,
    username TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS "user_balances" (
    id SERIAL PRIMARY KEY,
    user_id INT,
    balance INT NOT NULL,
    balance_achieve INT NOT NULL
);

CREATE TABLE IF NOT EXISTS "user_balance_histories" (
    id SERIAL PRIMARY KEY,
    user_balance_id INT,
    balance_before INT NOT NULL,
    balance_after INT NOT NULL,
    activity TEXT NOT NULL,
    type transaction_type NOT NULL,
    ip_address TEXT NOT NULL,
    location TEXT NOT NULL,
    user_agent TEXT NOT NULL,
    author TEXT
);

CREATE TABLE IF NOT EXISTS "bank_balances" (
    id SERIAL PRIMARY KEY,
    balance INT NOT NULL,
    balance_achieve INT NOT NULL,
    code TEXT NOT NULL,
    enable boolean NOT NULL
);

CREATE TABLE IF NOT EXISTS "bank_balance_histories" (
    id SERIAL PRIMARY KEY,
    bank_balance_id INT,
    balance_before INT NOT NULL,
    balance_after INT NOT NULL,
    activity TEXT NOT NULL,
    type transaction_type NOT NULL,
    ip_address TEXT NOT NULL,
    location TEXT NOT NULL,
    user_agent TEXT NOT NULL,
    author TEXT NOT NULL
);


ALTER TABLE "user_balances" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "user_balance_histories" ADD FOREIGN KEY ("user_balance_id") REFERENCES "user_balances" ("id");

ALTER TABLE "bank_balance_histories" ADD FOREIGN KEY ("bank_balance_id") REFERENCES "bank_balances" ("id");

-- +migrate Down
DROP TABLE IF EXISTS "users";
DROP TABLE IF EXISTS "user_balances";
DROP TABLE IF EXISTS "user_balance_histories";
DROP TABLE IF EXISTS "bank_balances";
DROP TABLE IF EXISTS "bank_balance_histories";