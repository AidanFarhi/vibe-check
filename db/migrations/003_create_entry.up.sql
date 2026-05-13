CREATE TABLE entry (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID        NOT NULL REFERENCES "user"(id),
    date       DATE        NOT NULL,
    depression SMALLINT    NOT NULL CHECK (depression  BETWEEN 1 AND 10),
    happiness  SMALLINT    NOT NULL CHECK (happiness   BETWEEN 1 AND 10),
    pain       SMALLINT    NOT NULL CHECK (pain        BETWEEN 1 AND 10),
    energy     SMALLINT    NOT NULL CHECK (energy      BETWEEN 1 AND 10),
    sleep      SMALLINT    NOT NULL CHECK (sleep       BETWEEN 1 AND 10),
    note       TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (user_id, date)
);
