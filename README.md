# claude-hoist

Hoist Claude Code project permissions to your user config.

When you approve tool permissions in a Claude Code project, they're stored in `.claude/settings.local.json` inside that project. If you work across multiple projects, you end up re-approving the same tools over and over.

`claude-hoist` reads your project's permission rules and merges them into your user-level `~/.claude/settings.local.json`, so approved tools carry across all projects.

## Install

```
go install github.com/jeffrydegrande/claude-hoist@latest
```

## Usage

Run from a project directory that has `.claude/settings.local.json`:

```bash
# See what's new (project permissions not yet in your user config)
claude-hoist show

# Preview the exact changes as a unified diff
claude-hoist diff

# Merge all new permissions into your user config
claude-hoist add

# Skip the confirmation prompt
claude-hoist add -y

# Step through each permission one by one (y/n/q)
claude-hoist step

# Open settings in $EDITOR
claude-hoist edit project
claude-hoist edit user
```

## How it works

1. Reads `.claude/settings.local.json` from the current directory
2. Reads `~/.claude/settings.local.json` (your user config)
3. Computes which `allow` and `deny` rules are new
4. Merges them into your user config (deduped and sorted)

No rules are ever removed. The merge is additive only.

## License

MIT
