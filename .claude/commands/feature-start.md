Set up a new feature branch and matching Neon DB branch.

Arguments: $ARGUMENTS (the branch name, e.g. `feature/my-feature`)

Steps — execute them in order, stopping and reporting any error before proceeding:

1. **Create and checkout the git branch**
   ```
   git checkout -b $ARGUMENTS
   ```

2. **Create the Neon branch** branched from `production` with the same name:
   ```
   neon branches create --name $ARGUMENTS --parent production
   ```

3. **Get the connection string** for the new branch:
   ```
   neon connection-string $ARGUMENTS --database-name vibecheck_db --role-name vibecheck_user
   ```

4. **Update `.env`** — replace the existing `DATABASE_URL` value with the connection string returned in step 3. The line must stay in the format:
   ```
   DATABASE_URL="<connection string>"
   ```

When done, confirm: git branch name, Neon branch name, and the new endpoint host from the connection string.
