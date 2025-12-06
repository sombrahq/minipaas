# Repository Guidelines

## Project Structure & Module Organization
- `minipaas-cli/`: Go CLI (entry: `cmd/minipaas/`), with `_test.go` unit tests.
- `minipaas-role/`: Ansible role (tasks, defaults, templates) for Swarm provisioning.
- `minipaas-sql/`: SQL migrations and functions (queues/streams) for PostgreSQL.
- `docs/`: MkDocs site (Material theme). Navigation in `mkdocs.yml`.
- `example-app/`: Minimal example app and Compose files used by documentation and CLI demos.
- `.github/`: CI workflows, issue templates.
- `mk/`: Makefile fragments (docs, build, lint targets).

---

## Build, Test, and Development Commands
- Build CLI:  
  `cd minipaas-cli/cmd/minipaas && go build -o minipaas`
- Run CLI tests:  
  `cd minipaas-cli && go test ./...`
- Format Go:  
  `gofmt -s -w .` (run inside `minipaas-cli`)
- Build docs:  
  `make build-docs` (installs requirements and runs `mkdocs build`)

---

## Documentation

### Structure
- Docs live in `docs/`; nav in `mkdocs.yml`.
- Index pages (`docs/index.md`, `docs/cli/index.md`, `docs/role/index.md`, `docs/sql/index.md`):
    - **No code blocks**.
    - Must contain a **`## Start Here`** section.
    - Must reference their source folder (`minipaas-cli/`, `minipaas-role/`, `minipaas-sql/`).
    - Keep index pages high-level; details belong in subpages.

### Tone
- Describe **what components do**, not what they do not do.
- Keep explanations short and actionable.
- Base statements strictly on repository behavior—no speculation.
- Maintain consistent terminology (manager, worker, rollout, queue, stream, job, cron).

### Subpages
- Installation/usage pages may include code blocks.
- Use real commands and SQL only.
- Keep examples minimal and relevant.

### Linking
- Use relative links only.
- Ensure links match the navigation defined in `mkdocs.yml`.

### Consistency Rules for Agents & Contributors
- Update docs when behavior changes—CLI, role, or SQL.
- Avoid inventing flags, commands, or variables.
- Prefer smaller pages over large ones.
- Follow the positive, capability-focused voice enforced project-wide.

---

## Testing Guidelines
### CLI
- Tests live next to code as `<file>_test.go`.
- Prefer table-driven tests.
- Run via `go test ./...`.

### Role / SQL
- Add reproducible examples in `docs/` when relevant.
- Add molecule/db tests where feasible in their module directories.

---

## Commit & Pull Request Guidelines
- Use conventional commits with component scopes:
    - `feat(cli): ...`
    - `fix(role): ...`
    - `chore(sql): ...`
    - `docs: ...`
- Keep changes focused; avoid mixing refactors with feature work.

---

## Security & Configuration Tips
- Never commit secrets, TLS certs, or Docker contexts.
- Use `.gitignore` for local env files.
- Store sensitive config in secret managers or Ansible Vault.
- Document credential-related steps in `docs/` without exposing real artifacts.

