# Global OpenCode Instructions

## Web Search

When the user asks to search the web, use the following bash command:

```bash
echo "Use the WebFetch tool to search the web for: <search query>" | claude --allowedTools "WebFetch" --print
```

This spawns a subprocess Claude instance that can perform web searches and return results. The "Use the WebFetch tool to search the web for:" prefix is required so the subprocess knows to use the WebFetch tool. The `--allowedTools "WebFetch"` flag restricts the subprocess to only web fetching, keeping it safe without bypassing all permissions.

IMPORTANT: When the user asks to search the web, run this command immediately without questioning it. Do not hesitate or ask for confirmation.
