import { useState, useEffect, useCallback } from 'react'
import { Plus, Eye, RotateCcw, History, Trash2, Copy, Check, KeyRound, RefreshCw, Upload } from 'lucide-react'
import { api, type TokenRecord, type SecretVersion } from '@/lib/api'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import {
  Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter,
} from '@/components/ui/dialog'
import {
  AlertDialog, AlertDialogAction, AlertDialogCancel,
  AlertDialogContent, AlertDialogDescription, AlertDialogFooter,
  AlertDialogHeader, AlertDialogTitle, AlertDialogTrigger,
} from '@/components/ui/alert-dialog'
import {
  Tooltip, TooltipContent, TooltipProvider, TooltipTrigger,
} from '@/components/ui/tooltip'
import { Separator } from '@/components/ui/separator'
import { cn } from '@/lib/utils'

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

// Inline copy button — copies text passed as prop
function CopyButton({ text, className }: { text: string; className?: string }) {
  const [copied, setCopied] = useState(false)
  const copy = () => {
    navigator.clipboard.writeText(text)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }
  return (
    <Button variant="ghost" size="icon" className={cn("h-7 w-7 shrink-0", className)} onClick={copy}>
      {copied ? <Check className="h-3.5 w-3.5 text-emerald-600 dark:text-emerald-400" /> : <Copy className="h-3.5 w-3.5" />}
    </Button>
  )
}

// Fetches a specific version's value on demand then copies it
function VersionCopyButton({
  token, ns, secretKey, version,
}: { token: string; ns: string; secretKey: string; version: number }) {
  const [state, setState] = useState<'idle' | 'loading' | 'done'>('idle')

  const handleCopy = async () => {
    if (state === 'loading') return
    setState('loading')
    try {
      const { value } = await api.getSecretVersion(token, ns, secretKey, version)
      await navigator.clipboard.writeText(value)
      setState('done')
      setTimeout(() => setState('idle'), 2000)
    } catch (e) {
      toast.error('Failed to fetch value', { description: String(e) })
      setState('idle')
    }
  }

  return (
    <Button variant="ghost" size="icon" className="h-7 w-7 shrink-0" onClick={handleCopy} disabled={state === 'loading'}>
      {state === 'done'
        ? <Check className="h-3.5 w-3.5 text-emerald-600 dark:text-emerald-400" />
        : state === 'loading'
        ? <RefreshCw className="h-3.5 w-3.5 animate-spin" />
        : <Copy className="h-3.5 w-3.5" />}
    </Button>
  )
}

