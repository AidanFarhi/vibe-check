Clean up a finished feature branch by switching to main, pulling latest, and deleting the branch.

Steps — execute them in order, stopping and reporting any error before proceeding:

1. **Capture the current branch name** so it can be deleted after switching:
   ```
   git branch --show-current
   ```

2. **Switch to main**:
   ```
   git checkout main
   ```

3. **Pull latest changes**:
   ```
   git pull
   ```

4. **Delete the feature branch**:
   ```
   git branch -d <branch from step 1>
   ```

5. **Delete the Neon branch** associated with the feature branch:
   - The Neon branch should share the same name as the git branch from step 1.
   - First, list existing Neon branches to confirm the match:
     ```
     neon branches list
     ```
   - If a branch with the exact same name exists, delete it:
     ```
     neon branches delete <branch name>
     ```
   - If no branch matches exactly, or if there is any ambiguity about which Neon branch to delete, **stop and ask the user to confirm** before proceeding.

When done, confirm: the git branch deleted, the Neon branch deleted, and the current HEAD commit on main.
