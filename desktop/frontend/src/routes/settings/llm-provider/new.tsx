import LLMProviderEditForm from '#/components/llm-provider-edit-form'
import { toast } from '#/components/ui/toast'
import type { Provider } from '#/dto/api'
import { getErrorMessage } from '#/lib/error'
import { ProviderService } from '#/services/provider'
import { createFileRoute, useRouter } from '@tanstack/react-router'
import { useEffect, useState } from 'react'

export const Route = createFileRoute('/settings/llm-provider/new')({
  component: RouteComponent,
})

function RouteComponent() {
  const router = useRouter()
  const [currentProvider, setCurrentProvider] = useState<Provider>({
    id: '',
    name: '',
    compatible: '',
    base_url: '',
    api_key: '',
    is_enable: false,
    available_models: [],
  })
  const handleSubmit = async (provider: Provider) => {
    try {
      await ProviderService.createProvider({
        name: provider.name,
        compatible: provider.compatible,
        base_url: provider.base_url,
        api_key: provider.api_key,
        is_enable: provider.is_enable,
        available_models: provider.available_models,
      })

      toast.add({
        type: 'success',
        description: 'Provider created successfully',
        priority: 'high',
      })

      router.navigate({
        to: '/settings/llm-provider/overview',
      })
    } catch (e) {
      toast.add({
        type: 'error',
        description: getErrorMessage(e),
        priority: 'high',
      })
    }
  }

  return (
    <div className="container mx-auto py-8 px-6 max-w-6xl">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">New Provider</h1>
        </div>
      </div>
      {currentProvider ? (
        <LLMProviderEditForm
          currentProvider={currentProvider}
          onCancel={() => alert('1')}
          onSubmit={handleSubmit}
        />
      ) : (
        'Error'
      )}
    </div>
  )
}
