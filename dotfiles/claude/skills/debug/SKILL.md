---
name: debug
description: Systematic debugging workflow — trace bugs from symptom to root cause using hypothesis-driven investigation
argument-hint: <description of the bug or error message>
allowed-tools: Read, Grep, Glob, Bash, Agent
---

You are debugging an issue. Use the `debugging-detective` agent for investigation and `codebase-analyzer` for understanding code structure.

## Investigation: $ARGUMENTS

### Step 1: Gather Context
Launch in parallel:
- **debugging-detective**: Analyze the error/symptom, form hypotheses, trace the code path
- **codebase-analyzer**: Map the relevant component's architecture, entry points, and data flow

### Step 2: Narrow Down
Based on findings:
- If it's a data issue → trace the data from source to symptom
- If it's a timing issue → look for race conditions, missing awaits
- If it's an environment issue → compare configs across environments
- If it's a regression → use `git bisect` to find the breaking commit

### Step 3: Fix
Once root cause is identified:
1. Launch the appropriate developer agent (typescript-pro, python-pro, etc.) to write the fix
2. Launch `test-engineer` to write a regression test that:
   - Fails WITHOUT the fix
   - Passes WITH the fix
3. Launch `security-auditor` to verify the fix doesn't introduce vulnerabilities

### Step 4: Verify
- Run the full test suite
- Grep for the same pattern elsewhere in the codebase
- If the same bug exists elsewhere, fix ALL instances

### Output
```
## Bug Resolution

**Symptom**: [What was reported]
**Root Cause**: [What was actually wrong — file:line]
**Fix**: [What was changed]
**Regression Test**: [Test that prevents recurrence]
**Related Fixes**: [Other instances of the same bug fixed]
```
