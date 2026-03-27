---
name: pr-comments
description: Check and resolve PR review comments from GitHub
---

Check for pending PR review comments and resolve them.

Steps:

1. Find the current PR:
   ```bash
   gh pr view --json number,url,title
   ```

2. Get all review comments:
   ```bash
   gh api repos/{owner}/{repo}/pulls/{number}/comments --jq '.[] | {id: .id, path: .path, line: .line, body: .body, user: .user.login, created: .created_at}'
   ```

3. Get review threads to check resolution status:
   ```bash
   gh api graphql -f query='
     query {
       repository(owner: "{owner}", name: "{repo}") {
         pullRequest(number: {number}) {
           reviewThreads(first: 100) {
             nodes {
               isResolved
               comments(first: 10) {
                 nodes {
                   body
                   author { login }
                   path
                   line
                 }
               }
             }
           }
         }
       }
     }
   '
   ```

4. For each unresolved comment:
   - Read the referenced file and line
   - Understand the reviewer's feedback
   - Make the requested change (or explain why not)
   - Report what was done

5. Summarize all changes made and any comments that need discussion.
