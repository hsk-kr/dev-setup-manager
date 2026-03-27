#!/usr/bin/env node
// Context monitor hook — warns Claude when context is running low
// Adapted from pablobfonseca/dotfiles
// Runs as PostToolUse hook, reads context metrics from bridge file

const fs = require('fs');
const path = require('path');

const DEBOUNCE_CALLS = 5;
const WARNING_THRESHOLD = 35;
const CRITICAL_THRESHOLD = 25;

let timeout;

async function main() {
  // Read stdin with timeout
  const chunks = [];
  const stdinTimeout = setTimeout(() => process.exit(0), 10000);

  try {
    for await (const chunk of process.stdin) {
      chunks.push(chunk);
    }
  } catch {
    process.exit(0);
  } finally {
    clearTimeout(stdinTimeout);
  }

  const input = Buffer.concat(chunks).toString().trim();
  if (!input) process.exit(0);

  let data;
  try {
    data = JSON.parse(input);
  } catch {
    process.exit(0);
  }

  const sessionId = data.session_id || process.env.CLAUDE_SESSION_ID || 'unknown';
  const stateFile = path.join('/tmp', `claude-ctx-monitor-${sessionId}.json`);

  // Read or init state
  let state = { callCount: 0, lastSeverity: null };
  try {
    state = JSON.parse(fs.readFileSync(stateFile, 'utf8'));
  } catch {}

  state.callCount = (state.callCount || 0) + 1;

  // Try to read context metrics from bridge file
  const bridgeFile = path.join('/tmp', `claude-ctx-${sessionId}.json`);
  let remaining = 100;
  try {
    const bridge = JSON.parse(fs.readFileSync(bridgeFile, 'utf8'));
    remaining = bridge.remaining_percentage || 100;
  } catch {
    // No bridge file — try to extract from input
    if (data.context_window) {
      remaining = data.context_window.remaining_percentage || 100;
    }
  }

  // Determine severity
  let severity = null;
  if (remaining <= CRITICAL_THRESHOLD) {
    severity = 'CRITICAL';
  } else if (remaining <= WARNING_THRESHOLD) {
    severity = 'WARNING';
  }

  // Debounce: skip unless severity changed or escalated
  const severityEscalated = severity && severity !== state.lastSeverity;
  if (!severityEscalated && state.callCount % DEBOUNCE_CALLS !== 0) {
    fs.writeFileSync(stateFile, JSON.stringify(state));
    process.exit(0);
  }

  state.lastSeverity = severity;
  fs.writeFileSync(stateFile, JSON.stringify(state));

  // Output warning to stderr (visible to Claude)
  if (severity === 'CRITICAL') {
    process.stderr.write(
      `\n⚠️ CRITICAL: Context at ${remaining}% remaining. Stop current work and inform the user. Use /handoff to preserve progress.\n`
    );
  } else if (severity === 'WARNING') {
    process.stderr.write(
      `\n⚠️ WARNING: Context at ${remaining}% remaining. Finalize current task soon.\n`
    );
  }
}

main().catch(() => process.exit(0));
