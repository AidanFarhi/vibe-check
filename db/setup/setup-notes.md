# Neon Setup Notes

Manual steps required in the Neon console before CI/CD and the app will function.

## 1. Create a Neon Project

Create a new project at console.neon.tech. The default branch (`main`) becomes the production database branch.

Note the **Project ID** — it is required in step 4.

## 2. Run Setup SQL

Connect to the default `neondb` database as the `neondb_owner` role using the Neon SQL editor or `psql`, then run the scripts in order:

```
db/setup/01_create_role.sql   — creates the vibecheck_user role
db/setup/02_create_database.sql — creates vibecheck_db owned by vibecheck_user
```

**Important:** update the password in `01_create_role.sql` before running — the placeholder value `change_me` is not a real password.

After running these scripts, note the connection string for `vibecheck_db` as `vibecheck_user` — this becomes `DATABASE_URL`.

## 3. Generate a Neon API Key

In the Neon console, go to **Account Settings → API Keys** and create a key. This is used by CI to create and delete preview branches per pull request.

## 4. Configure GitHub Secrets and Variables

In the GitHub repo under **Settings → Secrets and variables → Actions**:

| Kind     | Name               | Value                                      |
|----------|--------------------|--------------------------------------------|
| Secret   | `NEON_API_KEY`     | API key from step 3                        |
| Secret   | `DATABASE_URL`     | Production connection string from step 2   |
| Variable | `NEON_PROJECT_ID`  | Project ID from step 1                     |

`DATABASE_URL` is also a Fly.io secret (injected by the deploy workflow automatically), so it only needs to be stored in GitHub.

## Teardown

To fully remove the Neon-side setup, run `db/setup/teardown.sql` as `neondb_owner`. This drops `vibecheck_db` and `vibecheck_user`.

## How Preview Branches Work (automated, no manual steps)

On PR open/reopen/synchronize, CI creates a Neon branch named `preview/pr-{number}-{branch}` forked from `main` and passes its connection string to the Fly.io preview app. On PR close, both the Neon branch and the Fly.io preview app are destroyed.
