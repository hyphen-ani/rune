import { useState } from 'react'
import { Eye, EyeOff, Sun, Moon } from 'lucide-react'
import { useTheme } from 'next-themes'
import logo from '@/assets/logo.png'
import { api } from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { toast } from 'sonner'

interface Props {
  savedToken: string
  onSaveToken: (t: string) => void
  onClearToken: () => void
  onSealed: (sealed: boolean) => void
  sealed: boolean
}

export function LoginView({ savedToken, onSaveToken, onClearToken, onSealed, sealed }: Props) {
  const [tokenInput, setTokenInput] = useState('')
  const [passphrase, setPassphrase] = useState('')
  const [showToken, setShowToken] = useState(false)
  const [saving, setSaving] = useState(false)
  const [unsealing, setUnsealing] = useState(false)
  const { resolvedTheme, setTheme } = useTheme()

  // Phase 1: no token saved yet → ask for token
  // Phase 2: token saved → ask for passphrase (if sealed) or auto-enter (handled by App)
  const phase = savedToken ? 'vault' : 'token'

  async function handleSaveToken() {
    const t = tokenInput.trim()
    if (!t) return
    setSaving(true)
    try {
      await api.tokenMe(t)
      onSaveToken(t)
    } catch {
      toast.error('Invalid token — check and try again')
    } finally {
      setSaving(false)
    }
  }

  async function handleUnseal() {
    const p = passphrase.trim()
    if (!p) return
    setUnsealing(true)
    try {
      await api.unseal(p)
      onSealed(false)
      setPassphrase('')
    } catch (e) {
      toast.error('Wrong passphrase', { description: String(e) })
    } finally {
      setUnsealing(false)
    }
  }

  const maskedToken = savedToken.length > 16
    ? `${savedToken.slice(0, 5)}…${savedToken.slice(-6)}`
    : savedToken

  return (
    <div className="flex h-full w-full items-center justify-center bg-background">
      {/* Theme toggle */}
      <div className="absolute top-4 right-4">
        <Button
          variant="ghost"
          size="icon"
          className="h-8 w-8"
          onClick={() => setTheme(resolvedTheme === 'dark' ? 'light' : 'dark')}
          title="Toggle theme"
        >
          {resolvedTheme === 'dark' ? <Sun className="h-4 w-4" /> : <Moon className="h-4 w-4" />}
        </Button>
      </div>

      <div className="w-full max-w-sm space-y-6 px-4">
        {/* Logo */}
        <div className="flex flex-col items-center gap-2">
          <img src={logo} alt="rune" className="h-12 w-auto" />
        </div>

        {phase === 'token' ? (
          /* ── Step 1: save token ── */
          <div className="rounded-lg border border-border bg-card p-5 space-y-4">
            <div className="space-y-1">
              <h2 className="text-sm font-semibold">Set your access token</h2>
              <p className="text-xs text-muted-foreground">
                Saved locally — you won't be asked again unless you log out.
              </p>
            </div>
            <div className="space-y-3">
              <div className="space-y-1.5">
                <Label htmlFor="token" className="text-xs text-muted-foreground">Token</Label>
                <div className="relative">
                  <Input
                    id="token"
                    type={showToken ? 'text' : 'password'}
                    placeholder="rune.xxxxxxxxxxxxxxxx"
                    value={tokenInput}
                    onChange={e => setTokenInput(e.target.value)}
                    onKeyDown={e => e.key === 'Enter' && handleSaveToken()}
                    className="pr-9 font-mono text-sm"
                    autoFocus
                  />
                  <button
                    type="button"
                    onClick={() => setShowToken(v => !v)}
                    className="absolute right-2.5 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
                  >
                    {showToken ? <EyeOff className="h-3.5 w-3.5" /> : <Eye className="h-3.5 w-3.5" />}
                  </button>
                </div>
              </div>
              <Button
                className="w-full"
                onClick={handleSaveToken}
                disabled={!tokenInput.trim() || saving}
              >
                {saving ? 'Validating…' : 'Continue'}
              </Button>
            </div>
          </div>
        ) : (
          /* ── Step 2: unseal vault ── */
          <div className="rounded-lg border border-border bg-card p-5 space-y-4">
            {sealed ? (
              <>
                <div className="space-y-1">
                  <h2 className="text-sm font-semibold">Unlock vault</h2>
                  <p className="text-xs text-muted-foreground">Enter your passphrase to continue</p>
                </div>
                <div className="space-y-3">
                  <div className="space-y-1.5">
                    <Label htmlFor="passphrase" className="text-xs text-muted-foreground">Passphrase</Label>
                    <Input
                      id="passphrase"
                      type="password"
                      placeholder="••••••••••••"
                      value={passphrase}
                      onChange={e => setPassphrase(e.target.value)}
                      onKeyDown={e => e.key === 'Enter' && handleUnseal()}
                      autoFocus
                    />
                  </div>
                  <Button
                    className="w-full"
                    onClick={handleUnseal}
                    disabled={!passphrase.trim() || unsealing}
                  >
                    {unsealing ? 'Unlocking…' : 'Unlock vault'}
                  </Button>
                </div>
              </>
            ) : (
              /* Vault already unsealed — App.tsx auto-navigates, this is a brief flash */
              <div className="flex items-center justify-center py-4">
                <div className="h-4 w-4 animate-spin rounded-full border-2 border-muted border-t-foreground" />
              </div>
            )}

            {/* Token indicator + change link */}
            <div className="flex items-center justify-between pt-1 border-t border-border">
              <span className="text-xs text-muted-foreground font-mono">{maskedToken}</span>
              <button
                type="button"
                onClick={onClearToken}
                className="text-xs text-muted-foreground hover:text-foreground underline underline-offset-2 transition-colors"
              >
                Change token
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
