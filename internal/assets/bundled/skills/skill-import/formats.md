# Skill Import — Source Format Reference

Use this file to detect the source format and map its fields to the switchic schema.

---

## Background

All major AI coding tools (Claude Code, Cursor, GitHub Copilot, Windsurf) follow the
[Agent Skills](https://agentskills.io) open standard. Skills are stored as `SKILL.md` files
with YAML frontmatter inside a named subdirectory. The directory location varies per platform:

| Platform | Skill directory |
|----------|----------------|
| Claude Code | `.claude/skills/<name>/SKILL.md` |
| Cursor | `.cursor/skills/<name>/SKILL.md` or `.agents/skills/<name>/SKILL.md` |
| GitHub Copilot | `.github/skills/<name>/SKILL.md` or `.agents/skills/<name>/SKILL.md` |
| Windsurf | `.windsurf/skills/<name>/SKILL.md` or `.agents/skills/<name>/SKILL.md` |
| Generic | any `.agents/skills/<name>/SKILL.md` |

---

## Detection rules

Check in this order. Stop at the first match.

| # | Condition | Format |
|---|-----------|--------|
| 1 | File is a `.yaml`/`.yml` AND contains a `name:` key AND a `prompt:` key | Already switchic YAML — report as-is and skip conversion |
| 2 | File has YAML frontmatter (`---` … `---`) with a `name:` key | Agent Skills SKILL.md — detect platform from path (see below) |
| 3 | File is a `.md` or plain text with no YAML frontmatter | Generic markdown skill |

**Platform detection from path** (used for format 2 only):

| Path contains | Platform |
|---------------|----------|
| `.claude/` | Claude Code |
| `.cursor/` | Cursor |
| `.github/` | GitHub Copilot |
| `.windsurf/` | Windsurf |
| `.agents/` or unknown | Generic / Agent Skills |

---

## Format field mappings

### 1. Agent Skills SKILL.md (all platforms)

All platforms share this file shape per the [agentskills.io](https://agentskills.io/specification) spec:

```
---
name: my-skill
description: What this skill does and when to use it.
# --- platform-specific fields below ---
argument-hint: "[args]"          # Claude Code
disable-model-invocation: true   # Claude Code, Cursor
allowed-tools: "shell"           # Claude Code, GitHub Copilot (experimental)
paths: "src/**/*.ts"             # Cursor — glob scope
when_to_use: "..."               # Windsurf — extended activation guidance
license: MIT                     # agentskills.io / GitHub Copilot — optional
compatibility: "Requires git"    # agentskills.io / GitHub Copilot — optional
metadata:                        # agentskills.io / GitHub Copilot — optional
  key: value
---

# Prompt body starts here
...
```

**`name` validation (agentskills.io constraints):**

The `name` field must satisfy these rules — flag a warning if the source violates any of them:
- 1–64 characters
- Lowercase letters (`a-z`), digits (`0-9`), and hyphens (`-`) only
- Must not start or end with a hyphen
- Must not contain consecutive hyphens (`--`)
- Must match the parent directory name

**Core field mapping (applies to all platforms):**

| switchic field | Source |
|----------------|--------|
| `name` | frontmatter `name`; fallback: directory or filename stem |
| `description` | frontmatter `description`; fallback: synthesize from prompt first paragraph |
| `prompt` | everything after the closing `---` delimiter |

**Platform-specific field handling:**

| Field | Platform | Action |
|-------|----------|--------|
| `argument-hint` | Claude Code | → `claude.argument-hint` |
| `disable-model-invocation` | Claude Code / Cursor | → `claude.disable-model-invocation` |
| `allowed-tools` | Claude Code | → `claude.allowed-tools` (verbatim string) |
| `allowed-tools` | GitHub Copilot | → `copilot.allowed-tools` (verbatim string) |
| `license` | GitHub Copilot / agentskills.io | → `copilot.license` |
| `compatibility` | GitHub Copilot / agentskills.io | → `copilot.compatibility` |
| `metadata` | GitHub Copilot / agentskills.io | → `copilot.metadata` (preserve key-value map) |
| `paths` | Cursor | Do not carry over. Preserve as inline comment at top of `prompt`: `# Originally scoped to: <value>` |
| `when_to_use` | Windsurf | If non-empty and distinct from `description`, append to `description` as a second sentence |

**How to determine platform for field routing:**

When the source path is `.github/skills/`, route `allowed-tools`, `license`, `compatibility`,
and `metadata` into `copilot:`. When the source path is `.claude/skills/`, route `allowed-tools`,
`argument-hint`, and `disable-model-invocation` into `claude:`. When the path is ambiguous
(e.g. `.agents/skills/`), treat `license`, `compatibility`, and `metadata` as `copilot:` fields
(per agentskills.io spec), and treat `argument-hint` / `disable-model-invocation` as `claude:` fields.

---

### 2. Already switchic YAML

If `name:` and `prompt:` are both present as top-level YAML keys, the file is already in switchic
format. Report this to the user and do not overwrite unless they explicitly ask.

---

### 3. Generic markdown skill

Plain markdown with no YAML frontmatter.

| switchic field | Source |
|----------------|--------|
| `name` | directory name or filename stem, kebab-cased |
| `description` | first non-heading paragraph; fallback: synthesize from content |
| `prompt` | full markdown body |

---

## Post-extraction cleanup rules

Apply to the `prompt` field regardless of source format:

1. Strip all YAML frontmatter blocks (`---` … `---`) that leaked into the body.
2. Strip platform-specific metadata decorators:
   - HTML comments: `<!-- ... -->`
   - Cursor `@file`, `@folder`, `@codebase` reference lines (keep their semantic meaning as
     plain prose if it matters to the instruction).
3. Normalize line endings to `\n`.
4. Preserve all other markdown formatting (headings, lists, code blocks, tables).
5. Do not rewrite the user's instructions — preserve them verbatim after stripping artifacts.
