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

It provides AES-GCM encrypted storage, a secure seal/unseal lifecycle, and token-based authentication — all in a single binary with zero external dependencies. Whether you're running it locally on your machine or on your own server, Rune keeps your secrets encrypted at rest and inaccessible until explicitly unlocked.

> **Rune is not a replacement for enterprise secrets management.** It's a focused tool for developers who want a simple, auditable, and self-contained solution.

---

## Features

| Feature | Description |
|---|---|
| **AES-GCM Encryption** | Secrets are encrypted at rest using industry-standard authenticated encryption |
| **Argon2 Key Derivation** | Passphrase-based key derivation using Argon2 — resistant to brute-force attacks |
| **Seal / Unseal Lifecycle** | Vault starts sealed on every launch; no secrets are accessible until explicitly unsealed |
| **Token Authentication** | All protected endpoints require a valid token; tokens are hashed before storage |
| **CLI Interface** | Intuitive command-line tool for all vault operations |
| **Zero External Dependencies** | No cloud services, no agents, no sidecars — just a single binary |
| **BoltDB Storage** | Embedded key-value store (the same underlying tech powering etcd) |

---

## Architecture

Rune follows a simple client-server model:

```
┌─────────────┐     HTTP      ┌───────────────────-───┐
│   rune CLI  │ ────────────► │     rune-server       │
│  (client)   │               │  (vault daemon)       │
└─────────────┘               │                       │
                              │  ┌─────────────────┐  │
                              │  │   AES-GCM Vault │  │
                              │  │   (BoltDB)      │  │
                              │  └─────────────────┘  │
                              └─────────────-─────────┘
```

- **`rune-server`** — runs the vault daemon, manages encryption, and exposes an HTTP API
- **`rune`** — the CLI client that communicates with the server over HTTP

Both binaries are self-contained. The server is designed to run locally or on any machine you control.

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

Keys support path-style namespacing (e.g., `db/password`, `api/stripe/key`).

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
  login <token>       Authenticate with the vault and store token locally
  unseal              Unseal the vault using your passphrase
  seal                Seal the vault, making all secrets inaccessible
  status              Display the current seal status of the vault
  put <key> <value>   Store a secret at the given key path
  get <key>           Retrieve the secret stored at the given key path

OPTIONS:
  --help              Show help for any command
  --version           Print the current version
```

### Examples

```bash
# Check vault status before unsealing
rune status

# Store secrets using path-style keys
rune put app/prod/db-url postgres://user:pass@host/db
rune put app/prod/api-key sk-abc123

# Retrieve a secret and pipe it to another command
rune get app/prod/db-url | psql

# Seal after use
rune seal
```

---

## API Reference

The server exposes a local HTTP API. All endpoints return JSON.

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|:---:|
| `POST` | `/unseal` | Unseal the vault with a passphrase | ✗ |
| `GET` | `/status` | Return the current vault seal status | ✗ |
| `POST` | `/secret/put` | Store a secret | ✓ |
| `GET` | `/secret/get` | Retrieve a secret | ✓ |
| `POST` | `/seal` | Seal the vault | ✓ |

### Authentication

Protected endpoints require a token passed in the `Authorization` header:

```
Authorization: Bearer <your-token>
```

### Example Requests

```bash
# Check vault status
curl http://localhost:8200/status

# Unseal
curl -X POST http://localhost:8200/unseal \
  -H "Content-Type: application/json" \
  -d '{"passphrase": "your-passphrase"}'

# Store a secret
curl -X POST http://localhost:8200/secret/put \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"key": "db/password", "value": "my-secret"}'

# Retrieve a secret
curl http://localhost:8200/secret/get?key=db/password \
  -H "Authorization: Bearer <token>"
```

---

## Data Storage

Rune stores all data locally. No data ever leaves your machine unless you choose to run `rune-server` on a remote host.

```
~/.rune/
├── rune.db        # AES-GCM encrypted secrets (BoltDB)
└── config.json    # CLI authentication token
```

`rune.db` is an embedded BoltDB database. All secret values are encrypted before being written. The encryption key is derived at unseal time and held only in memory — it is never written to disk.

---

## Security Model

Rune's security design prioritizes simplicity and auditability over complexity.

**Encryption**
- Secrets are encrypted using **AES-256-GCM**, providing both confidentiality and integrity.
- Each secret is encrypted with a unique nonce to prevent ciphertext reuse.

**Key Derivation**
- The encryption key is derived from your passphrase using **Argon2id**, a memory-hard KDF that resists brute-force and GPU-based attacks.
- The derived key is held in memory only for the duration the vault is unsealed. It is never written to disk.

**Seal / Unseal**
- The vault starts **sealed** on every process launch.
- In a sealed state, no secret can be read or written — the in-memory key is zeroed out.
- Sealing is instantaneous and does not require a restart.

**Authentication**
- Tokens are **hashed before storage** using a one-way function. The server never stores raw tokens.
- All write and read operations on secrets require a valid, authenticated token.

**Threat model:** Rune protects secrets at rest from an attacker with access to the filesystem. It is not designed to protect against a compromised process, root-level access, or memory forensics against a live, unsealed vault.

---

## Roadmap

The following features are planned for future releases:

- [X] Token management (create, revoke, list tokens)
- [ ] Role-based access control (RBAC)
- [X] Namespaces for multi-tenant secret isolation
- [ ] Audit logging with tamper-evident records
- [ ] UI dashboard (Electron)
- [ ] TLS support for remote deployments
- [ ] Secret versioning and history

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
