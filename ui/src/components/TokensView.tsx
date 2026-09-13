import { useState, useEffect, useCallback } from 'react'
import { Plus, Trash2, Lock, Copy, Check, KeySquare } from 'lucide-react'
import { api, type TokenRecord } from '@/lib/api'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
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
  const [addOpen, setAddOpen] = useState(false)
  const [newName, setNewName] = useState('')
  const [newNs, setNewNs] = useState('default')
  const [creating, setCreating] = useState(false)
  const [revealedToken, setRevealedToken] = useState<string | null>(null)

  const isRoot = tokenInfo?.namespace === '*'

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
        setNewName('')
        setNewNs('default')
        setAddOpen(false)
        load()
      })
      .catch(e => toast.error('Create failed', { description: String(e) }))
    setCreating(false)
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
        <Button size="sm" onClick={() => setAddOpen(true)}>
          <Plus className="h-4 w-4" />
          Add Token
        </Button>
      </div>

      {/* Table */}
      <div className="rounded-lg border border-border bg-card">
        <Table>
          <TableHeader>
            <TableRow className="hover:bg-transparent">
              <TableHead>Name</TableHead>
              <TableHead>Namespace</TableHead>
              <TableHead>Created</TableHead>
              <TableHead>Status</TableHead>
              <TableHead className="w-16 text-right">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {tokens.length === 0 && (
              <TableRow>
                <TableCell colSpan={5} className="py-16">
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
      <Dialog open={addOpen} onOpenChange={open => { setAddOpen(open); if (!open) { setNewName(''); setNewNs('default') } }}>
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
              <Label htmlFor="token-ns">Namespace</Label>
              <Input
                id="token-ns"
                placeholder="e.g. production"
                value={newNs}
                onChange={e => setNewNs(e.target.value)}
                onKeyDown={e => e.key === 'Enter' && handleCreate()}
                className="font-mono text-sm"
              />
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

      {/* Token reveal modal (one-time) */}
      <Dialog open={!!revealedToken} onOpenChange={open => !open && setRevealedToken(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Token created</DialogTitle>
            <DialogDescription>
              Copy this token now — it will not be shown again.
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
