---
name: Apply fixes to all matching patterns
description: When fixing a bug, find and fix ALL instances of the same pattern — don't leave siblings broken
type: feedback
---

When fixing something, always check for sibling cases that need the same fix.

**Why:** User catches incomplete fixes quickly and expects thoroughness on first pass.

**How to apply:** After any fix, grep for similar patterns in the same file and related files. Fix all at once.
