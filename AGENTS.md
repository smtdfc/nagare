# Commit Rules for Agents

You are a strict, no-nonsense Git commit message generator and repository guardian.

These rules are mandatory. A commit that violates them is invalid and must be rewritten before it is pushed.

## 1. Required format

Every commit subject MUST match:

```
type(scope): subject
```

- `type` — one of the allowed types in section 2.
- `scope` — the owning domain path, written exactly as it appears in the repository.
- `subject` — imperative, concise, no trailing period, 72 characters or fewer.

The scope is mandatory. A commit without a scope is invalid.

## 2. Allowed types

| Type | Use for |
|---|---|
| `feat` | A new user-visible capability, endpoint, command, or public API. |
| `fix` | A bug fix or incorrect behavior. |
| `refactor` | Internal restructuring with no intended behavior change. |
| `perf` | A measurable performance improvement. |
| `test` | Tests, fixtures, or test-only helpers. |
| `docs` | Pure documentation: READMEs, guides, comments, examples, and `AGENTS.md`. Always use scope `docs`. |
| `chore` | Dependency and package updates, build files, CI, generated code, and release plumbing. |

Do not invent new types. CI and build changes use `chore(repo)`. Formatting-only changes use `chore(<scope>)`.

### 2.1 Documentation changes

All documentation changes use the `docs` domain, regardless of where the file lives:

```
docs(docs): explain gateway setup
docs(docs): update AGENTS.md commit rules
```

Do not use `docs(repo)`, `docs(core/tool)`, or a module scope for documentation. If a code change also updates comments or examples, the code commit owns that change; split out unrelated prose updates into `docs(docs)`.

### 2.2 Package and dependency updates

Dependency, lockfile, and package-manager updates MUST use `chore`. Never use `feat` or `fix`.

Scope = the module that owns the changed manifest:

```
chore(core): bump go dependencies
chore(pkgs/security): bump keyring dependency
chore(client/web): bump react-router
chore(crates/nagare_vector): bump serde
chore(repo): bump root pnpm dependencies
```

## 3. Scope rules

### 3.1 Scope grammar

```
<module>                                  module root only
<module>/<domain>                         normal case
<module>/<domain>/<subdomain>             deepest owning directory
```

Allowed module/domain prefixes (``docs`` is a logical domain, not a directory):

```
core  gateway  client  plugin  pkgs  cli  dtos  crates  installers  scripts  repo  docs
```

The scope MUST mirror the real directory name, including underscores:

- `core/model_provider` — correct
- `core/model-provider` — wrong
- `core/custom_errors` — correct
- `core/errors` — wrong

Never use a glob or wildcard. `core/*`, `gateway/*`, and `*` are not valid scopes.

### 3.2 Bare names are forbidden

A bare feature or technology name is ambiguous and is forbidden as a scope:

```
chat  plugin  tool  agent  task  config  auth  api  server  backend
frontend  ui  web  db  database  tests  misc  common  utils
```

Wrong:

```
fix(chat): reject empty message
feat(plugin): add telegram command
refactor(tool): split approval helper
```

Correct:

```
fix(gateway/chat): reject empty message
feat(plugin/builtins/telegram): add telegram command
refactor(core/tool/builtins): split approval helper
```

The same word may appear after a module prefix because the prefix removes the ambiguity.

### 3.3 Choose the deepest owning directory

Use the deepest directory that owns the change:

| Changed file | Scope |
|---|---|
| `core/tool/builtins/write_file.go` | `core/tool/builtins` |
| `core/tool/manager/bindings.go` | `core/tool/manager` |
| `core/tool/base_tool.go` | `core/tool` |
| `core/go.mod` | `core` |
| `gateway/main.go` | `gateway` |
| `gateway/chat/controller.go` | `gateway/chat` |
| `gateway/common/middlewares/auth.go` | `gateway/common/middlewares` |
| `pkgs/messages/hil.go` | `pkgs/messages` |
| `plugin/builtins/telegram/handler.go` | `plugin/builtins/telegram` |
| `client/web/vite.config.ts` | `client/web` |
| `README.md`, `AGENTS.md`, or any documentation file | `docs` |
| `.github/workflows/release.yml` | `repo` |

Module-only scopes (`core`, `gateway`, `client`, `plugin`, `pkgs`, `cli`, `dtos`, `crates`, `installers`, `scripts`) are reserved for module-root files and genuinely cross-domain changes inside that module.

Documentation is the exception to the deepest-directory rule: every documentation change uses `docs`, even when the file lives inside a module directory.

### 3.4 One commit, one scope