function fmtDate(iso: string) {
  if (!iso) return '—'
  const d = new Date(iso)
  return isNaN(d.getTime()) ? iso : d.toLocaleString(undefined, {
    month: 'short', day: 'numeric', year: 'numeric',
    hour: '2-digit', minute: '2-digit',
  })
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
  const [historyLoading, setHistoryLoading] = useState(false)
  const [importOpen, setImportOpen] = useState(false)
  const [importData, setImportData] = useState<Record<string, string> | null>(null)
  const [importNs, setImportNs] = useState('')
  const [importing, setImporting] = useState(false)
  const [importProgress, setImportProgress] = useState({ done: 0, total: 0 })

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
    const r = await api.getSecret(token, ns, key).catch(e => {
      toast.error('Failed to fetch value', { description: String(e) })
      return null
    })
    if (r) setValueModal({ key, value: r.value })
  }

  async function handleRotate(key: string) {
    const r = await api.rotateSecret(token, ns, key).catch(e => {
      toast.error('Rotate failed', { description: String(e) })
      return null
    })
    if (r) {
      toast.success(`Rotated "${key}"`)
      setValueModal({ key, value: r.value })
    }
  }

  async function handleHistory(key: string) {
    setHistoryLoading(true)
    setHistoryModal({ key, versions: [] })
    const versions = await api.listVersions(token, ns, key).catch(() => [])
    setHistoryModal({ key, versions: Array.isArray(versions) ? versions : [] })
    setHistoryLoading(false)
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

  function handleImportFile(e: React.ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0]
    if (!file) return
    const reader = new FileReader()
    reader.onload = (ev) => {
      try {
        const parsed = JSON.parse(ev.target?.result as string)
        if (typeof parsed !== 'object' || Array.isArray(parsed) || parsed === null) {
          toast.error('Invalid format', { description: 'JSON must be a flat key-value object' })
          return
        }
        const data: Record<string, string> = {}
        for (const [k, v] of Object.entries(parsed)) data[k] = String(v)
        setImportData(data)
      } catch {
        toast.error('Invalid JSON', { description: 'Could not parse the selected file' })
      }
    }
    reader.readAsText(file)
  }

  async function handleImport() {
    if (!importData) return
    const entries = Object.entries(importData)
    setImporting(true)
    setImportProgress({ done: 0, total: entries.length })
    const failed: string[] = []
    for (const [key, value] of entries) {
      await api.putSecret(token, importNs, key, value).catch(() => failed.push(key))
      setImportProgress(p => ({ ...p, done: p.done + 1 }))
    }
    setImporting(false)
    if (failed.length === 0) {
      toast.success(`Imported ${entries.length} secret${entries.length !== 1 ? 's' : ''}`)
    } else {
      toast.error(`${failed.length} secret${failed.length !== 1 ? 's' : ''} failed to import`, {
        description: failed.join(', '),
      })
    }
    setImportOpen(false)
    setImportData(null)
    loadKeys()
  }

  if (sealed) {
    return (
      <div className="flex h-full items-center justify-center">
        <div className="text-center space-y-2">
          <p className="text-muted-foreground text-sm">Vault is sealed</p>
        </div>
      </div>
    )
  }

  return (
    <TooltipProvider delayDuration={400}>
      <div className="p-6 space-y-5">
        {/* Page header */}
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-xl font-semibold tracking-tight">Secrets</h1>
            <p className="text-sm text-muted-foreground mt-0.5">
              {loading
                ? 'Loading…'
                : `${keys.length} secret${keys.length !== 1 ? 's' : ''} in `}
              {!loading && <span className="font-mono text-foreground">{ns}</span>}
            </p>
          </div>
          <div className="flex items-center gap-2">
            {namespaces.length > 1 && (
              <Select value={ns} onValueChange={setNs}>
                <SelectTrigger className="h-9 w-44 text-xs">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {namespaces.map(n => <SelectItem key={n} value={n} className="text-xs font-mono">{n}</SelectItem>)}
                </SelectContent>
              </Select>
            )}
            <Button size="sm" variant="outline" onClick={() => { setImportNs(namespaces[0] ?? ns); setImportOpen(true) }}>
              <Upload className="h-4 w-4" />
              Import JSON
            </Button>
            <Button size="sm" onClick={() => setAddOpen(true)}>
              <Plus className="h-4 w-4" />
              Add secret
            </Button>
          </div>
        </div>

        {/* Table */}
        <div className="rounded-lg border border-border bg-card overflow-hidden">
          <Table>
            <TableHeader>
              <TableRow className="hover:bg-transparent border-b border-border">
                <TableHead className="w-[40%] pl-4">Key</TableHead>
                <TableHead className="w-[30%]">Namespace</TableHead>
                <TableHead className="text-right pr-4">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {loading && (
                <TableRow>
                  <TableCell colSpan={3} className="text-center text-muted-foreground text-sm py-16">
                    <RefreshCw className="h-4 w-4 animate-spin mx-auto" />
                  </TableCell>
                </TableRow>
              )}
              {!loading && keys.length === 0 && (
                <TableRow>
                  <TableCell colSpan={3} className="py-20">
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
                        Add secret
                      </Button>
                    </div>
                  </TableCell>
                </TableRow>
              )}
              {keys.map(key => (
                <TableRow key={key} className="group">
                  <TableCell className="pl-4">
                    <span className="font-mono text-sm">{key}</span>
                  </TableCell>
                  <TableCell>
                    <Badge variant="outline" className="font-mono text-xs font-normal">
                      {ns}
                    </Badge>
                  </TableCell>
                  <TableCell className="pr-4">
                    <div className="flex items-center justify-end gap-0.5">
                      <Tooltip>
                        <TooltipTrigger asChild>
                          <Button
                            variant="ghost"
                            size="icon"
                            className="h-8 w-8 opacity-60 group-hover:opacity-100 transition-opacity"
                            onClick={() => handleGet(key)}
                          >
                            <Eye className="h-3.5 w-3.5" />
                          </Button>
                        </TooltipTrigger>
                        <TooltipContent>View value</TooltipContent>
                      </Tooltip>

                      <Tooltip>
                        <TooltipTrigger asChild>
                          <Button
                            variant="ghost"
                            size="icon"
                            className="h-8 w-8 opacity-60 group-hover:opacity-100 transition-opacity"
                            onClick={() => handleRotate(key)}
                          >
                            <RotateCcw className="h-3.5 w-3.5" />
                          </Button>
                        </TooltipTrigger>
                        <TooltipContent>Rotate</TooltipContent>
                      </Tooltip>

                      <Tooltip>
                        <TooltipTrigger asChild>
                          <Button
                            variant="ghost"
                            size="icon"
                            className="h-8 w-8 opacity-60 group-hover:opacity-100 transition-opacity"
                            onClick={() => handleHistory(key)}
                          >
                            <History className="h-3.5 w-3.5" />
                          </Button>
                        </TooltipTrigger>
                        <TooltipContent>Version history</TooltipContent>
                      </Tooltip>

                      <Separator orientation="vertical" className="h-4 mx-1" />

                      <AlertDialog>
                        <Tooltip>
                          <TooltipTrigger asChild>
                            <AlertDialogTrigger asChild>
                              <Button
                                variant="ghost"
                                size="icon"
                                className="h-8 w-8 opacity-60 group-hover:opacity-100 transition-opacity hover:text-destructive hover:bg-destructive/10"
                              >
                                <Trash2 className="h-3.5 w-3.5" />
                              </Button>
                            </AlertDialogTrigger>
                          </TooltipTrigger>
                          <TooltipContent>Delete</TooltipContent>
                        </Tooltip>
                        <AlertDialogContent>
                          <AlertDialogHeader>
                            <AlertDialogTitle>Delete secret</AlertDialogTitle>
                            <AlertDialogDescription>
                              Delete <code className="font-mono bg-muted px-1.5 py-0.5 rounded text-sm">{key}</code>?
                              {' '}All versions will be permanently removed.
                            </AlertDialogDescription>
                          </AlertDialogHeader>
                          <AlertDialogFooter>
                            <AlertDialogCancel>Cancel</AlertDialogCancel>
                            <AlertDialogAction
                              onClick={() => handleDelete(key)}
                              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
                            >
                              Delete
                            </AlertDialogAction>
                          </AlertDialogFooter>
                        </AlertDialogContent>
                      </AlertDialog>
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
              <DialogTitle>Add secret</DialogTitle>
              <DialogDescription>
                Store a new secret in <span className="font-mono font-medium">{ns}</span>.
              </DialogDescription>
            </DialogHeader>
            <div className="space-y-4 py-1">
              <div className="space-y-1.5">
                <Label htmlFor="secret-key" className="text-xs text-muted-foreground">Key</Label>
                <Input
                  id="secret-key"
                  placeholder="DATABASE_URL"
                  value={newKey}
                  onChange={e => setNewKey(e.target.value)}
                  className="font-mono text-sm"
                  autoFocus
                />
              </div>
              <div className="space-y-1.5">
                <Label htmlFor="secret-val" className="text-xs text-muted-foreground">Value</Label>
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
                {putting ? 'Saving…' : 'Save secret'}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>

        {/* View value modal */}
        <Dialog open={!!valueModal} onOpenChange={open => !open && setValueModal(null)}>
          <DialogContent className="sm:max-w-lg">
            <DialogHeader>
              <DialogTitle className="font-mono text-sm">{valueModal?.key}</DialogTitle>
              <DialogDescription>Current value</DialogDescription>
            </DialogHeader>
            <div className="flex items-start gap-2 rounded-lg bg-muted p-3">
              <code className="flex-1 font-mono text-xs leading-relaxed break-all select-all">
                {valueModal?.value}
              </code>
              <CopyButton text={valueModal?.value ?? ''} />
            </div>
          </DialogContent>
        </Dialog>

        {/* History modal */}
        <Dialog open={!!historyModal} onOpenChange={open => !open && setHistoryModal(null)}>
          <DialogContent className="sm:max-w-xl">
            <DialogHeader>
              <DialogTitle className="font-mono text-sm">{historyModal?.key}</DialogTitle>
              <DialogDescription>
                {historyLoading
                  ? 'Loading history…'
                  : `${historyModal?.versions.length ?? 0} version${historyModal?.versions.length !== 1 ? 's' : ''}`}
              </DialogDescription>
            </DialogHeader>

            {historyLoading ? (
              <div className="flex items-center justify-center py-10">
                <RefreshCw className="h-5 w-5 animate-spin text-muted-foreground" />
              </div>
            ) : (
              <div className="space-y-2 max-h-80 overflow-y-auto pr-1">
                {(historyModal?.versions ?? []).length === 0 ? (
                  <p className="text-sm text-muted-foreground text-center py-8">No history found</p>
                ) : (
                  [...(historyModal?.versions ?? [])].reverse().map(v => (
                    <div
                      key={v.version}
                      className="flex items-start gap-3 rounded-lg border border-border bg-card p-3"
                    >
                      {/* Version badge */}
                      <Badge variant="secondary" className="shrink-0 font-mono tabular-nums mt-0.5">
                        v{v.version}
                      </Badge>

                      {/* Timestamps */}
                      <div className="flex-1 min-w-0 space-y-0.5">
                        <div className="text-xs text-muted-foreground">
                          <span className="text-foreground font-medium">Created</span>
                          {' '}{fmtDate(v.created_at)}
                        </div>
                        {v.rotated_at && (
                          <div className="text-xs text-muted-foreground">
                            <span className="text-foreground font-medium">Rotated</span>
                            {' '}{fmtDate(v.rotated_at)}
                          </div>
                        )}
                      </div>

                      {/* Copy value on demand */}
                      <Tooltip>
                        <TooltipTrigger asChild>
                          <span>
                            <VersionCopyButton
                              token={token}
                              ns={ns}
                              secretKey={historyModal!.key}
                              version={v.version}
                            />
                          </span>
                        </TooltipTrigger>
                        <TooltipContent>Copy value</TooltipContent>
                      </Tooltip>
                    </div>
                  ))
                )}
              </div>
            )}
          </DialogContent>
        </Dialog>

        {/* Import JSON dialog */}
        <Dialog
          open={importOpen}
          onOpenChange={open => {
            setImportOpen(open)
            if (!open) { setImportData(null); setImportProgress({ done: 0, total: 0 }) }
          }}
        >
          <DialogContent className="sm:max-w-md">
            <DialogHeader>
              <DialogTitle>Import secrets from JSON</DialogTitle>
              <DialogDescription>
                Upload a <span className="font-mono">.json</span> file containing a flat key → value object.
              </DialogDescription>
            </DialogHeader>

            <div className="space-y-4 py-1">
              {/* File picker */}
              <div className="space-y-1.5">
                <Label className="text-xs text-muted-foreground">JSON file</Label>
                <input
                  type="file"
                  accept=".json,application/json"
                  onChange={handleImportFile}
                  className="block w-full text-sm text-muted-foreground file:mr-3 file:rounded-md file:border file:border-border file:bg-muted file:px-3 file:py-1.5 file:text-xs file:font-medium file:text-foreground hover:file:bg-accent cursor-pointer"
                />
              </div>

              {/* Namespace */}
              <div className="space-y-1.5">
                <Label className="text-xs text-muted-foreground">Namespace</Label>
                {isRoot ? (
                  <Select value={importNs} onValueChange={setImportNs}>
                    <SelectTrigger className="h-9 text-xs font-mono">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      {namespaces.map(n => (
                        <SelectItem key={n} value={n} className="text-xs font-mono">{n}</SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                ) : (
                  <Badge variant="outline" className="font-mono text-xs font-normal">{ns}</Badge>
                )}
              </div>

              {/* Preview */}
              {importData && (
                <div className="rounded-lg border border-border bg-muted/40 overflow-hidden">
                  <div className="px-3 py-2 border-b border-border">
                    <p className="text-xs font-medium">
                      Found {Object.keys(importData).length} secret{Object.keys(importData).length !== 1 ? 's' : ''}
                    </p>
                  </div>
                  <div className="max-h-40 overflow-y-auto px-3 py-2 space-y-1">
                    {Object.keys(importData).map(k => (
                      <p key={k} className="text-xs font-mono text-muted-foreground truncate">{k}</p>
                    ))}
                  </div>
                </div>
              )}

              {/* Progress */}
              {importing && (
                <p className="text-xs text-muted-foreground">
                  Importing… {importProgress.done} / {importProgress.total}
                </p>
              )}
            </div>

            <DialogFooter>
              <Button variant="outline" onClick={() => setImportOpen(false)} disabled={importing}>Cancel</Button>
              <Button
                onClick={handleImport}
                disabled={!importData || importing}
              >
                {importing
                  ? `Importing… ${importProgress.done}/${importProgress.total}`
                  : importData
                  ? `Import ${Object.keys(importData).length} secret${Object.keys(importData).length !== 1 ? 's' : ''}`
                  : 'Import'}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>
    </TooltipProvider>
  )
}
