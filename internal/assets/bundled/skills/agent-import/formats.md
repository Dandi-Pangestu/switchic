# Agent Import — Source Format Reference

Use this file to detect the source format and map its fields to the switchic schema.

---

## Background

AI coding platforms store agent definitions as markdown files with YAML frontmatter inside a
named directory. The directory location and frontmatter field names differ per platform.

| Platform | Agent file location |
|----------|---------------------|
| Claude Code | `.claude/agents/<name>.md` |
| GitHub Copilot | `.github/agents/<name>.agent.md` |
| Cursor | `.cursor/agents/<name>.md` *(future)* |
| Windsurf | `.windsurf/agents/<name>.md` *(future)* |
| Generic | `.agents/<name>.md` or any `<name>.md` |

---

## Detection rules

Check in this order. Stop at the first match.

| # | Condition | Format |
|---|-----------|--------|
| 1 | File is a `.yaml`/`.yml` AND contains a `name:` key AND an `instructions:` key | Already switchic YAML — report as-is and skip conversion |
| 2 | File has YAML frontmatter (`---` … `---`) with a `name:` key | Platform agent file — detect platform from path (see below) |
| 3 | File is a `.md` or plain text with no YAML frontmatter | Generic markdown agent |

**Platform detection from path** (used for format 2 only):

| Path contains | Filename ends with | Platform |
|---------------|--------------------|----------|
| `.claude/agents/` | `.md` | Claude Code |
| `.github/agents/` | `.agent.md` | GitHub Copilot |
| `.cursor/agents/` | `.md` | Cursor |
| `.windsurf/agents/` | `.md` | Windsurf |
| `.agents/` or unknown | any | Generic |

---

## Format field mappings

### 1. Claude Code — `.claude/agents/<name>.md`

Claude uses camelCase frontmatter keys. All fields map into the `claude:` block of the
switchic YAML (except `name`, `description`, `skills`, and the body).

**Source shape:**

```
---
name: my-agent
description: What this agent does.
tools: Read, Write, Bash
disallowedTools: WebFetch
model: sonnet
permissionMode: acceptEdits
maxTurns: 20
skills: commit-msg, other-skill
memory: project
background: true
effort: high
isolation: worktree
color: blue
initialPrompt: |
  Check memory before starting.
mcpServers:
  - github
hooks:
  PreToolUse:
    - matcher: Bash
      hooks:
        - type: command
          command: ./scripts/validate.sh
---

System prompt body here.
```

**Field mapping:**

| switchic field | Claude frontmatter key | Notes |
|----------------|------------------------|-------|
| `name` | `name` | Fallback: filename stem, kebab-cased |
| `description` | `description` | Synthesize from instructions if absent |
| `required_skills` | `skills` | Comma-separated string → YAML array |
| `instructions` | Body after closing `---` | Strip any leaked frontmatter |
| `claude.tools` | `tools` | Comma-separated string → YAML array |
| `claude.disallowed_tools` | `disallowedTools` | Comma-separated string → YAML array |
| `claude.model` | `model` | Carry over verbatim |
| `claude.permission_mode` | `permissionMode` | camelCase → snake_case |
| `claude.max_turns` | `maxTurns` | Integer |
| `claude.effort` | `effort` | Carry over verbatim |
| `claude.isolation` | `isolation` | Carry over verbatim |
| `claude.background` | `background` | Boolean; omit if false |
| `claude.memory` | `memory` | Carry over verbatim |
| `claude.color` | `color` | Carry over verbatim |
| `claude.initial_prompt` | `initialPrompt` | camelCase → snake_case |
| `claude.mcp_servers` | `mcpServers` | camelCase → snake_case; preserve YAML structure |
| `claude.hooks` | `hooks` | Preserve YAML structure verbatim |

**Comma-separated string → YAML array:**

Claude serializes `tools`, `disallowedTools`, and `skills` as comma-separated strings on a
single line. Split on `, ` (comma + space) or `,` (comma only), trim whitespace from each
element, and emit as a YAML sequence:

```
# Claude source
tools: Read, Write, Bash

# switchic output
claude:
  tools:
    - Read
    - Write
    - Bash
```

---

### 2. GitHub Copilot — `.github/agents/<name>.agent.md`

Copilot uses kebab-case frontmatter keys. Platform-specific fields map into the `copilot:` block
of the switchic YAML.

**Source shape:**

```
---
name: my-agent
description: What this agent does.
tools: ['execute', 'read', 'edit', 'search', 'web', 'todo']
target: github-copilot
model: gpt-4o
disable-model-invocation: true
user-invocable: false
mcp-servers:
  my-server:
    type: http
    url: https://api.example.com/mcp
metadata:
  team: platform
  version: "1.0"
---

System prompt body here.
```

**Field mapping:**

| switchic field | Copilot frontmatter key | Notes |
|----------------|-------------------------|-------|
| `name` | `name` | Fallback: filename stem, strip `.agent.md` suffix |
| `description` | `description` | Synthesize from instructions if absent |
| `instructions` | Body after closing `---` | Strip any leaked frontmatter |
| `copilot.tools` | `tools` | Flow sequence `['a', 'b']` or comma-string → YAML array |
| `copilot.target` | `target` | Carry over verbatim; omit if absent |
| `copilot.model` | `model` | Carry over verbatim |
| `copilot.disable-model-invocation` | `disable-model-invocation` | Boolean; omit if false |
| `copilot.user-invocable` | `user-invocable` | Boolean; omit if true (true is the default) |
| `copilot.mcp-servers` | `mcp-servers` | Preserve YAML structure verbatim |
| `copilot.metadata` | `metadata` | Preserve key-value map verbatim |

**Tools format note:** Copilot serializes `tools` as a YAML flow sequence with single-quoted
strings: `['execute', 'read']`. Parse this as a list and emit as a YAML array in the switchic
output:

```
# Copilot source
tools: ['execute', 'read', 'edit']

# switchic output
copilot:
  tools:
    - execute
    - read
    - edit
```

---

### 3. Already switchic YAML

If `name:` and `instructions:` are both present as top-level YAML keys, the file is already in
switchic format. Report this to the user and do not overwrite unless they explicitly ask.

---

### 4. Generic markdown agent

Plain markdown with no YAML frontmatter.

| switchic field | Source |
|----------------|--------|
| `name` | Filename stem, kebab-cased |
| `description` | First non-heading paragraph; fallback: synthesize from content |
| `instructions` | Full markdown body |

No `claude:` block is emitted for generic markdown since no platform-specific fields are known.

---

## Future platforms (not yet implemented)

When support for these platforms is added, extend the detection rules and add a section here.

| Platform | Expected agent location | Notes |
|----------|------------------------|-------|
| Cursor | `.cursor/agents/<name>.md` | Likely uses snake_case; confirm when spec is available |
| Windsurf | `.windsurf/agents/<name>.md` | Field names TBD |

When implementing a new platform, follow this pattern:
1. Add a detection rule in the table above.
2. Add a "Format field mappings" section for that platform.
3. Map platform-specific fields to either the shared block (`name`, `description`,
   `required_skills`, `instructions`) or a new platform block (e.g. `cursor:`).

---

## Post-extraction cleanup rules

Apply to the `instructions` field regardless of source format:

1. Strip all YAML frontmatter blocks (`---` … `---`) that leaked into the body.
2. Strip platform-specific metadata decorators:
   - HTML comments: `<!-- ... -->`
3. Normalize line endings to `\n`.
4. Preserve all other markdown formatting (headings, lists, code blocks, tables).
5. Do not rewrite the instructions — preserve them verbatim after stripping artifacts.
