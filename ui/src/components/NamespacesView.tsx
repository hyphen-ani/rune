import { useState, useEffect, useCallback } from 'react'
import { Plus, Trash2, Lock, FolderKey } from 'lucide-react'
import { api, type TokenRecord } from '@/lib/api'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import {
  Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter,
} from '@/components/ui/dialog'
import {
  AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent,
  AlertDialogDescription, AlertDialogFooter, AlertDialogHeader,
  AlertDialogTitle, AlertDialogTrigger,
} from '@/components/ui/alert-dialog'

interface Props {
  token: string
  tokenInfo: TokenRecord | null
}

export function NamespacesView({ token, tokenInfo }: Props) {
  const [namespaces, setNamespaces] = useState<string[]>([])
  const [addOpen, setAddOpen] = useState(false)
  const [newName, setNewName] = useState('')
  const [creating, setCreating] = useState(false)

  const isRoot = tokenInfo?.namespace === '*'

  const load = useCallback(async () => {
    const list = await api.listNamespaces(token).catch(() => [])
    setNamespaces(Array.isArray(list) ? list : [])
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
    if (!name) return
    setCreating(true)
    await api.createNamespace(token, name)
      .then(() => {
        toast.success(`Created "${name}"`)
        setNewName('')
        setAddOpen(false)
        load()
      })
      .catch(e => toast.error('Create failed', { description: String(e) }))
    setCreating(false)
  }

  async function handleDelete(name: string) {
    await api.deleteNamespace(token, name)
      .then(() => { toast.success(`Deleted "${name}"`); load() })
      .catch(e => toast.error('Delete failed', { description: String(e) }))
  }

  return (
    <div className="p-6 space-y-6">
      {/* Page header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-xl font-semibold tracking-tight">Namespaces</h1>
          <p className="text-sm text-muted-foreground mt-0.5">
            {namespaces.length} namespace{namespaces.length !== 1 ? 's' : ''}
          </p>
        </div>
        <Button size="sm" onClick={() => setAddOpen(true)}>
          <Plus className="h-4 w-4" />
          Add Namespace
        </Button>
      </div>

      {/* Table */}
      <div className="rounded-lg border border-border bg-card">
        <Table>
          <TableHeader>
            <TableRow className="hover:bg-transparent">
              <TableHead>Name</TableHead>
              <TableHead className="w-16 text-right">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {namespaces.length === 0 && (
              <TableRow>
                <TableCell colSpan={2} className="py-16">
                  <div className="flex flex-col items-center gap-3 text-center">
                    <div className="rounded-full bg-muted p-3">
                      <FolderKey className="h-5 w-5 text-muted-foreground" />
                    </div>
                    <div>
                      <p className="text-sm font-medium">No namespaces yet</p>
                      <p className="text-xs text-muted-foreground mt-0.5">Create a namespace to isolate secrets by team or environment</p>
                    </div>
                    <Button size="sm" variant="outline" onClick={() => setAddOpen(true)}>
                      <Plus className="h-4 w-4" />
                      Add Namespace
                    </Button>
                  </div>
                </TableCell>
              </TableRow>
            )}
            {namespaces.map(ns => (
              <TableRow key={ns}>
                <TableCell className="font-mono text-sm">{ns}</TableCell>
                <TableCell className="text-right">
                  <AlertDialog>
                    <AlertDialogTrigger asChild>
                      <Button variant="ghost" size="icon" className="h-8 w-8 text-muted-foreground hover:text-destructive">
                        <Trash2 className="h-3.5 w-3.5" />
                      </Button>
                    </AlertDialogTrigger>
                    <AlertDialogContent>
                      <AlertDialogHeader>
                        <AlertDialogTitle>Delete namespace</AlertDialogTitle>
                        <AlertDialogDescription>
                          Delete <code className="font-mono bg-muted px-1 rounded">{ns}</code>? All secrets inside will be lost.
                        </AlertDialogDescription>
                      </AlertDialogHeader>
                      <AlertDialogFooter>
                        <AlertDialogCancel>Cancel</AlertDialogCancel>
                        <AlertDialogAction onClick={() => handleDelete(ns)}>Delete</AlertDialogAction>
                      </AlertDialogFooter>
                    </AlertDialogContent>
                  </AlertDialog>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>

      {/* Add namespace dialog */}
      <Dialog open={addOpen} onOpenChange={open => { setAddOpen(open); if (!open) setNewName('') }}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>Add Namespace</DialogTitle>
            <DialogDescription>
              Namespaces isolate secrets by team, service, or environment.
            </DialogDescription>
          </DialogHeader>
          <div className="space-y-1.5 py-2">
            <Label htmlFor="ns-name">Name</Label>
            <Input
              id="ns-name"
              placeholder="e.g. production"
              value={newName}
              onChange={e => setNewName(e.target.value)}
              onKeyDown={e => e.key === 'Enter' && handleCreate()}
              className="font-mono text-sm"
              autoFocus
            />
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setAddOpen(false)}>Cancel</Button>
            <Button onClick={handleCreate} disabled={!newName.trim() || creating}>
              {creating ? 'Creating…' : 'Create Namespace'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
