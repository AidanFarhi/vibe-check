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

When done, confirm: the branch that was deleted and the current HEAD commit on main.
