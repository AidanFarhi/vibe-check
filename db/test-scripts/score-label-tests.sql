-- Score Label Test Data
-- All test users share the password: TestPass1!
-- Today = 2026-05-31  |  Yesterday = 2026-05-30
-- Score formula: (happiness + energy + sleep + (11-depression) + (11-pain)) / 5.0
--
-- Run each block independently, then log in as that user to verify the label.
-- Run the CLEANUP block at the end when done.


-- ============================================================
-- TEST 1: Great ↑ up from yesterday
-- Expected label: "Great ↑ up from yesterday"
-- Today score  = 9.0  (Great,    happiness=9 energy=9 sleep=9 depression=2 pain=2)
-- Yesterday score = 6.0  (Moderate, happiness=6 energy=6 sleep=6 depression=5 pain=5)
-- ============================================================
WITH u AS (
    INSERT INTO "user" (email, password_hash)
    VALUES ('test_great_up@test.com', '$2a$10$dwBpOxzL42kMfLBWBld9FumxO2AMbTahM9n1Ggqz1msJWyN2TUdzC')
    RETURNING id
)
INSERT INTO entry (user_id, date, depression, happiness, pain, energy, sleep, score)
SELECT id, '2026-05-31'::date, 2, 9, 2, 9, 9, 9.00 FROM u
UNION ALL
SELECT id, '2026-05-30'::date, 5, 6, 5, 6, 6, 6.00 FROM u;


-- ============================================================
-- TEST 2: Moderate ↓ down from yesterday
-- Expected label: "Moderate ↓ down from yesterday"
-- Today score    = 6.0  (Moderate, happiness=6 energy=6 sleep=6 depression=5 pain=5)
-- Yesterday score = 9.0  (Great,    happiness=9 energy=9 sleep=9 depression=2 pain=2)
-- ============================================================
WITH u AS (
    INSERT INTO "user" (email, password_hash)
    VALUES ('test_moderate_down@test.com', '$2a$10$dwBpOxzL42kMfLBWBld9FumxO2AMbTahM9n1Ggqz1msJWyN2TUdzC')
    RETURNING id
)
INSERT INTO entry (user_id, date, depression, happiness, pain, energy, sleep, score)
SELECT id, '2026-05-31'::date, 5, 6, 5, 6, 6, 6.00 FROM u
UNION ALL
SELECT id, '2026-05-30'::date, 2, 9, 2, 9, 9, 9.00 FROM u;


-- ============================================================
-- TEST 3: Low — same as yesterday
-- Expected label: "Low same as yesterday"
-- Today score    = 3.0  (Low, happiness=3 energy=3 sleep=3 depression=8 pain=8)
-- Yesterday score = 3.0  (Low, same metrics)
-- ============================================================
WITH u AS (
    INSERT INTO "user" (email, password_hash)
    VALUES ('test_low_same@test.com', '$2a$10$dwBpOxzL42kMfLBWBld9FumxO2AMbTahM9n1Ggqz1msJWyN2TUdzC')
    RETURNING id
)
INSERT INTO entry (user_id, date, depression, happiness, pain, energy, sleep, score)
SELECT id, '2026-05-31'::date, 8, 3, 8, 3, 3, 3.00 FROM u
UNION ALL
SELECT id, '2026-05-30'::date, 8, 3, 8, 3, 3, 3.00 FROM u;


-- ============================================================
-- TEST 4: Poor ↑ up from yesterday
-- Expected label: "Poor ↑ up from yesterday"
-- Today score    = 2.40 (Poor, happiness=3 energy=3 sleep=3 depression=10 pain=9)
-- Yesterday score = 1.60 (Poor, happiness=2 energy=2 sleep=2 depression=10 pain=10)
-- ============================================================
WITH u AS (
    INSERT INTO "user" (email, password_hash)
    VALUES ('test_poor_up@test.com', '$2a$10$dwBpOxzL42kMfLBWBld9FumxO2AMbTahM9n1Ggqz1msJWyN2TUdzC')
    RETURNING id
)
INSERT INTO entry (user_id, date, depression, happiness, pain, energy, sleep, score)
SELECT id, '2026-05-31'::date, 10, 3, 9, 3, 3, 2.40 FROM u
UNION ALL
SELECT id, '2026-05-30'::date, 10, 2, 10, 2, 2, 1.60 FROM u;


-- ============================================================
-- TEST 5: No yesterday — tier label only (no directional suffix)
-- Expected label: "Moderate"
-- Today score = 6.0  (Moderate, happiness=6 energy=6 sleep=6 depression=5 pain=5)
-- No yesterday entry
-- ============================================================
WITH u AS (
    INSERT INTO "user" (email, password_hash)
    VALUES ('test_no_yesterday@test.com', '$2a$10$dwBpOxzL42kMfLBWBld9FumxO2AMbTahM9n1Ggqz1msJWyN2TUdzC')
    RETURNING id
)
INSERT INTO entry (user_id, date, depression, happiness, pain, energy, sleep, score)
SELECT id, '2026-05-31'::date, 5, 6, 5, 6, 6, 6.00 FROM u;


-- ============================================================
-- TEST 6: No today entry — label should be empty
-- Expected label: (nothing shown)
-- No entry for today; one old entry exists just to confirm it is not used
-- ============================================================
WITH u AS (
    INSERT INTO "user" (email, password_hash)
    VALUES ('test_no_today@test.com', '$2a$10$dwBpOxzL42kMfLBWBld9FumxO2AMbTahM9n1Ggqz1msJWyN2TUdzC')
    RETURNING id
)
INSERT INTO entry (user_id, date, depression, happiness, pain, energy, sleep, score)
SELECT id, '2026-05-29'::date, 5, 6, 5, 6, 6, 6.00 FROM u;


-- ============================================================
-- CLEANUP — run after testing to remove all test data
-- Entries are deleted automatically via ON DELETE CASCADE
-- ============================================================
DELETE FROM "user"
WHERE email IN (
    'test_great_up@test.com',
    'test_moderate_down@test.com',
    'test_low_same@test.com',
    'test_poor_up@test.com',
    'test_no_yesterday@test.com',
    'test_no_today@test.com'
);
