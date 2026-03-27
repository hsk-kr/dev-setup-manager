---
name: test-plan
description: Generate a comprehensive test plan for a feature, module, or PR, then optionally write the tests using the test-engineer agent
allowed-tools: Read, Grep, Glob, Bash, Agent
---

You are creating a test plan. Use the `test-engineer` agent for writing tests and `codebase-analyzer` for understanding the code.

## Process

1. **Understand what to test**:
   - If a PR/branch: analyze the diff to understand what changed
   - If a feature: read the implementation to understand all code paths
   - If a module: map all public APIs and their contracts

2. **Analyze with codebase-analyzer** (parallel):
   - Identify entry points, dependencies, and side effects
   - Map all code paths including error paths
   - Find existing tests to avoid duplication

3. **Generate the test plan**:

```
## Test Plan: [Feature/Module Name]

### Unit Tests
| # | Test Case | Input | Expected Output | Priority |
|---|-----------|-------|-----------------|----------|
| 1 | [description] | [input] | [expected] | HIGH |

### Integration Tests
| # | Test Case | Components | Expected Behavior | Priority |
|---|-----------|------------|-------------------|----------|

### Edge Cases
| # | Scenario | Why It Matters |
|---|----------|---------------|

### Not Testing (and why)
- [Thing] — [reason: too trivial, already covered by X, etc.]
```

4. **Ask user**: "Want me to write these tests now?"

5. **If yes**: Launch `test-engineer` agent(s) in parallel:
   - One agent per test file/module for maximum parallelism
   - Each agent gets the test plan + relevant source code context
   - After tests are written, run them to verify they pass

6. **After tests written**: Launch `security-auditor` to verify tests don't introduce security issues (no hardcoded credentials in test fixtures, etc.)
