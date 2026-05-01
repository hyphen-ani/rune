# Rune

![version](https://img.shields.io/badge/version-v0.1.0-blue)

<p align="center">
  <img src="./assets/logo.png" alt="Rune Logo" width="220"/>
</p>

<p align="center">
  <strong>A minimal, secure, self-hosted secrets manager.</strong>
</p>

<p align="center">
  Lightweight. Encrypted. Local-first.
</p>

---

## Overview

Rune is a lightweight alternative to traditional secrets managers designed for simplicity and control.

It provides encrypted storage for secrets, a secure seal/unseal lifecycle, and token-based authentication, all without requiring heavy infrastructure or external dependencies.

Rune is built for developers who want a simple, self-hosted solution for managing secrets locally or on their own infrastructure.

---

## Features

- AES-GCM encryption at rest
- Argon2-based key derivation
- Secure seal / unseal workflow
- Token-based authentication
- CLI for interacting with secrets
- Local-first, zero external dependencies
- BoltDB storage (same underlying tech used in etcd)

---

## Architecture


Rune follows a simple client-server model:

- `rune-server` runs the vault
- `rune` CLI interacts with it over HTTP

---

## Installation

### Using Homebrew

```bash
brew install runelock
```

---

## Getting Started

### 1. Start the server

```bash
rune-server
```

On first run:

- You will be prompted to set a passphrase
- A root token will be generated

```
Root Token (SAVE THIS): <token>
```

---

### 2. Login

```bash
rune login <token>
```

Token is stored securely at:

```
~/.rune/config.json
```

---

### 3. Unseal the vault

```bash
rune unseal
```

---

### 4. Store a secret

```bash
rune put db/password my-secret
```

---

### 5. Retrieve a secret

```bash
rune get db/password
```

---

### 6. Seal the vault

```bash
rune seal
```

---

## Data Storage

Rune stores all data locally:

```
~/.rune/
├── rune.db        # encrypted secrets (BoltDB)
├── config.json    # CLI auth token
```

---

## Security Model

Rune follows a simple but strong security model:

- Secrets are encrypted using AES-GCM
- Encryption keys are derived from a passphrase using Argon2
- The vault starts in a sealed state
- No secrets are accessible until unsealed
- Tokens are hashed before being stored
- All protected endpoints require valid authentication

---

## API Overview

| Endpoint        | Description              | Auth Required |
|----------------|--------------------------|--------------|
| `/unseal`      | Unseal vault             | No           |
| `/status`      | Check vault status       | No           |
| `/secret/put`  | Store secret             | Yes          |
| `/secret/get`  | Retrieve secret          | Yes          |
| `/seal`        | Seal vault               | Yes          |

---

## CLI Commands

```
rune login <token>
rune unseal
rune seal
rune status
rune put <key> <value>
rune get <key>
```

---

## Roadmap

- Token management (create, revoke, list)
- Role-based access control (RBAC)
- Namespaces for secrets
- Audit logging
- UI dashboard (Electron)

---

## Contributing

Contributions are welcome.

If you find a bug or want to improve Rune:

1. Fork the repository
2. Create a new branch
3. Submit a pull request

---

## License

MIT License

---

## Philosophy

Rune is designed with a simple principle:

> Secure systems do not need to be complex.

It prioritizes clarity, minimalism, and control over feature bloat.

---

## Acknowledgements

Inspired by tools like:

- HashiCorp Vault
- AWS Secrets Manager
- Redis

---

<p align="center">
  Built with focus on simplicity and real-world usability.
</p>
