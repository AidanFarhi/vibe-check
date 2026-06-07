-- Chart Filter Test Data
-- User: test_filters@test.com  |  Password: TestPass1!
-- Seeds one entry per day for the last 365 days (relative to CURRENT_DATE).
-- Use this to validate the 7D, 30D, 3M, 6M, and 1Y chart filters and the
-- weekly/monthly bucket averaging.
--
-- Metric values follow simple modulo patterns so each metric has its own
-- rhythm; all stay within the 1–10 domain. Score is derived from the formula:
--   (happiness + energy + sleep + (11 - depression) + (11 - pain)) / 5.0


WITH u AS (
    INSERT INTO "user" (email, password_hash)
    VALUES ('test_filters@test.com', '$2a$10$dwBpOxzL42kMfLBWBld9FumxO2AMbTahM9n1Ggqz1msJWyN2TUdzC')
    RETURNING id
),
metrics AS (
    SELECT
        u.id AS user_id,
        (CURRENT_DATE - gs.n)::date                            AS date,
        (3 + ABS((gs.n % 14) - 7) / 2)::smallint               AS depression,  -- 3..6
        (3 + (gs.n % 7))::smallint                             AS happiness,   -- 3..9
        (2 + ((gs.n + 2) % 7))::smallint                       AS pain,        -- 2..8
        (4 + ((gs.n + 4) % 6))::smallint                       AS energy,      -- 4..9
        (5 + (gs.n % 5))::smallint                             AS sleep        -- 5..9
    FROM u, generate_series(0, 364) AS gs(n)
)
INSERT INTO entry (user_id, date, depression, happiness, pain, energy, sleep, score)
SELECT
    user_id, date, depression, happiness, pain, energy, sleep,
    ROUND(
        (happiness::numeric + energy + sleep + (11 - depression) + (11 - pain)) / 5.0,
        2
    )
FROM metrics;