- Do not mix unrelated domains in one commit. Split them.
- If one logical change touches several domains inside the same module, use the module scope and list every affected domain in the body.
- If one logical change touches several modules, split by module. If the change is an atomic contract change, use the source-of-truth domain as the scope and list regenerated consumers in the body.
- `repo` is a sentinel for workspace-level files only: root build files, `go.work`, CI, and release plumbing. It is not a scope for application features or documentation.

### 3.5 Generated and build output

- `client/messages` is generated from `pkgs/messages` by `pnpm gen:ts`. Commit the regenerated output together with the source change and use the source scope, for example `feat(pkgs/messages): add approval request payload`.
- If only the generator output changes, use `chore(client/messages)`.
- Do not commit `target/`, `dist/`, `node_modules/`, or other local build artifacts.

## 4. Scope catalogue

This catalogue is exhaustive for the current repository. When a change introduces a new domain directory, add its scope here in the same change that introduces it. Do not commit a new directory under an unlisted or improvised scope. Adding a new top-level module is a `repo` change.

### 4.1 `core`

| Scope | Use for |
|---|---|
| `core` | Module-root files (`go.mod`, `go.sum`) and cross-domain changes that cannot be split. |
| `core/agent` | Agent runtime, execution loop, tool executor, agent pool, and agent state. |
| `core/chat` | Agent invocation and the chat worker pipeline. |
| `core/config` | General configuration. Sub-scope: `core/config/manager`. |
| `core/context` | Execution and request-scoped context types. |
| `core/credential` | Credential model and storage. Sub-scope: `core/credential/manager`. |
| `core/custom_errors` | Domain and sentinel error types. |
| `core/event_bus` | Internal event bus. |
| `core/hil` | Human-in-the-loop approval manager and pending-request state. |
| `core/logger` | Logging setup and logger wrappers. |
| `core/mappers` | Mapping between persistence, domain, and message types. |
| `core/media` | Media handling. Sub-scope: `core/media/upload`. |
| `core/memory` | Vector memory and platform-specific memory implementations. |
| `core/message` | Message channel types. |
| `core/model_provider` | Provider manager, base adapter, provider adapters, and provider config. Sub-scopes: `core/model_provider/manager`, `core/model_provider/provider`, `core/model_provider/provider/adapters`. |
| `core/persistence` | Database layer. Sub-scopes: `core/persistence/database`, `core/persistence/database/entities`, `core/persistence/database/repositories`. |
| `core/plugin` | Plugin host, handshake, forwarding, scopes, and lifecycle. Sub-scopes: `core/plugin/host`, `core/plugin/manager`, `core/plugin/internals`. |
| `core/prompt` | Prompt builder and system prompt. |
| `core/session` | Session model, state, and manager. Sub-scope: `core/session/manager`. |
| `core/setup` | Dependency wiring and bootstrap hooks. Sub-scope: `core/setup/hooks`. |
| `core/task` | Task model, triggers, repeat rules, scheduler, and cron worker. Sub-scopes: `core/task/manager`, `core/task/workers`. |
| `core/tool` | Tool contracts and types. Sub-scopes: `core/tool/builtins`, `core/tool/manager`, `core/tool/plugin`, `core/tool/registry`. |

### 4.2 `gateway`

| Scope | Use for |
|---|---|
| `gateway` | Module-root files (`main.go`, `dix.config.json`, `go.mod`, `go.sum`). |
| `gateway/app` | HTTP app bootstrap, route registration, and error handling. |
| `gateway/auth` | Authentication controller and routes. |
| `gateway/chat` | Chat REST/WebSocket controllers, service, and worker. |
| `gateway/common` | Shared gateway infrastructure. Sub-scopes: `gateway/common/config`, `gateway/common/custom_errors`, `gateway/common/guards`, `gateway/common/middlewares`, `gateway/common/websocket`. |
| `gateway/credential` | Credential endpoints and service. |
| `gateway/model_provider` | Model provider endpoints and service. |
| `gateway/plugin` | Plugin endpoints and service. |
| `gateway/settings` | Settings endpoints and service. |
| `gateway/utils` | Gateway response and utility helpers. |

### 4.3 `client`

| Scope | Use for |
|---|---|
| `client/desktop` | Wails desktop shell and desktop build configuration. |
| `client/dtos` | TypeScript DTOs shared by client packages. |
| `client/messages` | Generated TypeScript message types. |
| `client/services` | TypeScript API/service layer. |
| `client/ui` | Shared React components and styles. |
| `client/web` | React web application, routes, and web build configuration. |

### 4.4 `plugin`

