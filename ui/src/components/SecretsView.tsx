import { useState, useEffect, useCallback } from 'react'
import { Plus, MoreHorizontal, Eye, RotateCcw, History, Trash2, Copy, Check, KeyRound } from 'lucide-react'
import { api, type TokenRecord, type SecretVersion } from '@/lib/api'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import {
  DropdownMenu, DropdownMenuContent, DropdownMenuItem,
  DropdownMenuSeparator, DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import {
  Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter,
} from '@/components/ui/dialog'
import {
  AlertDialog, AlertDialogAction, AlertDialogCancel,
  AlertDialogContent, AlertDialogDescription, AlertDialogFooter,
  AlertDialogHeader, AlertDialogTitle, AlertDialogTrigger,
} from '@/components/ui/alert-dialog'

interface Props {
  token: string
  tokenInfo: TokenRecord | null
  sealed: boolean
}

interface ValueModal {
  key: string
  value: string
}

interface HistoryModal {
  key: string
  versions: SecretVersion[]
}

function CopyButton({ text }: { text: string }) {
  const [copied, setCopied] = useState(false)
  const copy = () => {
    navigator.clipboard.writeText(text)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }
  return (
    <Button variant="ghost" size="icon" className="h-7 w-7 shrink-0" onClick={copy}>
      {copied ? <Check className="h-3.5 w-3.5 text-emerald-400" /> : <Copy className="h-3.5 w-3.5" />}
    </Button>
  )
}

export function SecretsView({ token, tokenInfo, sealed }: Props) {
  const [namespaces, setNamespaces] = useState<string[]>([])
  const [ns, setNs] = useState('')
  const [keys, setKeys] = useState<string[]>([])
  const [loading, setLoading] = useState(false)
  const [addOpen, setAddOpen] = useState(false)
  const [newKey, setNewKey] = useState('')
  const [newVal, setNewVal] = useState('')
  const [putting, setPutting] = useState(false)
  const [valueModal, setValueModal] = useState<ValueModal | null>(null)
  const [historyModal, setHistoryModal] = useState<HistoryModal | null>(null)

  const isRoot = tokenInfo?.namespace === '*'

  const loadNamespaces = useCallback(async () => {
    if (isRoot) {
      const list = await api.listNamespaces(token).catch(() => ['default'])
      setNamespaces(list)
      if (!ns) setNs(list[0] ?? 'default')
    } else {
      const n = tokenInfo?.namespace ?? 'default'
      setNamespaces([n])
      setNs(n)
    }
  }, [token, tokenInfo, isRoot, ns])

  const loadKeys = useCallback(async () => {
    if (!ns) return
    setLoading(true)
    const list = await api.listSecrets(token, ns).catch(() => [])
    setKeys(Array.isArray(list) ? list : [])
    setLoading(false)
  }, [token, ns])

  useEffect(() => { loadNamespaces() }, [loadNamespaces])
  useEffect(() => { if (ns) loadKeys() }, [loadKeys, ns])

  async function handleGet(key: string) {
    const r = await api.getSecret(token, ns, key).catch(e => { toast.error('Error', { description: String(e) }); return null })
    if (r) setValueModal({ key, value: r.value })
  }

  async function handleRotate(key: string) {
    const r = await api.rotateSecret(token, ns, key).catch(e => { toast.error('Rotate failed', { description: String(e) }); return null })
    if (r) setValueModal({ key, value: r.value })
  }

  async function handleHistory(key: string) {
    const versions = await api.listVersions(token, ns, key).catch(() => [])
    setHistoryModal({ key, versions })
  }

  async function handleDelete(key: string) {
    await api.deleteSecret(token, ns, key)
      .then(() => { toast.success(`Deleted "${key}"`); loadKeys() })
      .catch(e => toast.error('Delete failed', { description: String(e) }))
  }

  async function handlePut() {
    if (!newKey.trim() || !newVal.trim()) return
    setPutting(true)
    await api.putSecret(token, ns, newKey.trim(), newVal.trim())
      .then(() => {
        toast.success(`Saved "${newKey.trim()}"`)
        setNewKey('')
        setNewVal('')
        setAddOpen(false)
        loadKeys()
      })
      .catch(e => toast.error('Save failed', { description: String(e) }))
    setPutting(false)
  }

  if (sealed) {
    return (
      <div className="flex h-full items-center justify-center">
        <div className="text-center space-y-2">
          <p className="text-muted-foreground text-sm">Vault is sealed</p>
          <p className="text-xs text-muted-foreground">Go to login to unseal</p>
        </div>
      </div>
    )
  }

  return (
    <div className="p-6 space-y-6">
      {/* Page header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-xl font-semibold tracking-tight">Secrets</h1>
          <p className="text-sm text-muted-foreground mt-0.5">
            {loading ? 'Loading…' : `${keys.length} key${keys.length !== 1 ? 's' : ''}`} in{' '}
            <span className="font-mono">{ns}</span>
          </p>
        </div>
        <div className="flex items-center gap-3">
          {namespaces.length > 1 && (
            <Select value={ns} onValueChange={setNs}>
              <SelectTrigger className="h-9 w-44 text-xs">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {namespaces.map(n => <SelectItem key={n} value={n}>{n}</SelectItem>)}
              </SelectContent>
            </Select>
          )}
          <Button size="sm" onClick={() => setAddOpen(true)}>
            <Plus className="h-4 w-4" />
            Add Secret
          </Button>
        </div>
      </div>

      {/* Table */}
      <div className="rounded-lg border border-border bg-card">
        <Table>
          <TableHeader>
            <TableRow className="hover:bg-transparent">
              <TableHead className="w-[60%]">Key</TableHead>
              <TableHead className="text-right">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {loading && (
              <TableRow>
                <TableCell colSpan={2} className="text-center text-muted-foreground text-sm py-12">
                  Loading…
                </TableCell>
              </TableRow>
            )}
            {!loading && keys.length === 0 && (
              <TableRow>
                <TableCell colSpan={2} className="py-16">
                  <div className="flex flex-col items-center gap-3 text-center">
                    <div className="rounded-full bg-muted p-3">
                      <KeyRound className="h-5 w-5 text-muted-foreground" />
                    </div>
                    <div>
                      <p className="text-sm font-medium">No secrets yet</p>
                      <p className="text-xs text-muted-foreground mt-0.5">Add your first secret to get started</p>
                    </div>
                    <Button size="sm" variant="outline" onClick={() => setAddOpen(true)}>
                      <Plus className="h-4 w-4" />
                      Add Secret
                    </Button>
                  </div>
                </TableCell>
              </TableRow>
            )}
            {keys.map(key => (
              <TableRow key={key}>
                <TableCell className="font-mono text-sm">{key}</TableCell>
                <TableCell className="text-right">
                  <div className="flex items-center justify-end gap-1">
                    <Button variant="ghost" size="sm" className="h-8 px-3 text-xs" onClick={() => handleGet(key)}>
                      <Eye className="h-3.5 w-3.5 mr-1" />View
                    </Button>
                    <DropdownMenu>
                      <DropdownMenuTrigger asChild>
                        <Button variant="ghost" size="icon" className="h-8 w-8">
                          <MoreHorizontal className="h-4 w-4" />
                        </Button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent align="end">
                        <DropdownMenuItem onClick={() => handleRotate(key)}>
                          <RotateCcw className="h-3.5 w-3.5" />Rotate
                        </DropdownMenuItem>
                        <DropdownMenuItem onClick={() => handleHistory(key)}>
                          <History className="h-3.5 w-3.5" />History
                        </DropdownMenuItem>
                        <DropdownMenuSeparator />
                        <AlertDialog>
                          <AlertDialogTrigger asChild>
                            <DropdownMenuItem
                              className="text-destructive focus:text-destructive"
                              onSelect={e => e.preventDefault()}
                            >
                              <Trash2 className="h-3.5 w-3.5" />Delete
                            </DropdownMenuItem>
                          </AlertDialogTrigger>
                          <AlertDialogContent>
                            <AlertDialogHeader>
                              <AlertDialogTitle>Delete secret</AlertDialogTitle>
                              <AlertDialogDescription>
                                Delete <code className="font-mono bg-muted px-1 rounded">{key}</code>? This cannot be undone.
                              </AlertDialogDescription>
                            </AlertDialogHeader>
                            <AlertDialogFooter>
                              <AlertDialogCancel>Cancel</AlertDialogCancel>
                              <AlertDialogAction onClick={() => handleDelete(key)}>Delete</AlertDialogAction>
                            </AlertDialogFooter>
                          </AlertDialogContent>
                        </AlertDialog>
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </div>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>

      {/* Add secret dialog */}
      <Dialog open={addOpen} onOpenChange={open => { setAddOpen(open); if (!open) { setNewKey(''); setNewVal('') } }}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>Add Secret</DialogTitle>
            <DialogDescription>
              Store a new secret in <span className="font-mono font-medium">{ns}</span>.
            </DialogDescription>
          </DialogHeader>
          <div className="space-y-4 py-2">
            <div className="space-y-1.5">
              <Label htmlFor="secret-key">Key</Label>
              <Input
                id="secret-key"
                placeholder="e.g. DATABASE_URL"
                value={newKey}
                onChange={e => setNewKey(e.target.value)}
                className="font-mono text-sm"
                autoFocus
              />
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="secret-val">Value</Label>
              <Input
                id="secret-val"
                type="password"
                placeholder="Secret value"
                value={newVal}
                onChange={e => setNewVal(e.target.value)}
                onKeyDown={e => e.key === 'Enter' && handlePut()}
              />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setAddOpen(false)}>Cancel</Button>
            <Button onClick={handlePut} disabled={!newKey.trim() || !newVal.trim() || putting}>
              {putting ? 'Saving…' : 'Save Secret'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Value modal */}
      <Dialog open={!!valueModal} onOpenChange={open => !open && setValueModal(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle className="font-mono text-sm">{valueModal?.key}</DialogTitle>
            <DialogDescription>Secret value</DialogDescription>
          </DialogHeader>
          <div className="flex items-center gap-2">
            <code className="flex-1 font-mono text-xs bg-muted rounded p-3 break-all">
              {valueModal?.value}
            </code>
            <CopyButton text={valueModal?.value ?? ''} />
          </div>
        </DialogContent>
      </Dialog>

      {/* History modal */}
      <Dialog open={!!historyModal} onOpenChange={open => !open && setHistoryModal(null)}>
        <DialogContent className="max-w-lg">
          <DialogHeader>
            <DialogTitle className="font-mono text-sm">{historyModal?.key} — History</DialogTitle>
            <DialogDescription>{historyModal?.versions.length ?? 0} versions</DialogDescription>
          </DialogHeader>
          <div className="space-y-2 max-h-72 overflow-auto">
            {historyModal?.versions.map(v => (
              <div key={v.version} className="flex items-start gap-3 rounded border border-border p-3">
                <Badge variant="secondary" className="shrink-0">v{v.version}</Badge>
                <code className="flex-1 font-mono text-xs break-all">{v.value}</code>
                <CopyButton text={v.value} />
              </div>
            ))}
            {historyModal?.versions.length === 0 && (
              <p className="text-sm text-muted-foreground text-center py-4">No history</p>
            )}
          </div>
        </DialogContent>
      </Dialog>
    </div>
  )
}
