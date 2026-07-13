# skill-import

Convert a `SKILL.md` from any AI coding tool into a switchic-compatible YAML skill definition.

## What it does

All major AI coding tools (Claude Code, Cursor, GitHub Copilot, Windsurf) follow the
[Agent Skills](https://agentskills.io) open standard — skills are `SKILL.md` files with YAML
frontmatter stored inside a named subdirectory. This skill reads one of those files,
auto-detects the platform, and writes a valid `skill.yaml` that switchic can register
and use.

## Supported platforms

| Platform | Skill location |
|----------|---------------|
| Claude Code | `.claude/skills/<name>/SKILL.md` |
| Cursor | `.cursor/skills/<name>/SKILL.md` |
| GitHub Copilot | `.github/skills/<name>/SKILL.md` |
| Windsurf | `.windsurf/skills/<name>/SKILL.md` |
| Generic | `.agents/skills/<name>/SKILL.md` or any `SKILL.md` |

## Setup

```bash
switchic add skill skill-import
switchic switch claude
```

## Usage

```
/skill-import <skill-file-path> [output-path]
```

**Claude Code:**
```
/skill-import .claude/skills/my-skill/SKILL.md
```

**Cursor:**
```
/skill-import .cursor/skills/refactor/SKILL.md
```

**GitHub Copilot:**
```
/skill-import .github/skills/deploy/SKILL.md
```

**Windsurf:**
```
/skill-import .windsurf/skills/code-reviewer/SKILL.md
```

**With an explicit output path:**
```
/skill-import .cursor/skills/refactor/SKILL.md ./internal/assets/bundled/skills/refactor/skill.yaml
```

## Output path resolution

1. If you pass an explicit output path, it is used.
2. If you pass only an input path, the output is written to the same directory as the input,
   named `skill.yaml`.
3. If you pass nothing, the output is written to `./skill.yaml` in the current directory.

The skill will warn you before overwriting an existing file.

## Output format

```yaml
name: my-skill
description: What this skill does and when to invoke it.
argument-hint: "[optional args]"   # only when present in source
allowed-tools: "Read Write"        # only when present in source
prompt: |
  # Full skill instruction body
  ...
```

Only fields present in the source are carried over. Platform-specific fields (`paths`,
`when_to_use`, `license`, `metadata`) are handled per the rules in `formats.md`.

## Sibling files and folders

If the source `SKILL.md` has other files or folders next to it (reference docs, templates,
`examples/`, etc.), they are copied alongside the generated `skill.yaml`, preserving their
relative paths and contents unchanged. switchic's skill loader picks up anything sitting next
to a `skill.yaml` automatically, so the imported skill ships with its full reference set on
every platform. VCS/editor cruft (`.git`, `.DS_Store`) is skipped; everything else is copied
regardless of whether the prompt body links to it.

## Multi-skill files

If the source file contains multiple skill blocks, the skill will ask whether to combine them
into one YAML or split into separate files.
