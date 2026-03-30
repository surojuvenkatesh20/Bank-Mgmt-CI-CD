CREATE TABLE accounts(
    id bigserial PRIMARY KEY,
    owner varchar NOT NULL,
    balance bigint NOT NULL,
    currency varchar NOT NULL,
    created_at timestamp default now()
);


CREATE TABLE entries(
    id bigserial PRIMARY KEY,
    account_id bigserial REFERENCES accounts (id),
    amount bigint NOT NULL,
    created_at timestamp default now()
);


CREATE TABLE transfers(
    id bigserial PRIMARY KEY,
    from_account_id bigserial references accounts (id),
    -- from_account_id bigserial,
    to_account_id bigserial references accounts (id),
    -- to_account_id bigserial,
    amount bigint NOT NULL,
    created_at timestamp default now()
);


CREATE INDEX ON accounts (id);
CREATE INDEX ON entries (id);
CREATE INDEX ON transfers (id);
CREATE INDEX ON transfers (from_account_id);
CREATE INDEX ON transfers (to_account_id);
CREATE INDEX ON transfers (from_account_id, to_account_id);
