-- +migrate Up notransaction
ALTER TABLE "bank_balances" ADD CONSTRAINT "code_unique" unique ("code");

-- +migrate Down
DROP INDEX "bank_balances" CASCADE;