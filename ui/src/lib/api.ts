export interface NamespaceCount {
  namespace: string
  count: number
}

export interface MonthCount {
  month: string
  count: number
}

export interface VaultStats {
  total_secrets: number
  total_namespaces: number
  rotated_secrets: number
  total_tokens: number
  active_tokens: number
  revoked_tokens: number
  secrets_by_namespace: NamespaceCount[]
  tokens_by_month: MonthCount[]
}

export interface NamespacePermission {
  namespace: string
  permissions: string[]
}

export interface TokenRecord {
  id: string
  name: string
  namespace: string                      // legacy field; present on root and system tokens
  namespaces?: NamespacePermission[]     // fine-grained multi-ns permissions (new tokens)
  created_at: string
  revoked: boolean
  system?: boolean
  read_only?: boolean
}

export const ALL_PERMS = ['read', 'write', 'delete', 'rotate'] as const
export type Permission = typeof ALL_PERMS[number]

export interface SecretVersion {
  version: number
  created_at: string
  updated_at: string
  rotated_at: string
}

const hdr = (token: string) => ({
  headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
})

async function check(r: Response) {
  if (!r.ok) {
    const text = await r.text().catch(() => r.statusText)
    throw new Error(text || `HTTP ${r.status}`)
  }
  return r
}

export const api = {
  status: () => fetch('/status').then(r => r.text()),

  unseal: (passphrase: string) =>
    fetch('/unseal', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ passphrase }),
    }).then(check),

  seal: (token: string) =>
    fetch('/seal', { method: 'POST', ...hdr(token) }).then(check),

  tokenMe: (token: string): Promise<TokenRecord> =>
    fetch('/token/me', hdr(token)).then(check).then(r => r.json()),

  listSecrets: (token: string, ns: string): Promise<string[]> =>
    fetch(`/secret/list?namespace=${encodeURIComponent(ns)}`, hdr(token))
      .then(check).then(r => r.json()),

  getSecret: (token: string, ns: string, key: string): Promise<{ value: string }> =>
    fetch(`/secret/get?namespace=${encodeURIComponent(ns)}&key=${encodeURIComponent(key)}`, hdr(token))
      .then(check).then(r => r.json()),

  putSecret: (token: string, ns: string, key: string, value: string) =>
    fetch('/secret/put', {
      method: 'POST', ...hdr(token),
      body: JSON.stringify({ namespace: ns, key, value }),
    }).then(check),

  deleteSecret: (token: string, ns: string, key: string) =>
    fetch('/secret/delete', {
      method: 'POST', ...hdr(token),
      body: JSON.stringify({ namespace: ns, key }),
    }).then(check),

  rotateSecret: (token: string, ns: string, key: string): Promise<{ value: string }> =>
    fetch(`/secret/rotate?namespace=${encodeURIComponent(ns)}&key=${encodeURIComponent(key)}`, {
      method: 'POST', ...hdr(token),
    }).then(check).then(r => r.json()),

  listVersions: (token: string, ns: string, key: string): Promise<SecretVersion[]> =>
    fetch(`/secret/history?namespace=${encodeURIComponent(ns)}&key=${encodeURIComponent(key)}`, hdr(token))
      .then(check).then(r => r.json()),

  getSecretVersion: (token: string, ns: string, key: string, version: number): Promise<{ value: string }> =>
    fetch(`/secret/version?namespace=${encodeURIComponent(ns)}&key=${encodeURIComponent(key)}&version=${version}`, hdr(token))
      .then(check).then(r => r.json()),

  getAllSecrets: (token: string, ns: string): Promise<Record<string, string>> =>
    fetch(`/secret/all?namespace=${encodeURIComponent(ns)}`, hdr(token))
      .then(check).then(r => r.json()),

  listNamespaces: (token: string): Promise<string[]> =>
    fetch('/namespace/list', hdr(token)).then(check).then(r => r.json()),

  createNamespace: (token: string, name: string) =>
    fetch('/namespace/create', {
      method: 'POST', ...hdr(token),
      body: JSON.stringify({ name }),
    }).then(check),

  deleteNamespace: (token: string, name: string) =>
    fetch('/namespace/delete', {
      method: 'POST', ...hdr(token),
      body: JSON.stringify({ name }),
    }).then(check),

  listTokens: (token: string): Promise<TokenRecord[]> =>
    fetch('/token/list', hdr(token)).then(check).then(r => r.json()),

  createToken: (token: string, name: string, namespaces: NamespacePermission[]): Promise<{ token: string; record: TokenRecord }> =>
    fetch('/token/create', {
      method: 'POST', ...hdr(token),
      body: JSON.stringify({ name, namespaces }),
    }).then(check).then(r => r.json()),

  createSystemToken: (token: string, name: string, namespace: string): Promise<{ token: string; record: TokenRecord }> =>
    fetch('/token/create-system', {
      method: 'POST', ...hdr(token),
      body: JSON.stringify({ name, namespace }),
    }).then(check).then(r => r.json()),

  revokeToken: (token: string, id: string) =>
    fetch('/token/revoke', {
      method: 'POST', ...hdr(token),
      body: JSON.stringify({ id }),
    }).then(check),

  getStats: (token: string): Promise<VaultStats> =>
    fetch('/stats', hdr(token)).then(check).then(r => r.json()),
}
