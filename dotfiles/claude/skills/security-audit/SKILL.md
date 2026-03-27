---
name: security-audit
description: Run a security audit on recent code changes or the entire codebase, checking for OWASP Top 10, secrets exposure, and dependency vulnerabilities
allowed-tools: Read, Grep, Glob, Bash, Agent
---

You are performing a security audit. Use the `security-auditor` agent for deep analysis. Coordinate the audit as follows:

## Audit Workflow

1. **Determine scope**: Ask the user if they want to audit:
   - Recent changes only (`git diff` since last commit/PR)
   - A specific file or directory
   - The entire codebase

2. **Launch parallel security checks** using TeamCreate or parallel Agent calls:

   **Agent 1 — Code Security (security-auditor)**:
   - Scan for injection vulnerabilities (SQL, XSS, command, path traversal)
   - Check authentication and authorization patterns
   - Review error handling for information leakage
   - Check for hardcoded secrets

   **Agent 2 — Dependency Security (dependency-manager)**:
   - Run `npm audit` / `pip-audit` / equivalent
   - Check for known CVEs in dependencies
   - Flag outdated packages with security implications

   **Agent 3 — Configuration Security (devops-engineer)** (if infra files exist):
   - Check Docker/K8s configs for security misconfigurations
   - Verify secrets management patterns
   - Review CI/CD pipeline for secret exposure

3. **Compile results** into a unified report:

```
## Security Audit Report

### Summary
- Critical: X | High: Y | Medium: Z | Low: W

### Critical Findings
[...]

### High Findings
[...]

### Medium Findings
[...]

### Low Findings
[...]

### Dependencies
[...]

### Recommendations (Prioritized)
1. [Fix immediately]
2. [Fix this sprint]
3. [Improve over time]
```

4. **Offer to fix** critical and high findings immediately using appropriate developer agents.
