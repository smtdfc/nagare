import { DataTable } from '#/components/data-table.tsx'
import { Button } from '#/components/ui/button.tsx'
import { Spinner } from '#/components/ui/spinner.tsx'
import type { Provider } from '#/dto/api.ts'
import { ProviderService } from '#/services/provider.ts'
import { createFileRoute, useNavigate, useRouter } from '@tanstack/react-router'
import { useEffect, useMemo, useState } from 'react'
import { Plus } from 'lucide-react'
import { createColumns } from '#/components/columns/provider-columns.tsx'

export const Route = createFileRoute('/settings/llm-provider/overview')({
  component: RouteComponent,
  staticData: {
    breadcrumb: 'LLM Provider settings',
  },
})

function RouteComponent() {
  const navigate = useNavigate()
  const [isLoading, setIsLoading] = useState<boolean>(true)
  const [providers, setProviders] = useState<Provider[]>([])
  const router = useRouter()
  useEffect(() => {
    ;(async () => {
      try {
        const data = await ProviderService.getListProvider()
        setProviders(data.providers || [])
      } catch (err) {
        console.error('Failed to load providers:', err)
      } finally {
        setIsLoading(false)
      }
    })()
  }, [])

  const columns = useMemo(() => {
    return createColumns(
      (provider) => {
        router.navigate({
          to: '/settings/llm-provider/edit/$id',
          params: { id: provider.id },
        })
      },
      (provider) => {},
    )
  }, [])

  return (
    <div className="container mx-auto py-8 px-6 max-w-6xl">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">LLM Providers</h1>
          <p className="text-sm text-muted-foreground mt-1">
            Manage your AI provider endpoints and configurations.
          </p>
        </div>
        <Button
          className="flex items-center gap-2 "
          onClick={() => navigate('/settings/llm-provider/new')}
        >
          <Plus className="size-4" />
          Add Provider
        </Button>
      </div>

      {isLoading ? (
        <div className="flex items-center justify-center h-[50vh] w-full rounded-lg border border-dashed">
          <Spinner className="size-8 text-primary animate-spin" />
        </div>
      ) : (
        <div className=" overflow-hidden">
          <DataTable columns={columns} data={providers} />
        </div>
      )}
    </div>
  )
}
