Stage all uncommitted changes, commit with a concise message, and push to the remote branch.

## Steps

1. Run `git status` and `git diff` in parallel to understand what changed.
2. Draft a commit message:
   - Title: `<area>: <what changed>` — one short phrase, no period
   - Body: 1–3 bullets, one sentence each, describing *why* not just *what*
   - Favor brevity — omit the body entirely if the title is self-explanatory
3. Stage all modified and untracked files relevant to the change (`git add <files>` — avoid `git add -A` if sensitive files might be present).
4. Commit using a HEREDOC so formatting is preserved:
   ```
   git commit -m "$(cat <<'EOF'
   <title>

   <optional body>

   Co-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>
   EOF
   )"
   ```
5. Push to the current remote branch: `git push`.
6. Report the commit hash and pushed branch in one line.
