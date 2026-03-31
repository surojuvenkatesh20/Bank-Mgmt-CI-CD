ALTER Table IF EXISTS "accounts" DROP CONSTRAINT IF EXISTS "owner_currency_key";

ALTER Table IF EXISTS "accounts" DROP CONSTRAINT IF EXISTS "accounts_owner_fkey";

DROP Table IF EXISTS "users";