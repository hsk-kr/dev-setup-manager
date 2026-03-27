---
name: refactor
description: Identify code smells and execute safe, incremental refactoring with test verification at every step
argument-hint: <file, module, or description of what to refactor>
allowed-tools: Read, Grep, Glob, Bash, Agent
---

You are performing a safe refactoring. Use `refactoring-surgeon` for the actual changes and `test-engineer` to ensure safety.

## Refactoring: $ARGUMENTS

### Step 1: Analyze (Parallel)
- **codebase-analyzer**: Map the component's structure, dependencies, and callers
- **refactoring-surgeon**: Identify code smells and plan specific refactoring techniques
- **test-engineer**: Check existing test coverage — are there enough tests to refactor safely?

### Step 2: Ensure Safety Net
If test coverage is insufficient:
- Launch `test-engineer` to write characterization tests BEFORE touching anything
- Tests must capture current behavior (even if the behavior is "wrong")
- Verify tests pass

### Step 3: Execute Refactoring
Launch `refactoring-surgeon` with the plan:
- One refactoring operation at a time
- Run tests after each step
- If tests fail, revert and try a smaller step

### Step 4: Cross-Check (Parallel)
After refactoring is complete:
- **test-engineer**: Verify all tests still pass, suggest additional tests for the new structure
- **security-auditor**: Verify no security regressions from the restructuring
- **code-review skill**: Review the final diff for quality

### Output
```
## Refactoring Report

**Scope**: [What was refactored]
**Code Smells Fixed**: [List with before/after]
**Tests**: [X existing passed, Y new tests added]
**Security**: [No regressions / findings]
**Files Changed**: [List]
```
