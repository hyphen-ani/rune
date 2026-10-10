import { useState, useEffect, useCallback } from 'react'
import { Plus, Trash2, Lock, Copy, Check, KeySquare, Bot } from 'lucide-react'
import { api, type TokenRecord } from '@/lib/api'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Badge } from '@/components/ui/badge'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter } from '@/components/ui/dialog'
import {
  AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent,
  AlertDialogDescription, AlertDialogFooter, AlertDialogHeader,
  AlertDialogTitle, AlertDialogTrigger,
} from '@/components/ui/alert-dialog'

interface Props {
  token: string
  tokenInfo: TokenRecord | null
}

function CopyButton({ text }: { text: string }) {
  const [copied, setCopied] = useState(false)
  const copy = () => {
    navigator.clipboard.writeText(text)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }
  return (
    <Button variant="outline" size="sm" onClick={copy} className="gap-1.5">
      {copied ? <Check className="h-3.5 w-3.5 text-emerald-400" /> : <Copy className="h-3.5 w-3.5" />}
      {copied ? 'Copied' : 'Copy'}
    </Button>
  )
}

export function TokensView({ token, tokenInfo }: Props) {
  const [tokens, setTokens] = useState<TokenRecord[]>([])
  const [namespaces, setNamespaces] = useState<string[]>([])

  // Regular token dialog
  const [addOpen, setAddOpen] = useState(false)
  const [newName, setNewName] = useState('')
  const [newNs, setNewNs] = useState('')
  const [creating, setCreating] = useState(false)

  // System token dialog
  const [sysAddOpen, setSysAddOpen] = useState(false)
  const [sysNewName, setSysNewName] = useState('')
  const [sysNewNs, setSysNewNs] = useState('')
  const [sysCreating, setSysCreating] = useState(false)

  // Reveal modal — shared between regular and system token creation
  const [revealedToken, setRevealedToken] = useState<string | null>(null)
  const [revealedIsSystem, setRevealedIsSystem] = useState(false)

  const isRoot = tokenInfo?.namespace === '*'

  const loadNamespaces = useCallback(async () => {
    const list = await api.listNamespaces(token).catch(() => [])
    const ns = Array.isArray(list) ? list : []
    setNamespaces(ns)
    setNewNs(prev => prev || ns[0] || '')
    setSysNewNs(prev => prev || ns[0] || '')
  }, [token])

  useEffect(() => { if (isRoot) loadNamespaces() }, [isRoot, loadNamespaces])

  const load = useCallback(async () => {
    const list = await api.listTokens(token).catch(() => [])
    setTokens(Array.isArray(list) ? list : [])
  }, [token])

  useEffect(() => { if (isRoot) load() }, [isRoot, load])

  if (!isRoot) {
    return (
      <div className="flex h-full items-center justify-center">
        <div className="text-center space-y-2">
          <Lock className="h-8 w-8 text-muted-foreground/40 mx-auto" />
          <p className="text-sm text-muted-foreground">Root access required</p>
        </div>
      </div>
    )
  }

  async function handleCreate() {
    const name = newName.trim()
    const ns = newNs.trim()
    if (!name || !ns) return
    setCreating(true)
    await api.createToken(token, name, ns)
      .then(r => {
        setRevealedToken(r.token)
        setRevealedIsSystem(false)
        setNewName('')
        setNewNs(namespaces[0] ?? '')
        setAddOpen(false)
        load()
      })
      .catch(e => toast.error('Create failed', { description: String(e) }))
    setCreating(false)
  }

  async function handleCreateSystem() {
    const name = sysNewName.trim()
    const ns = sysNewNs.trim()
    if (!name || !ns) return
    setSysCreating(true)
    await api.createSystemToken(token, name, ns)
      .then(r => {
        setRevealedToken(r.token)
        setRevealedIsSystem(true)
        setSysNewName('')
        setSysNewNs(namespaces[0] ?? '')
        setSysAddOpen(false)
        load()
      })
      .catch(e => toast.error('Create failed', { description: String(e) }))
    setSysCreating(false)
  }

  async function handleRevoke(id: string) {
    await api.revokeToken(token, id)
      .then(() => { toast.success(`Revoked ${id}`); load() })
      .catch(e => toast.error('Revoke failed', { description: String(e) }))
  }

  function fmt(date: string) {
    return new Date(date).toLocaleDateString(undefined, { month: 'short', day: 'numeric', year: 'numeric' })
  }

  return (
    <div className="p-6 space-y-6">
      {/* Page header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-xl font-semibold tracking-tight">Tokens</h1>
          <p className="text-sm text-muted-foreground mt-0.5">
            {tokens.length} token{tokens.length !== 1 ? 's' : ''}
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button size="sm" variant="outline" onClick={() => setSysAddOpen(true)}>
            <Bot className="h-4 w-4" />
            System Token
          </Button>
          <Button size="sm" onClick={() => setAddOpen(true)}>
            <Plus className="h-4 w-4" />
            Add Token
          </Button>
        </div>
      </div>

      {/* Table */}
      <div className="rounded-lg border border-border bg-card overflow-hidden">
        <Table>
          <TableHeader>
            <TableRow className="hover:bg-transparent">
              <TableHead>Name</TableHead>
              <TableHead>Namespace</TableHead>
              <TableHead>Type</TableHead>
              <TableHead>Created</TableHead>
              <TableHead>Status</TableHead>
              <TableHead className="w-16 text-right">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {tokens.length === 0 && (
              <TableRow>
                <TableCell colSpan={6} className="py-16">
                  <div className="flex flex-col items-center gap-3 text-center">
                    <div className="rounded-full bg-muted p-3">
                      <KeySquare className="h-5 w-5 text-muted-foreground" />
                    </div>
                    <div>
                      <p className="text-sm font-medium">No tokens yet</p>
                      <p className="text-xs text-muted-foreground mt-0.5">Create tokens to grant scoped access to namespaces</p>
                    </div>
                    <Button size="sm" variant="outline" onClick={() => setAddOpen(true)}>
                      <Plus className="h-4 w-4" />
                      Add Token
                    </Button>
                  </div>
                </TableCell>
              </TableRow>
            )}
            {tokens.map(t => (
              <TableRow key={t.id}>
                <TableCell className="text-sm font-medium">{t.name}</TableCell>
                <TableCell>
                  <Badge variant="outline" className="font-mono text-xs">
                    {t.namespace === '*' ? 'root' : t.namespace}
                  </Badge>
                </TableCell>
                <TableCell>
                  {t.system ? (
                    <Badge variant="secondary" className="text-xs gap-1 text-sky-400 border-sky-400/30 bg-sky-400/10">
                      <Bot className="h-3 w-3" />
                      system
                    </Badge>
                  ) : (
                    <Badge variant="outline" className="text-xs text-muted-foreground">
                      regular
                    </Badge>
                  )}
                </TableCell>
                <TableCell className="text-sm text-muted-foreground">{fmt(t.created_at)}</TableCell>
                <TableCell>
                  <Badge variant={t.revoked ? 'destructive' : 'secondary'} className="text-xs">
                    {t.revoked ? 'revoked' : 'active'}
                  </Badge>
                </TableCell>
                <TableCell className="text-right">
                  {!t.revoked && t.id !== 'root' && (
                    <AlertDialog>
                      <AlertDialogTrigger asChild>
                        <Button variant="ghost" size="icon" className="h-8 w-8 text-muted-foreground hover:text-destructive">
                          <Trash2 className="h-3.5 w-3.5" />
                        </Button>
                      </AlertDialogTrigger>
                      <AlertDialogContent>
                        <AlertDialogHeader>
                          <AlertDialogTitle>Revoke token</AlertDialogTitle>
                          <AlertDialogDescription>
                            Revoke token <span className="font-medium">{t.name}</span>? It will immediately lose access.
                          </AlertDialogDescription>
                        </AlertDialogHeader>
                        <AlertDialogFooter>
                          <AlertDialogCancel>Cancel</AlertDialogCancel>
                          <AlertDialogAction onClick={() => handleRevoke(t.id)}>Revoke</AlertDialogAction>
                        </AlertDialogFooter>
                      </AlertDialogContent>
                    </AlertDialog>
                  )}
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>

      {/* Add token dialog */}
      <Dialog open={addOpen} onOpenChange={open => { setAddOpen(open); if (!open) { setNewName(''); setNewNs(namespaces[0] ?? '') } }}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>Add Token</DialogTitle>
            <DialogDescription>
              Create a token with scoped access to a namespace.
            </DialogDescription>
          </DialogHeader>
          <div className="space-y-4 py-2">
            <div className="space-y-1.5">
              <Label htmlFor="token-name">Name</Label>
              <Input
                id="token-name"
                placeholder="e.g. ci-deploy"
                value={newName}
                onChange={e => setNewName(e.target.value)}
                className="text-sm"
                autoFocus
              />
            </div>
            <div className="space-y-1.5">
              <Label>Namespace</Label>
              <Select value={newNs} onValueChange={setNewNs}>
                <SelectTrigger className="font-mono text-sm">
                  <SelectValue placeholder="Select a namespace" />
                </SelectTrigger>
                <SelectContent>
                  {namespaces.map(ns => (
                    <SelectItem key={ns} value={ns} className="font-mono text-xs">{ns}</SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setAddOpen(false)}>Cancel</Button>
            <Button onClick={handleCreate} disabled={!newName.trim() || !newNs.trim() || creating}>
              {creating ? 'Creating…' : 'Create Token'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Add system token dialog */}
      <Dialog open={sysAddOpen} onOpenChange={open => { setSysAddOpen(open); if (!open) { setSysNewName(''); setSysNewNs(namespaces[0] ?? '') } }}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2">
              <Bot className="h-4 w-4 text-sky-400" />
              Add System Token
            </DialogTitle>
            <DialogDescription>
              Read-only SDK token. Can read secrets even when the vault is sealed. Use with the Rune Spring Boot starter or other SDKs.
            </DialogDescription>
          </DialogHeader>
          <div className="space-y-4 py-2">
            <div className="space-y-1.5">
              <Label htmlFor="sys-token-name">Name</Label>
              <Input
                id="sys-token-name"
                placeholder="e.g. spring-prod"
                value={sysNewName}
                onChange={e => setSysNewName(e.target.value)}
                className="text-sm"
                autoFocus
              />
            </div>
            <div className="space-y-1.5">
              <Label>Namespace</Label>
              <Select value={sysNewNs} onValueChange={setSysNewNs}>
                <SelectTrigger className="font-mono text-sm">
                  <SelectValue placeholder="Select a namespace" />
                </SelectTrigger>
                <SelectContent>
                  {namespaces.map(ns => (
                    <SelectItem key={ns} value={ns} className="font-mono text-xs">{ns}</SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <p className="text-xs text-muted-foreground rounded-md bg-muted px-3 py-2">
              This token is <span className="font-medium">read-only</span> — it cannot write, delete, or rotate secrets.
            </p>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setSysAddOpen(false)}>Cancel</Button>
            <Button onClick={handleCreateSystem} disabled={!sysNewName.trim() || !sysNewNs.trim() || sysCreating}>
              {sysCreating ? 'Creating…' : 'Create System Token'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Token reveal modal (one-time) — shared for both token types */}
      <Dialog open={!!revealedToken} onOpenChange={open => !open && setRevealedToken(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{revealedIsSystem ? 'System token created' : 'Token created'}</DialogTitle>
            <DialogDescription>
              {revealedIsSystem
                ? 'Copy this token and set it as RUNE_SYSTEM_TOKEN in your SDK config. It will not be shown again.'
                : 'Copy this token now — it will not be shown again.'}
            </DialogDescription>
          </DialogHeader>
          <div className="space-y-3">
            <code className="block font-mono text-xs bg-muted rounded p-3 break-all">
              {revealedToken}
            </code>
            <CopyButton text={revealedToken ?? ''} />
          </div>
        </DialogContent>
      </Dialog>
    </div>
  )
}
