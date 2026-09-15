# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

### Backend (Go)
```bash
go test ./tests/...       # run all tests
go run ./cmd/server       # run server (dev, skips UI build)
make server-build         # build UI then compile rune-server binary
make cli-build            # compile rune CLI binary
make build-ui             # build React UI only (cd ui && npm run build)
```

### Frontend (`ui/`)
```bash
npm run dev      # Vite dev server with hot-reload (proxies API to :8080)
npm run build    # tsc + vite build (output: ui/dist/, embedded into server binary)
npm run lint     # oxlint
```

## Architecture

### Go Backend

Three-layer design: HTTP handler → service → storage.

| Layer | Path | Role |
|---|---|---|
| Entry / routes | `cmd/server/main.go`, `internal/server/server.go` | Wires middleware, registers all routes |
| Handlers | `internal/api/handler.go` | One `Handler` struct, one method per endpoint |
| Services | `internal/service/` | `SecretService`, `TokenService`, `NamespaceService` |
| Storage | `internal/storage/bolt.go` | BoltDB wrapper; DB lives at `~/.rune/rune.db` |
| Crypto | `internal/crypto/` | AES-256-GCM (`Encrypt`/`Decrypt`) + Argon2id key derivation |
| Seal manager | `internal/seal/seal.go` | Holds AES key in memory; sealing clears it |
| Auth | `internal/middleware/auth.go`, `middleware/namespace.go` | Bearer token → SHA256 → BoltDB lookup; `AuthorizeNamespace` for write ops |

**Secret key format in BoltDB:** `{namespace}/{key}` (latest), `{namespace}/{key}@v{N}` (versioned).

**Auth flow:** Bearer token → SHA256 hash → lookup `TokenRecord` → check `Revoked` → inject into request context. Root token: `id = "root"`, `namespace = "*"`.

**Adding an endpoint:** add method to `Handler`, register in `server.go`, call `middleware.AuthorizeNamespace(r, ns)` before any write.

**No batch endpoint exists.** Bulk operations loop over `SecretService.Put` individually.

### React UI

Vite + React 19 + Tailwind v4 + shadcn/ui (Radix primitives in `components/ui/` — don't edit directly).

**`App.tsx`** owns top-level state: token, sealed flag, current view. The `View` union type is the routing source of truth.

**`lib/api.ts`** — all server calls. Every method uses `hdr(token)` for auth headers and the `check()` helper to throw on non-2xx.

**Adding a view:**
1. Add string literal to `View` type in `App.tsx`
2. Add nav entry in `app-sidebar.tsx` (`rootOnly: true` for admin-only items)
3. Add label in `site-header.tsx` `VIEW_LABELS`
4. Create `ComponentView.tsx`
5. Render in `App.tsx`

**Dialog pattern** (canonical: "Add secret" in `SecretsView.tsx`):
- State resets in `onOpenChange`
- Submit button disabled until inputs are valid
- Async handler with loading state (`putting`, `creating`, etc.)
- `toast.success` / `toast.error('…', { description: String(e) })` for feedback
- Reload data on success

**Namespace access:** `tokenInfo.namespace === "*"` means root → call `api.listNamespaces()` for dropdown. Non-root → fixed to `tokenInfo.namespace`, no dropdown.
