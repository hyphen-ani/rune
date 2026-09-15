<p align="center">
  <img src="./assets/logo.png" alt="Rune Logo" width="200"/>
</p>

<h1 align="center">Rune</h1>

<p align="center">
  <strong>A minimal, secure, self-hosted secrets manager.</strong><br/>
  Lightweight. Encrypted. Local-first.
</p>

<p align="center">
  <img src="https://img.shields.io/badge/version-v0.1.0-blue" alt="version"/>
  <img src="https://img.shields.io/badge/license-MIT-green" alt="license"/>
  <img src="https://img.shields.io/badge/go-%3E%3D1.21-00ADD8?logo=go" alt="go version"/>
  <img src="https://img.shields.io/badge/platform-macOS%20%7C%20Linux%20%7C%20Windows-blue" alt="platform"/>
</p>

---

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Architecture](#architecture)
- [Installation](#installation)
- [Getting Started](#getting-started)
- [CLI Reference](#cli-reference)
- [API Reference](#api-reference)
- [Data Storage](#data-storage)
- [Security Model](#security-model)
- [Roadmap](#roadmap)
- [Contributing](#contributing)
- [Philosophy](#philosophy)
- [License](#license)

---

## Overview

Rune is a lightweight, self-hosted secrets manager built for developers who want full control over their secrets without the overhead of complex infrastructure.

It provides AES-256-GCM encrypted storage, a secure seal/unseal lifecycle, namespace-based isolation, secret versioning, and token-based authentication — all in a single binary with zero external dependencies. Access your vault through the built-in web UI or the `rune` CLI.

> **Rune is not a replacement for enterprise secrets management.** It's a focused tool for developers who want a simple, auditable, and self-contained solution.

---

## Features

| Feature | Description |
|---|---|
| **AES-256-GCM Encryption** | Secrets are encrypted at rest using industry-standard authenticated encryption |
| **Argon2id Key Derivation** | Passphrase-based key derivation — resistant to brute-force and GPU attacks |
| **Seal / Unseal Lifecycle** | Vault starts sealed on every launch; no secrets accessible until explicitly unlocked |
| **Token Authentication** | All protected endpoints require a valid token; tokens are SHA-256 hashed before storage |
| **Namespaces** | Isolate secrets by team, service, or environment with fine-grained token access |
| **Secret Versioning** | Full history and point-in-time retrieval for every secret |
| **Secret Rotation** | Rotate secret values on demand; previous versions are preserved |
| **Web UI** | Built-in React dashboard for managing secrets, tokens, namespaces, and vault stats |
| **JSON Import** | Bulk-import secrets from a `.json` file directly into any namespace |
| **CLI Interface** | Full-featured command-line tool for all vault operations |
| **Zero External Dependencies** | No cloud services, no agents, no sidecars — just a single binary |
| **BoltDB Storage** | Embedded key-value store — the same underlying engine powering etcd |

---

## Architecture

Rune ships two binaries that work together:

```
┌──────────────────┐              ┌───────────────────────────┐
│   rune CLI       │   HTTP API   │       rune-server         │
│   (client)       │ ───────────► │       :8080               │
└──────────────────┘              │                           │
                                  │  ┌─────────────────────┐  │
┌──────────────────┐              │  │   AES-256-GCM Vault │  │
│   Web Browser    │   HTTP API   │  │   (BoltDB)          │  │
│   /ui/           │ ───────────► │  └─────────────────────┘  │
└──────────────────┘              └───────────────────────────┘
```

- **`rune-server`** — the vault daemon. Manages encryption, exposes the HTTP API, and serves the embedded web UI at `/ui/`
- **`rune`** — the CLI client that communicates with the server over HTTP
- **Web UI** — a React + Tailwind dashboard embedded in the server binary, accessible in any browser

---

## Installation

### Homebrew (macOS / Linux)

```bash
brew install runelock
```

### Verify the installation

```bash
rune --version
rune-server --version
```

> Binary releases for other platforms are available on the [Releases](https://github.com/hyphen-ani/rune/releases) page.

---

## Getting Started

### Step 1 — Start the server

```bash
rune-server
```

On first launch, Rune will prompt you to set a passphrase and generate a root token:

```
? Set a passphrase: ••••••••••••••

✓ Vault initialized.
✓ Root token generated.
  Root Token (SAVE THIS): rune.xxxxxxxxxxxxxxxxxxxxxxxx
  Store this token securely. It will not be shown again.
```

> ⚠️ **Save your root token immediately.** It cannot be recovered after this point.

The web UI is available at **[http://localhost:8080/ui/](http://localhost:8080/ui/)** once the server is running.

---

### Step 2 — Login

Authenticate the CLI with your root token:

```bash
rune login <token>
```

The token is stored locally at `~/.rune/config.json`.

---

### Step 3 — Unseal the vault

The vault starts in a **sealed** state on every launch. Unseal it with your passphrase:

```bash
rune unseal
```

```
? Enter passphrase: ••••••••••••••
✓ Vault unsealed.
```

---

### Step 4 — Store a secret

```bash
rune put db/password my-secret-value
```

---

### Step 5 — Retrieve a secret

```bash
rune get db/password
```

```
my-secret-value
```

---

### Step 6 — Seal the vault

When you're done, seal the vault to make all secrets inaccessible:

```bash
rune seal
```

```
✓ Vault sealed.
```

---

## CLI Reference

```
USAGE:
  rune <command> [arguments]

COMMANDS:
  login <token>               Authenticate with the vault and store token locally
  unseal                      Unseal the vault using your passphrase
  seal                        Seal the vault, making all secrets inaccessible
  status                      Display the current seal status of the vault
  put <key> <value>           Store a secret (use -n <namespace> to target a namespace)
  get <key>                   Retrieve a secret (use -n <namespace> for a namespace)
  list                        List all secrets (use -n <namespace> for a namespace)
  delete <key>                Delete a secret (use -n <namespace> for a namespace)
  namespace create <name>     Create a new namespace
  namespace list              List all namespaces
  namespace delete <name>     Delete a namespace
  token create <name>         Create a new token (use --namespace to restrict access)
  token list                  List all tokens
  token revoke <id>           Revoke a token by ID

OPTIONS:
  -n, --namespace <name>      Target namespace for secret operations
  --help                      Show help for any command
  --version                   Print the current version
```

### Examples

```bash
# Start the server
rune-server

# Unseal the vault
rune unseal

# Authenticate
rune login rune.xxxxxxxxxxxxxxxxxxxxxxxx

# Secrets — default namespace
rune put db/password my-secret-password
rune put api/key sk_live_abc123
rune get db/password
rune list
rune delete db/password

# Secrets — named namespace
rune namespace create springboot
rune put db/password supersecret -n springboot
rune get db/password -n springboot
rune list -n springboot
rune delete db/password -n springboot
rune namespace delete springboot

# Tokens
rune token create backend-service --namespace springboot
rune token list
rune token revoke <token-id>

# Seal when done
rune seal
```

---

## API Reference

All endpoints listen on `http://localhost:8080`. Protected endpoints require a Bearer token.

```
Authorization: Bearer <your-token>
```

### Public

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/status` | Vault seal status |
| `POST` | `/unseal` | Unseal the vault |
| `GET` | `/health` | Health check |

### Secrets

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/secret/put` | Store or update a secret |
| `GET` | `/secret/get` | Retrieve the latest value of a secret |
| `POST` | `/secret/delete` | Delete a secret |
| `GET` | `/secret/list` | List secrets in a namespace |
| `POST` | `/secret/rotate` | Rotate a secret and preserve the old version |
| `GET` | `/secret/history` | List all versions of a secret |
| `GET` | `/secret/version` | Retrieve a specific version of a secret |

### Namespaces

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/namespace/create` | Create a namespace |
| `GET` | `/namespace/list` | List all namespaces |
| `POST` | `/namespace/delete` | Delete a namespace |

### Tokens

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/token/me` | Get info about the current token |
| `POST` | `/token/create` | Create a new token |
| `POST` | `/token/revoke` | Revoke a token |
| `GET` | `/token/list` | List all tokens (root only) |

### Vault

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/seal` | Seal the vault |
| `GET` | `/stats` | Vault statistics (root only) |

### Example Requests

```bash
# Check vault status
curl http://localhost:8080/status

# Unseal
curl -X POST http://localhost:8080/unseal \
  -H "Content-Type: application/json" \
  -d '{"passphrase": "your-passphrase"}'

# Store a secret
curl -X POST http://localhost:8080/secret/put \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"namespace": "default", "key": "db/password", "value": "my-secret"}'

# Retrieve a secret
curl "http://localhost:8080/secret/get?namespace=default&key=db/password" \
  -H "Authorization: Bearer <token>"

# Rotate a secret
curl -X POST "http://localhost:8080/secret/rotate?namespace=default&key=db/password" \
  -H "Authorization: Bearer <token>"
```

---

## Data Storage

Rune stores all data locally. No data ever leaves your machine unless you choose to run `rune-server` on a remote host.

```
~/.rune/
├── rune.db        # AES-256-GCM encrypted secrets (BoltDB)
└── config.json    # CLI authentication token
```

Secret keys are stored in BoltDB using the format `{namespace}/{key}` for the latest value and `{namespace}/{key}@v{N}` for versioned history. All values are encrypted before being written. The encryption key is derived at unseal time and held only in memory — it is never written to disk.

---

## Security Model

Rune's security design prioritizes simplicity and auditability over complexity.

**Encryption**
- Secrets are encrypted using **AES-256-GCM**, providing both confidentiality and integrity.
- Each secret is encrypted with a unique nonce to prevent ciphertext reuse.

**Key Derivation**
- The encryption key is derived from your passphrase using **Argon2id**, a memory-hard KDF that resists brute-force and GPU-based attacks.
- The derived key is held in memory only while the vault is unsealed. It is never written to disk.

**Seal / Unseal**
- The vault starts **sealed** on every process launch.
- In a sealed state, no secret can be read or written — the in-memory key is zeroed out.
- Sealing is instantaneous and does not require a restart.

**Authentication**
- Tokens are **SHA-256 hashed before storage**. The server never stores raw tokens.
- All read and write operations on secrets require a valid, authenticated token.
- Namespace tokens are restricted to their assigned namespace; only the root token (`namespace: *`) has unrestricted access.

**Threat model:** Rune protects secrets at rest from an attacker with access to the filesystem. It is not designed to protect against a compromised process, root-level access, or memory forensics against a live, unsealed vault.

---

## Roadmap

- [x] Token management (create, revoke, list)
- [x] Namespaces for multi-tenant secret isolation
- [x] Secret versioning and history
- [x] Secret rotation
- [x] Web UI dashboard
- [x] JSON bulk import
- [ ] Audit logging with tamper-evident records
- [ ] TLS support for remote deployments
- [ ] Role-based access control (RBAC)
- [ ] HashiCorp Vault import

---

## Contributing

Contributions are welcome and appreciated.

1. **Fork** the repository
2. **Create** a feature branch: `git checkout -b feat/my-feature`
3. **Commit** your changes: `git commit -m "feat: add my feature"`
4. **Push** to your branch: `git push origin feat/my-feature`
5. **Open** a pull request

Please open an issue first for significant changes or new features so we can discuss the approach before you invest time in the implementation.

### Reporting Issues

If you discover a security vulnerability, **do not open a public issue**. Please email [security@runelock.io](mailto:security@runelock.io) instead.
For bugs and feature requests, use the [GitHub issue tracker](https://github.com/hyphen-ani/rune/issues).

---

## Philosophy

> *Secure systems do not need to be complex.*

Rune exists because most secrets managers either require cloud infrastructure, heavy runtimes, or significant operational overhead. Rune is a single binary. Its source is small enough to read in an afternoon.

It is built with three principles in mind:

**Minimalism** — every feature must justify its existence. Complexity is a liability.
**Control** — your secrets run on your infrastructure. Nothing phones home.
**Clarity** — the security model is simple enough to be understood and audited by a single developer.

Rune is inspired by [HashiCorp Vault](https://www.vaultproject.io/) and [AWS Secrets Manager](https://aws.amazon.com/secrets-manager/), but intentionally scoped for local and small-team use cases where those tools are too heavy.

---

## License

Rune is released under the [MIT License](./LICENSE).

---

<p align="center">
  Built with focus on simplicity and real-world usability.
</p>
