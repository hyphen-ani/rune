import { useEffect, useState } from 'react'
import { KeyRound, Layers, Coins, RotateCcw, ShieldCheck, RefreshCw } from 'lucide-react'
import {
  Bar, BarChart, XAxis, YAxis, CartesianGrid, LabelList,
  Pie, PieChart,
} from 'recharts'
import {
  ChartContainer, ChartTooltip, ChartTooltipContent,
  type ChartConfig,
} from '@/components/ui/chart'
import {
  Card, CardContent, CardDescription, CardHeader, CardTitle,
} from '@/components/ui/card'
import { api, type VaultStats } from '@/lib/api'
import { toast } from 'sonner'

interface Props {
  token: string
}

const secretsConfig: ChartConfig = {
  count: { label: 'Secrets', color: 'var(--chart-1)' },
  label: { color: 'var(--background)' },
}


const pieConfig: ChartConfig = {
  value:   { label: 'Tokens' },
  active:  { label: 'Active',  color: 'var(--chart-2)' },
  revoked: { label: 'Revoked', color: 'var(--chart-5)' },
}

function StatCard({
  label, value, icon: Icon,
}: { label: string; value: number; icon: React.ElementType }) {
  return (
    <div className="rounded-xl border border-border bg-card p-4 flex flex-col gap-2">
      <div className="flex items-center justify-between">
        <span className="text-xs text-muted-foreground">{label}</span>
        <div className="rounded-md bg-muted p-1.5">
          <Icon className="size-3.5 text-muted-foreground" />
        </div>
      </div>
      <span className="text-3xl font-bold tracking-tight">{value}</span>
    </div>
  )
}

function EmptyChart({ message }: { message: string }) {
  return (
    <div className="flex h-full items-center justify-center">
      <p className="text-xs text-muted-foreground">{message}</p>
    </div>
  )
}


export function DashboardView({ token }: Props) {
  const [stats, setStats] = useState<VaultStats | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    api.getStats(token)
      .then(setStats)
      .catch(e => toast.error('Failed to load stats', { description: String(e) }))
      .finally(() => setLoading(false))
  }, [token])

  if (loading) {
    return (
      <div className="flex h-full items-center justify-center">
        <RefreshCw className="size-4 animate-spin text-muted-foreground" />
      </div>
    )
  }

  if (!stats) return null

  const tokenStatusData = [
    { status: 'active',  value: stats.active_tokens,  fill: 'var(--color-active)' },
    { status: 'revoked', value: stats.revoked_tokens, fill: 'var(--color-revoked)' },
  ].filter(d => d.value > 0)

  return (
    <div className="p-6 space-y-6">
      {/* Stat cards */}
      <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-5 gap-3">
        <StatCard label="Total Secrets"   value={stats.total_secrets}    icon={KeyRound}    />
        <StatCard label="Namespaces"      value={stats.total_namespaces} icon={Layers}      />
        <StatCard label="Active Tokens"   value={stats.active_tokens}    icon={ShieldCheck} />
        <StatCard label="Revoked Tokens"  value={stats.revoked_tokens}   icon={Coins}       />
        <StatCard label="Rotated Secrets" value={stats.rotated_secrets}  icon={RotateCcw}   />
      </div>

      {/* Row 2: bar + donut */}
      <div className="grid grid-cols-1 lg:grid-cols-5 gap-4">
        <Card className="lg:col-span-3">
          <CardHeader>
            <CardTitle>Secrets by namespace</CardTitle>
            <CardDescription>Distribution across namespaces</CardDescription>
          </CardHeader>
          <CardContent>
            <ChartContainer config={secretsConfig} className="h-52 w-full">
              {(stats.secrets_by_namespace?.length ?? 0) > 0 ? (
                <BarChart
                  data={stats.secrets_by_namespace}
                  layout="vertical"
                  margin={{ left: 0, right: 48 }}
                >
                  <CartesianGrid horizontal={false} />
                  <YAxis dataKey="namespace" type="category" hide />
                  <XAxis dataKey="count" type="number" hide />
                  <ChartTooltip cursor={false} content={<ChartTooltipContent indicator="line" />} />
                  <Bar dataKey="count" fill="var(--color-count)" radius={4}>
                    <LabelList
                      dataKey="namespace"
                      position="insideLeft"
                      offset={8}
                      className="fill-(--color-label)"
                      fontSize={12}
                    />
                    <LabelList
                      dataKey="count"
                      position="right"
                      offset={8}
                      className="fill-foreground"
                      fontSize={12}
                    />
                  </Bar>
                </BarChart>
              ) : (
                <EmptyChart message="No secrets yet" />
              )}
            </ChartContainer>
          </CardContent>
        </Card>

        <Card className="lg:col-span-2">
          <CardHeader>
            <CardTitle>Token status</CardTitle>
            <CardDescription>Active vs revoked</CardDescription>
          </CardHeader>
          <CardContent>
            <ChartContainer config={pieConfig} className="h-52 w-full">
              {tokenStatusData.length > 0 ? (
                <PieChart>
                  <ChartTooltip content={<ChartTooltipContent nameKey="status" hideLabel />} />
                  <Pie
                    data={tokenStatusData}
                    dataKey="value"
                    nameKey="status"
                    labelLine={false}
                    label={({ payload, cx, cy, x, y, textAnchor, dominantBaseline }) => (
                      <text
                        cx={cx}
                        cy={cy}
                        x={x}
                        y={y}
                        textAnchor={textAnchor}
                        dominantBaseline={dominantBaseline}
                        fill="var(--foreground)"
                        fontSize={13}
                        fontWeight={600}
                      >
                        {payload.value}
                      </text>
                    )}
                  />
                </PieChart>
              ) : (
                <EmptyChart message="No tokens yet" />
              )}
            </ChartContainer>
          </CardContent>
        </Card>
      </div>

    </div>
  )
}
