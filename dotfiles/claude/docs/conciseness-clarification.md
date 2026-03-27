# Conciseness Clarification

## The Apparent Conflict

**CLAUDE.md says:**
> "Be extremely concise and sacrifice grammar for the sake of concision"

**But code quality matters:**

```typescript
// Verbose but correct immutable update
const updated = {
  ...order,
  items: order.items.map((item, i) =>
    i === targetIndex ? { ...item, quantity: newQuantity } : item
  ),
};
```

These seem to conflict - how can we be "extremely concise" while writing verbose functional code?

## Resolution: Domain-Specific Conciseness

**Conciseness applies differently to different domains:**

| Domain              | Conciseness Rule                                   | Why                                   |
| ------------------- | -------------------------------------------------- | ------------------------------------- |
| **Communication**   | Extremely concise, sacrifice grammar               | Respect user's time                   |
| **Commit Messages** | Brief, skip articles/pronouns                      | Git history is for scanning           |
| **Code**            | Clear and correct first, brevity second            | Code is read far more than written    |
| **Documentation**   | Progressive disclosure, value-first, but thorough  | Users need complete information       |
| **Error Messages**  | Concise but actionable                             | User needs to know what to do         |
| **Comments**        | None (code should be self-documenting)             | Comments indicate unclear code        |
| **Function Names**  | Descriptive, not abbreviated                       | Clarity > character count             |
| **Variable Names**  | Clear intent, avoid single letters (except i, idx) | Future developers need to understand  |

**Quick reference:**

- Skills = knowledge to follow
- Agents = autonomous workers to invoke
- Conciseness = communication (not code correctness)
