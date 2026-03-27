---
name: dependency-check
description: Audit project dependencies for vulnerabilities, outdated packages, unused imports, and bundle size impact
allowed-tools: Read, Bash, Glob, Grep, Agent
---

You are performing a dependency audit. Use the `dependency-manager` agent for the core analysis.

## Process

1. **Detect package manager**: Look for `package.json`, `requirements.txt`, `pyproject.toml`, `go.mod`, `Cargo.toml`, `Gemfile`.

2. **Launch dependency-manager agent** with instructions to:
   - Run vulnerability scan (npm audit, pip-audit, govulncheck, etc.)
   - List outdated packages (major/minor/patch separately)
   - Find unused dependencies
   - Check bundle size impact for frontend deps

3. **Present results** in priority order:
   - Critical vulnerabilities (fix NOW)
   - High vulnerabilities (fix this week)
   - Major outdated (plan upgrade)
   - Unused (remove)
   - Minor/patch outdated (batch upgrade)

4. **Offer to fix**:
   - "Want me to fix the critical vulnerabilities now?"
   - "Want me to remove unused dependencies?"
   - "Want me to upgrade patch/minor versions?"

5. **For each fix**: Run tests after the change to verify nothing broke.
