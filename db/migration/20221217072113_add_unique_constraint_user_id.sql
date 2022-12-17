-- +migrate Up notransaction
ALTER TABLE "user_balances" ADD CONSTRAINT "user_id_unique" unique ("user_id");

-- +migrate Down
DROP INDEX "user_balances" CASCADE;