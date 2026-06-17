# agent-import

Convert an agent definition file from any AI coding platform into a switchic-compatible YAML agent definition.

## What it does

AI coding platforms store agent definitions as markdown files with YAML frontmatter in a
named directory (e.g. `.claude/agents/<name>.md`). This skill reads one of those files,
auto-detects the platform, maps the platform's field names and casing to the switchic schema,
and writes a valid `agent.yaml` that switchic can register and use.

## Supported platforms

| Platform | Agent file location |
|----------|---------------------|
| Claude Code | `.claude/agents/<name>.md` |
| GitHub Copilot | `.github/agents/<name>.agent.md` |
| Cursor | `.cursor/agents/<name>.md` *(future)* |
| Windsurf | `.windsurf/agents/<name>.md` *(future)* |
| Generic | any `<name>.md` with YAML frontmatter |

## Setup

```bash
switchic add skill agent-import
switchic switch claude
```

## Usage

```
/agent-import <agent-file-path> [output-path]
```

**Claude Code:**
```
/agent-import .claude/agents/planner.md
```

**With an explicit output path:**
```
/agent-import .claude/agents/planner.md ./internal/assets/bundled/agents/planner.yaml
```

## Output path resolution

1. If you pass an explicit output path, it is used.
2. If you pass only an input path, the output is written to the same directory as the input,
   named `agent.yaml`.
3. If you pass nothing, the output is written to `./agent.yaml` in the current directory.

The skill will warn you before overwriting an existing file.

## Output format

```yaml
name: planner
description: What this agent does and when to delegate to it.
required_skills:               # only when source has a skills list
  - commit-msg
instructions: |
  Full system prompt body.

claude:                        # only when claude-specific fields are present
  tools:
    - Read
    - Write
  model: sonnet
  permission_mode: acceptEdits
  max_turns: 20
```

Only fields present in the source are carried over. Claude's camelCase keys (`permissionMode`,
`maxTurns`, `disallowedTools`, `mcpServers`, `initialPrompt`) are converted to snake_case and
nested under the `claude:` block. Comma-separated tool lists are expanded to YAML arrays.
See `formats.md` for the complete field mapping.
