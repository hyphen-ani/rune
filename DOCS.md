# Rune — Complete Guide

Rune is a **local-first secrets vault** you run yourself. Your secrets never leave your machine unless you choose to move them. Everything is encrypted at rest with AES-256-GCM, and the vault can be sealed so that even if someone has access to the database file, they cannot read anything without the passphrase.

---

## Table of Contents

1. [How Rune Works](#how-rune-works)
2. [Installation](#installation)
3. [Starting for the First Time](#starting-for-the-first-time)
4. [Core Concepts](#core-concepts)
5. [CLI Commands](#cli-commands)
   - [Vault Control](#vault-control)
   - [Authentication](#authentication)
   - [Secrets](#secrets)
   - [Namespaces](#namespaces)
   - [Tokens](#tokens)
   - [System Tokens](#system-tokens)
6. [Web UI](#web-ui)
7. [Token Permissions](#token-permissions)
8. [Use Cases](#use-cases)
9. [Java Spring Boot SDK](#java-spring-boot-sdk)
10. [HTTP API Quick Reference](#http-api-quick-reference)
11. [Security Model](#security-model)
12. [Tips & Best Practices](#tips--best-practices)

---

## How Rune Works

```
┌─────────────────────────────────────────────────────────┐
│                    rune-server                          │
│                                                         │
│   ┌──────────┐   ┌──────────┐   ┌──────────────────┐   │
│   │  HTTP    │   │  Auth    │   │  Seal Manager    │   │
│   │  API     │──▶│Middleware│──▶│  (AES key in     │   │
│   │ :2710    │   │ (Bearer) │   │   memory only)   │   │
│   └──────────┘   └──────────┘   └──────────────────┘   │
│         │                               │               │
│         ▼                               ▼               │
│   ┌──────────────────────────────────────────────────┐  │
│   │              BoltDB  (~/.rune/rune.db)           │  │
│   │       AES-256-GCM encrypted secrets at rest      │  │
│   └──────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────┘
        ▲                               ▲
        │                               │
  rune (CLI)                      Browser / SDK
```

Rune has two binaries:

| Binary | Role |
|---|---|
| `rune-server` | The vault daemon. Runs the HTTP API on port **2710**, manages encryption, stores data in BoltDB |
| `rune` | The CLI client. Talks to the server over HTTP using a saved token |

**The seal/unseal cycle** is the key idea:

- On first start, you create a **passphrase**. This derives the AES-256 encryption key via Argon2id and the key is held **in memory only**.
- The vault starts **sealed** — the key is not loaded, no secrets can be read.
- You **unseal** by entering your passphrase. The key is derived again and loaded into memory.
- You can **seal** at any time. The key is wiped from memory. The database stays on disk but is unreadable without unsealing again.
- If the server restarts, you must unseal again.

---

## Installation

### Prerequisites

- Go 1.22+
- Node.js 18+ and npm (for building the web UI)

### Build and Install from Source

```bash
# Clone the repository
git clone https://github.com/your-org/rune.git
cd rune

# Build and install both binaries (requires sudo to move to /usr/local/bin)
make server-install   # builds UI + server binary → /usr/local/bin/rune-server
make cli-install      # builds CLI binary → /usr/local/bin/rune

# Or build without installing
make server-build     # produces ./rune-server
make cli-build        # produces ./rune
```

### What Gets Created

When `rune-server` starts for the first time it creates:

```
~/.rune/
├── rune.db        ← BoltDB database (all secrets, tokens, namespaces)
└── config.json    ← CLI token saved by `rune login`
```

---

## Starting for the First Time

### 1. Start the server

```bash
rune-server
```

On the **very first run** (no database exists), Rune prompts you to create a passphrase:

```
Set A Passphrase (Required For Sealing and Unsealing):
```

Type your passphrase and press Enter. The characters are **hidden** — nothing is echoed to the terminal.

Rune then prints your **root token** — copy and save it somewhere safe:

```
====================================
Root Token (SAVE THIS): rnk_a1b2c3d4e5f6...
====================================
Store this token securely. It will not be shown again.
```

The server banner then appears:

```
  rune  local-first secrets vault

  ➜  api   http://localhost:2710
  ➜  ui    http://localhost:2710/ui/
  ➜  db    /Users/you/.rune/rune.db

  ⬡  vault sealed  —  run rune unseal or open the web UI
```

### 2. Unseal the vault

The vault starts **sealed** after every server start. Unseal it with:

```bash
rune unseal
```

```
Enter passphrase:
✓ Vault Unsealed Successfully
```

### 3. Log in with the root token

Save your root token to the CLI so you don't have to pass it every time:

```bash
rune login rnk_a1b2c3d4e5f6...
✓ Token saved successfully
```

You're ready to go.

---

## Core Concepts

### Sealed vs Unsealed

| State | What works |
|---|---|
| **Sealed** | `rune status`, `rune unseal`, Web UI unseal page |
| **Unsealed** | Everything — read, write, delete, rotate secrets |

Think of "sealed" as the vault being locked. Even with physical access to `rune.db`, an attacker sees only ciphertext — the key exists only in memory and only while unsealed.

### Secrets

A secret is a **key → value** pair stored encrypted in a namespace. Keys are strings (e.g. `DATABASE_URL`, `api/stripe/key`). Values are arbitrary strings. Every write creates a new **version** — old values are kept and can be retrieved by version number.

```
namespace: production
  DATABASE_URL   = "postgres://..." (v3, current)
                   "postgres://..." (v2)
                   "postgres://..." (v1)
  STRIPE_KEY     = "sk_live_..."   (v1, current)
```

### Namespaces

Namespaces are logical buckets for secrets. They let you:

- Separate secrets by **environment** (`production`, `staging`, `dev`)
- Separate secrets by **service** (`backend`, `frontend`, `worker`)
- Grant a token access to one namespace without seeing others

The `default` namespace always exists and cannot be deleted.

### Tokens

Every API call requires a **Bearer token**. There are three token types:

| Type | Prefix | Access |
|---|---|---|
| **Root** | `rnk_` | Full access to everything |
| **Regular** | `tkn_` | Fine-grained: per-namespace, per-permission |
| **System** | `stk_` | Read-only, works even while vault is sealed (for SDKs) |

---

## CLI Commands

### Vault Control

#### `rune status`
Check whether the vault is sealed or unsealed. No token required.

```bash
rune status
# 🔒 Vault is Sealed
# 🔓 Vault is Unsealed
```

#### `rune unseal`
Unseal the vault. Prompts for your passphrase (hidden input).

```bash
rune unseal
# Enter passphrase:
# ✓ Vault Unsealed Successfully
```

#### `rune seal`
Seal the vault. Wipes the encryption key from memory immediately.

```bash
rune seal
# ✓ Vault Sealed Successfully
```

> **Tip:** Seal the vault when you leave your desk in shared environments. Nobody can read secrets until you unseal again.

---

### Authentication

#### `rune login [token]`
Save a token to `~/.rune/config.json`. All subsequent CLI commands use this token automatically.

```bash
rune login rnk_a1b2c3d4e5f6...
# ✓ Token saved successfully
```

To switch accounts, just run `login` again with a different token.

---

### Secrets

All secret commands accept `-n` / `--namespace` to target a specific namespace. The default is `default`.

#### `rune put [key] [value]`

Store a secret. If the key already exists, creates a new version.

```bash
# Store in the default namespace
rune put DATABASE_URL "postgres://localhost/mydb"

# Store in a specific namespace
rune put DATABASE_URL "postgres://prod-host/mydb" -n production

# Multi-word values — always quote them
rune put JWT_SECRET "my secret phrase with spaces" -n backend

# Warning: avoid unquoted $ — shell will expand it
rune put API_KEY "$SECRET"        # ❌ shell expands $SECRET
rune put API_KEY '$SECRET'        # ✓ literal string
```

#### `rune get [key]`

Retrieve and decrypt a secret value.

```bash
rune get DATABASE_URL
# [SECRET] postgres://localhost/mydb

rune get DATABASE_URL -n production
# [SECRET] postgres://prod-host/mydb

# Get a specific historical version
rune get DATABASE_URL --version 2
# [SECRET] postgres://old-host/mydb
```

#### `rune list`

List all secret keys in a namespace (values are never shown).

```bash
rune list
# Stored Secrets:
#   [SECRET] DATABASE_URL
#   [SECRET] REDIS_URL
#   [SECRET] JWT_SECRET

rune list -n production
```

#### `rune delete [key]`

Delete a secret and all its versions. Asks for confirmation.

```bash
rune delete OLD_API_KEY -n staging
# Are you sure you want to delete secret 'OLD_API_KEY' (namespace: staging)? (y/N): y
# ✓ Secret deleted
```

#### `rune rotate [key]`

Generate a cryptographically random new value for an existing secret. Useful for keys, passwords, and tokens that you want to refresh without choosing a value manually.

```bash
rune rotate API_KEY -n production
#   SECRET ROTATED
#   ─────────────────────────────────────────────
#   Namespace:   production
#   Key:         API_KEY
#   New Value:   f8a2c91e4d...
#   ─────────────────────────────────────────────
# ✓ Secret rotated successfully
```

The old value is preserved as a previous version and is still accessible via `--version`.

#### `rune history [key]`

Show the full version history of a secret with timestamps.

```bash
rune history DATABASE_URL -n production
#   SECRET HISTORY
#   ─────────────────────────────────────────────────────────
#   VERSION    CREATED                    ROTATED
#   ─────────────────────────────────────────────────────────
#   1          2025-01-10T09:00:00Z       -
#   2          2025-01-10T09:00:00Z       2025-06-01T14:30:00Z
#   3          2025-01-10T09:00:00Z       2025-10-15T09:00:00Z
#   ─────────────────────────────────────────────────────────
```

To retrieve a specific version's value:

```bash
rune get DATABASE_URL --version 1 -n production
```

---

### Namespaces

#### `rune namespace create [name]`

Create a new namespace. Requires root access.

```bash
rune namespace create production
rune namespace create staging
rune namespace create payments-service
```

#### `rune namespace list`

List all namespaces.

```bash
rune namespace list
# Namespaces:
#   [NAMESPACES] default
#   [NAMESPACES] production
#   [NAMESPACES] staging
```

#### `rune namespace delete [name]`

Delete a namespace and all secrets it contains. Asks for confirmation. The `default` namespace cannot be deleted.

```bash
rune namespace delete old-staging
# Are you sure you want to delete namespace 'old-staging'? (y/N): y
# ✓ Namespace deleted
```

---

### Tokens

Token management requires the **root token**.

#### `rune token create [name]`

Create a regular access token with fine-grained namespace permissions.

The `-n` flag takes the format `namespace:permissions` where permissions is a comma-separated list of `read`, `write`, `delete`, `rotate`. Omitting permissions after `:` grants all four. The flag is repeatable for multi-namespace access.

```bash
# Read-only access to production
rune token create read-only-prod -n production:read

# Read + write to production (no delete, no rotate)
rune token create ci-deploy -n production:read,write

# Full access to production
rune token create admin-prod -n production:read,write,delete,rotate

# Access to multiple namespaces with different permissions
rune token create backend-svc \
  -n production:read \
  -n staging:read,write

# Full access to a namespace (omit permissions after colon)
rune token create dev-token -n staging
```

Output:
```
✓ Token created successfully

  ID:          tkn_1728000000000
  Name:        ci-deploy
  Created At:  2025-10-15T09:00:00Z

  production           read, write

  Token: tkn_a1b2c3d4e5f6...
```

> **Important:** Copy the token now. It is shown only once.

#### `rune token list`

List all tokens with their status.

```bash
rune token list
# Tokens:
#   ID                           NAME                 CREATED AT                TYPE       STATUS
#   ------------------------------------------------------------------------------------------------
#   tkn_1728000000000            ci-deploy            2025-10-15T09:00:00Z      regular    active
#   stk_1728000000001            java-sdk             2025-10-15T10:00:00Z      system     active
#   tkn_1728000000002            old-token            2025-09-01T08:00:00Z      regular    revoked
```

#### `rune token revoke [token-id]`

Revoke a token immediately. Revoked tokens are rejected by the server.

```bash
rune token revoke tkn_1728000000002
# ✓ Token revoked
```

---

### System Tokens

System tokens (`stk_` prefix) are **read-only** tokens that work **even while the vault is sealed**. They are designed for SDKs, background services, and CI runners that need secrets at startup before a human has unsealed the vault.

System tokens achieve sealed-access through a separate encryption layer: each secret is re-encrypted with a key derived from the system token itself, so the token alone is sufficient to decrypt.

#### `rune system-token create [name]`

```bash
rune system-token create java-sdk -n production

# ✓ System token created successfully

#   ID:          stk_1728000000001
#   Name:        java-sdk
#   Namespace:   production
#   Read-only:   true
#   Created At:  2025-10-15T10:00:00Z
#   Token:       stk_a1b2c3d4e5...

#   Use this token in your SDK config. It cannot write or delete secrets.
```

System tokens are scoped to a **single namespace**. To access multiple namespaces, create one system token per namespace.

---

## Web UI

Open `http://localhost:2710/ui/` in your browser while the server is running.

### Unseal from the browser

If the vault is sealed, you'll see the unseal screen. Enter your passphrase and click Unseal — the vault unlocks and you land on the Secrets view.

### Secrets view

The Secrets view has two modes toggled at the top:

**Table mode** — lists all keys in the selected namespace with action buttons:
- 👁 **View** — reveal the decrypted value
- ↻ **Rotate** — generate a new random value
- 🕐 **History** — browse all versions with timestamps and copy-on-demand
- 🗑 **Delete** — permanently remove the key and all versions

**JSON mode** — loads all secrets in the namespace as an editable JSON document in a syntax-highlighted code editor. Edit the JSON directly and click **Save**:
- Added keys → created as new secrets
- Changed values → updated (new version created)
- Removed keys → confirmation dialog before permanent deletion

Action buttons are automatically **disabled** based on what your token can do. A read-only token sees the JSON editor in read-only mode with a "read-only" badge.

### Import

Click **Import** to bulk-load secrets from a file or paste:
- **`.env` files** — standard `KEY=VALUE` format, comments (`#`) and quoted values supported
- **`.json` files** — flat `{"KEY": "value"}` object
- **Paste** — drop raw `.env` or JSON text directly into the paste box

---

## Token Permissions

| Permission | Allowed operations |
|---|---|
| `read` | `get`, `list`, `history`, `version` |
| `write` | `put` (create and update) |
| `delete` | `delete` |
| `rotate` | `rotate` |

Permissions are **per-namespace**. A single token can have different permissions across namespaces:

```bash
# This token can read from production but fully manage staging
rune token create mixed-env \
  -n production:read \
  -n staging:read,write,delete,rotate
```

### Permission matrix by token type

| Token type | Read | Write | Delete | Rotate | Sealed access |
|---|---|---|---|---|---|
| Root | ✓ | ✓ | ✓ | ✓ | ✗ |
| Regular | per-token | per-token | per-token | per-token | ✗ |
| System | ✓ | ✗ | ✗ | ✗ | ✓ |

---

## Use Cases

### Local development

Store all your dev credentials in Rune instead of `.env` files that accidentally end up in git.

```bash
# One-time setup
rune-server &
rune unseal
rune login rnk_...

# Store your dev secrets
rune put DATABASE_URL "postgres://localhost/myapp_dev"
rune put REDIS_URL "redis://localhost:6379"
rune put STRIPE_KEY "sk_test_..."
rune put SENDGRID_API_KEY "SG...."

# Use in scripts
export DATABASE_URL=$(rune get DATABASE_URL)
export REDIS_URL=$(rune get REDIS_URL)
```

### Environment separation

Keep production and staging secrets completely isolated:

```bash
rune namespace create production
rune namespace create staging

rune put DATABASE_URL "postgres://prod-db/app" -n production
rune put DATABASE_URL "postgres://staging-db/app" -n staging

# Create a staging-only token for your staging CI
rune token create staging-ci -n staging:read,write
```

### CI/CD pipeline

Create a read-only token per environment. Inject it as a CI secret. Your pipeline fetches what it needs at deploy time — no plaintext secrets in your repo or CI config.

```bash
# Create a deploy token (read + write so it can update secrets during deploy)
rune token create github-actions-prod -n production:read,write
# → tkn_abc123...

# In your GitHub Actions workflow:
# secrets.RUNE_TOKEN = tkn_abc123...

# In the workflow step:
rune login ${{ secrets.RUNE_TOKEN }}
export DATABASE_URL=$(rune get DATABASE_URL -n production)
export API_KEY=$(rune get STRIPE_KEY -n production)
```

### Multi-service architecture

Each service gets its own namespace and token. A service can only see its own secrets.

```bash
rune namespace create auth-service
rune namespace create payment-service
rune namespace create notification-service

# Auth service gets its JWT secret and DB
rune put JWT_SECRET "..." -n auth-service
rune put DATABASE_URL "..." -n auth-service
rune token create auth-svc-token -n auth-service:read

# Payment service gets its Stripe key and DB
rune put STRIPE_KEY "..." -n payment-service
rune put DATABASE_URL "..." -n payment-service
rune token create payment-svc-token -n payment-service:read
```

### Secret rotation

Rotate a secret and ensure all services pick up the new value at the next restart:

```bash
# Generate a fresh random value
rune rotate STRIPE_WEBHOOK_SECRET -n production

# Need the old value? Still accessible
rune history STRIPE_WEBHOOK_SECRET -n production
rune get STRIPE_WEBHOOK_SECRET --version 2 -n production
```

### Read-only audit access

Give an auditor access to see what keys exist and what values are, without being able to change anything:

```bash
rune token create auditor -n production:read
# Share just the token — they can list and view, but not write, delete, or rotate
```

### Sealed-vault SDK access

For services that start before a human is available to unseal (e.g. a Java microservice in production), use a system token:

```bash
rune system-token create payment-service-prod -n payment-service
# → stk_xyz789...
```

The service uses this token in its SDK config. Even if the vault is sealed, the service can read its secrets because system tokens carry their own separate encryption layer.

---

## Java Spring Boot SDK

Rune ships with a Spring Boot starter that injects secrets directly into your Spring `Environment`, making them available via `@Value`, `application.properties` references, and `Environment.getProperty()`.

### Add the dependency

```xml
<dependency>
  <groupId>io.rune</groupId>
  <artifactId>rune-spring-boot-starter</artifactId>
  <version>0.1.0</version>
</dependency>
```

### Configure in `application.properties`

```properties
rune.server-url=http://localhost:2710
rune.token=stk_a1b2c3...
rune.namespace=production
rune.enabled=true
```

### Use secrets as normal Spring properties

```java
@Value("${DATABASE_URL}")
private String databaseUrl;

@Value("${STRIPE_KEY}")
private String stripeKey;
```

The starter fetches all secrets from the configured namespace at startup and injects them into the Spring `Environment`. No code changes are needed in your existing services — just swap `.env` for the Rune starter.

> **Tip:** Use a **system token** (`stk_`) for Spring Boot services so the app can start even if the Rune vault is sealed. System tokens are read-only and scoped to one namespace.

---

## HTTP API Quick Reference

All endpoints live at `http://localhost:2710`. Protected endpoints require `Authorization: Bearer <token>`.

| Method | Path | Description |
|---|---|---|
| `GET` | `/status` | Vault sealed/unsealed status |
| `POST` | `/unseal` | Unseal the vault |
| `POST` | `/seal` | Seal the vault |
| `GET` | `/secret/list?namespace=X` | List secret keys |
| `GET` | `/secret/get?namespace=X&key=Y` | Get a secret value |
| `GET` | `/secret/all?namespace=X` | Get all secrets with values |
| `POST` | `/secret/put` | Create or update a secret |
| `POST` | `/secret/delete` | Delete a secret |
| `POST` | `/secret/rotate?namespace=X&key=Y` | Rotate a secret |
| `GET` | `/secret/history?namespace=X&key=Y` | List secret versions |
| `GET` | `/secret/version?namespace=X&key=Y&version=N` | Get a specific version |
| `GET` | `/namespace/list` | List namespaces |
| `POST` | `/namespace/create` | Create a namespace |
| `POST` | `/namespace/delete` | Delete a namespace |
| `GET` | `/token/me` | Get current token info |
| `POST` | `/token/create` | Create a regular token |
| `POST` | `/token/create-system` | Create a system token |
| `POST` | `/token/revoke` | Revoke a token |
| `GET` | `/token/list` | List all tokens |

### Example curl calls

```bash
# Check status
curl http://localhost:2710/status

# Unseal
curl -X POST http://localhost:2710/unseal \
  -H "Content-Type: application/json" \
  -d '{"passphrase":"your-passphrase"}'

# Store a secret
curl -X POST http://localhost:2710/secret/put \
  -H "Authorization: Bearer rnk_..." \
  -H "Content-Type: application/json" \
  -d '{"namespace":"production","key":"DATABASE_URL","value":"postgres://..."}'

# Get a secret
curl "http://localhost:2710/secret/get?namespace=production&key=DATABASE_URL" \
  -H "Authorization: Bearer rnk_..."

# Get all secrets in a namespace (with values)
curl "http://localhost:2710/secret/all?namespace=production" \
  -H "Authorization: Bearer rnk_..."
```

---

## Security Model

| Layer | Mechanism |
|---|---|
| **Encryption at rest** | AES-256-GCM with a unique nonce per secret |
| **Key derivation** | Argon2id (memory-hard) from your passphrase + a random salt |
| **Key storage** | In-memory only, never written to disk |
| **Transport** | HTTP (add a reverse proxy with TLS for production) |
| **Authentication** | Bearer tokens — SHA-256 hashed before storage, never stored in plaintext |
| **Token revocation** | Immediate — revoked tokens are rejected on every request |
| **Sealed access** | System tokens use a second encryption layer keyed to the token itself |
| **Namespace isolation** | Enforced in both the service layer and storage key format |

### What "sealed" actually protects against

When sealed, the AES key is wiped from memory. An attacker who gains access to `~/.rune/rune.db` sees BoltDB buckets with ciphertext blobs. Without the key — derived from your passphrase — the data is unreadable. The passphrase itself is never stored anywhere.

### What Rune does not protect against

- An attacker with access to the running process (they could dump memory)
- A compromised `rune-server` binary
- Unencrypted HTTP traffic on a network (add TLS via nginx or Caddy for remote deployments)

---

## Tips & Best Practices

**Seal before stepping away.** One command wipes the key from memory and protects everything until you unseal again.

**Never commit tokens.** Treat `stk_`, `tkn_`, and `rnk_` tokens like passwords. Use environment variables or secret managers to inject them into CI.

**Use namespaces to scope damage.** If a token leaks, the attacker only gets what's in its namespace(s). Don't dump everything into `default`.

**Rotate regularly.** `rune rotate` generates a cryptographically random value and keeps the old one in history. Prefer rotated secrets to manually chosen ones for API keys and passwords.

**One token per service.** Each service gets its own token with only the permissions it needs. A read-only microservice should never hold a `write` or `delete` token.

**Back up `~/.rune/rune.db`.** The database contains all your secrets, encrypted. Back it up regularly. Without the passphrase it's useless to an attacker, but without the backup it's useless to you too.

**For teams, put Rune behind TLS.** The server itself speaks HTTP. For shared team deployments, put Caddy or nginx in front:

```
# Caddyfile example
vault.yourteam.internal {
    reverse_proxy localhost:2710
}
```

**Use system tokens for always-on services.** If a service restarts at 3 AM, no one is available to unseal the vault. System tokens let the service start regardless.

**Use the JSON editor for bulk work.** Instead of running `rune put` twenty times when onboarding a new environment, paste your `.env` file into the Import dialog or switch to JSON mode and paste directly into the editor.
