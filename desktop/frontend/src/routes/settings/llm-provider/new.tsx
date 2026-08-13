import LLMProviderEditForm from '#/components/llm-provider-edit-form'
import { toast } from '#/components/ui/toast'
import type { Provider } from '#/dto/api'
import { getErrorMessage } from '#/lib/error'
import { createFileRoute } from '@tanstack/react-router'
import { useEffect, useState } from 'react'

export const Route = createFileRoute('/settings/llm-provider/new')({
  component: RouteComponent,
})

function RouteComponent() {
  const [currentProvider, setCurrentProvider] = useState<Provider>({
    id: '',
    name: '',
    compatible: '',
    base_url: '',
    api_key: '',
    is_enable: false,
    available_models: [],
  })

  useEffect(() => {
    const fetchProvider = async () => {
      try {
        // const details = await ProviderService.getProviderById(id)
        // setCurrentProvider(details)
      } catch (e) {
        toast.add({
          type: 'error',
          description: getErrorMessage(e),
          priority: 'high',
        })
      }
    }

    fetchProvider()
  }, [])

  return (
    <div className="container mx-auto py-8 px-6 max-w-6xl">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">New Provider</h1>
        </div>
      </div>
      {currentProvider ? (
        <LLMProviderEditForm provider={currentProvider} />
      ) : (
        'Error'
      )}
    </div>
  )
}
