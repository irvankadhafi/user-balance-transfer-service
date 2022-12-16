-- +migrate Up notransaction
ALTER TABLE "users" ADD CONSTRAINT "username_email_unique" unique ("email", "username");

-- +migrate Down
DROP INDEX "username_email_unique" CASCADE;