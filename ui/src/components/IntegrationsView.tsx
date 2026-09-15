import { ExternalLink, ArrowRight, FileText } from 'lucide-react'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import hashicorpIcon from '@/assets/hashicorp.png'
import awsIcon from '@/assets/aws.png'
import azureIcon from '@/assets/azure.png'
import dopplerIcon from '@/assets/doppler.png'
import googleIcon from '@/assets/google.png'

interface Integration {
  id: string
  name: string
  description: string
  icon: React.ReactNode
  docsUrl: string | null
}

function Img({ src, alt }: { src: string; alt: string }) {
  return <img src={src} alt={alt} className="size-5 object-contain" />
}

const INTEGRATIONS: Integration[] = [
  {
    id: 'hashicorp',
    name: 'HashiCorp Vault',
    description: 'Import KV secrets, namespaces, and ACL policy mappings.',
    icon: <Img src={hashicorpIcon} alt="HashiCorp" />,
    docsUrl: 'https://developer.hashicorp.com/vault/docs',
  },
  {
    id: 'aws',
    name: 'AWS Secrets Manager',
    description: 'Pull secrets from AWS Secrets Manager into a Rune namespace.',
    icon: <Img src={awsIcon} alt="AWS" />,
    docsUrl: 'https://docs.aws.amazon.com/secretsmanager/',
  },
  {
    id: 'gcp',
    name: 'GCP Secret Manager',
    description: 'Import secrets from Google Cloud Secret Manager.',
    icon: <Img src={googleIcon} alt="GCP" />,
    docsUrl: 'https://cloud.google.com/secret-manager/docs',
  },
  {
    id: 'azure',
    name: 'Azure Key Vault',
    description: 'Migrate secrets and keys from Azure Key Vault into Rune.',
    icon: <Img src={azureIcon} alt="Azure" />,
    docsUrl: 'https://learn.microsoft.com/en-us/azure/key-vault/',
  },
  {
    id: 'doppler',
    name: 'Doppler',
    description: 'Sync secrets from a Doppler project into Rune namespaces.',
    icon: <Img src={dopplerIcon} alt="Doppler" />,
    docsUrl: 'https://docs.doppler.com/',
  },
  {
    id: 'dotenv',
    name: '.env File',
    description: 'Import secrets from a .env file. Each KEY=VALUE becomes a secret.',
    icon: <FileText className="size-4 text-muted-foreground" />,
    docsUrl: null,
  },
]

export function IntegrationsView() {
  return (
    <div className="p-6 space-y-6">
      {/* Page header */}
      <div>
        <h1 className="text-xl font-semibold tracking-tight">Integrations</h1>
        <p className="text-sm text-muted-foreground mt-0.5">
          Connect Rune Vault to external secret sources and import pipelines.
        </p>
      </div>

      {/* Cards grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
        {INTEGRATIONS.map(integration => (
          <div
            key={integration.id}
            className="rounded-lg border border-border bg-card p-4 flex flex-col gap-3"
          >
            {/* Icon + name + badge */}
            <div className="flex items-center justify-between gap-2">
              <div className="flex items-center gap-2.5">
                <div className="flex size-7 shrink-0 items-center justify-center rounded-md border border-border bg-muted">
                  {integration.icon}
                </div>
                <span className="text-sm font-medium leading-tight">{integration.name}</span>
              </div>
              <Badge variant="secondary" className="shrink-0 text-[10px] px-1.5 py-0">
                Soon
              </Badge>
            </div>

            {/* Description */}
            <p className="text-xs text-muted-foreground leading-relaxed">
              {integration.description}
            </p>

            {/* Footer */}
            <div className="flex items-center justify-between mt-auto">
              {integration.docsUrl ? (
                <a
                  href={integration.docsUrl}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="inline-flex items-center gap-1 text-xs text-muted-foreground hover:text-foreground transition-colors"
                >
                  <ExternalLink className="size-3" />
                  Docs
                </a>
              ) : (
                <span />
              )}
              <Button size="sm" variant="outline" disabled className="gap-1 text-xs h-6 px-2">
                Configure
                <ArrowRight className="size-3" />
              </Button>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}
