-- +migrate Up notransaction
CREATE TABLE IF NOT EXISTS "sessions" (
    "id" SERIAL PRIMARY KEY,
    "user_id" INT,
    "access_token" TEXT NOT NULL,
    "refresh_token" TEXT NOT NULL,
    "access_token_expired_at" TIMESTAMP NOT NULL,
    "refresh_token_expired_at" TIMESTAMP NOT NULL,
    "user_agent" text NOT NULL,
    "location" text NOT NULL,
    "ip_address" text NOT NULL,
    "updated_at" TIMESTAMP NOT NULL DEFAULT 'now()',
    "created_at" TIMESTAMP NOT NULL DEFAULT 'now()'
);

ALTER TABLE "sessions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "sessions" ADD CONSTRAINT "access_token_unique" unique ("access_token");
ALTER TABLE "sessions" ADD CONSTRAINT "refresh_token_unique" unique ("refresh_token");

-- +migrate Down
DROP TABLE IF EXISTS "sessions";