| Scope | Use for |
|---|---|
| `plugin/builder` | Plugin compiler, packager, metadata, and build tooling. |
| `plugin/client` | Plugin runtime client, handshake, tool bridge, and model provider bridge. |
| `plugin/dtos` | Plugin-facing DTOs. |
| `plugin/metadata` | Plugin metadata contract. |
| `plugin/builtins/telegram` | Telegram built-in plugin. |

### 4.5 `pkgs`

| Scope | Use for |
|---|---|
| `pkgs/event_bus` | Reusable event bus package. |
| `pkgs/helpers` | JSON, cast, slice, string, and path helpers. |
| `pkgs/ipc` | Unix and Windows IPC transport, dialer, and listener. |
| `pkgs/messages` | Wire message types, including HIL approval messages. |
| `pkgs/paths` | Application path resolution. |
| `pkgs/security` | Keyring, RSA, token, and payload security helpers. |
| `pkgs/system` | OS integrations: browser, files, power, processes, and volume. |

### 4.6 `cli`, `dtos`, `crates`, `installers`, `scripts`

| Scope | Use for |
|---|---|
| `cli` | CLI entrypoint and root command. |
| `cli/cmd` | CLI subcommands (`gateway`, `plugin`, `auth`, `doctor`, `web`, `desktop`). |
| `cli/helpers` | Gateway process, signal, token, and web UI helpers. |
| `dtos/rest` | REST request and response DTOs. |
| `dtos/websocket` | WebSocket DTOs. |
| `crates/nagare_plugin` | Rust plugin crate. |
| `crates/nagare_vector` | Rust vector extension crate. |
| `installers/linux` | Linux installer. |
| `installers/window` | Windows installer. The directory is named `window`; keep that spelling in the scope. |
| `scripts` | Repository code generation and build scripts. |

### 4.7 Repository root

| Scope | Use for |
|---|---|
| `repo` | Root build files, `go.work`, root `Cargo.toml`, root `package.json`, `pnpm` workspace files, `.github/`, and release configuration. |

### 4.8 `docs`

| Scope | Use for |
|---|---|
| `docs` | All documentation regardless of file location: root and per-module READMEs, guides, `docs/` pages, `AGENTS.md`, changelog prose, and standalone examples. Always use type `docs`. |

## 5. Subject and body rules

- Write the subject in the imperative mood: `add`, `fix`, `remove`, `reject`, `split` — not `added`, `fixes`, or `removed`.
- Lowercase the first word unless it is a proper noun.
- Do not end the subject with a period.
- Do not use emojis.
- Keep the subject under 72 characters.
- Describe the outcome, not a list of files.
- Use the body to explain why the change is needed and any non-obvious trade-off.
- Use a `BREAKING CHANGE:` footer for incompatible changes.
- Reference issues with a `Refs:` footer when relevant.

## 6. Agent commit workflow

1. Inspect the work before staging:
   ```bash
   git status --short
   git diff --stat
   git diff
   ```
2. Stage explicit paths by name. Never use `git add -A` or `git add .`; they sweep unrelated user work into the commit.
3. Keep one logical intent per commit. If the change contains a fix, a refactor, and tests for a different domain, split it.
4. Run the relevant tests for the affected module before committing. Do not commit failing code unless the user explicitly asks for a work-in-progress commit.
5. Never amend, rebase, force-push, or push unless the user explicitly requests it.
6. Never revert or overwrite uncommitted changes that you did not create.
7. Never commit secrets, credentials, tokens, or local build output.
8. If a contract change requires regenerated DTOs or messages, regenerate them and include them in the same commit under the source scope.

## 7. Examples

Valid:

```
feat(core/tool/builtins): add volume control tool
fix(gateway/chat): reject empty message payloads
refactor(core/model_provider/provider): extract adapter retry helper
perf(core/memory): cache vector index lookups
test(core/tool/builtins): cover process kill approval paths
docs(docs): document local gateway setup
docs(docs): update AGENTS.md commit rules
chore(pkgs/security): bump keyring dependency
chore(client/web): bump react-router
feat(plugin/builtins/telegram): deduplicate updates by message id
feat(client/web): add approval prompt
feat(crates/nagare_vector): expose batch search
chore(repo): add release workflow
```

Invalid:

```
fix(chat): reject empty message
feat(plugin): add telegram command
refactor(tool): split approval helper
chore(core/*): update modules
feat(core): add volume control tool
fix(tests): repair tool tests
docs(repo): document local gateway setup
fix(client/web): bump react-router
```

The invalid examples are ambiguous, too broad, use a glob, use a bare domain, put documentation outside the `docs` domain, or use a non-`chore` type for a package update.
