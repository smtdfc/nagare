import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Field, FieldGroup, FieldLabel, FieldSet } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import type { Provider } from '#/dto/api.ts'
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from './ui/select'
import { ProviderService } from '#/services/provider'
import { toast } from './ui/toast'

type LLMProviderEditForm = {
  currentProvider: Provider
  onSubmit: (provider: Provider) => void
  onCancel: () => void
}
export default function LLMProviderEditForm({
  currentProvider,
  onSubmit,
  onCancel,
}: LLMProviderEditForm) {
  const [provider, setProvider] = useState<Provider>(currentProvider)
  const [models, setModels] = useState<string[]>(provider.available_models)
  const [newModelInput, setNewModelInput] = useState('')

  const handleAddModel = () => {
    if (newModelInput.trim() && !models.includes(newModelInput.trim())) {
      setModels([...models, newModelInput.trim()])
      setNewModelInput('')
    }
  }

  const handleRemoveModel = (indexToRemove: number) => {
    setModels(models.filter((_, index) => index !== indexToRemove))
  }

  const handleUpdateModel = (value: string, indexToUpdate: number) => {
    const updated = [...models]
    updated[indexToUpdate] = value
    setModels(updated)
  }

  const handleSubmit = () => {
    onSubmit({ ...provider, available_models: models })
  }

  const handleFetchModel = async () => {
    try {
      const models = await ProviderService.fetchModel({
        compatible: provider.compatible,
        base_url: provider.base_url,
        api_key: provider.api_key,
      })
      setProvider({ ...provider, available_models: models ?? [] })
      setModels(models ?? [])
      toast.add({
        title: 'Success',
        description: 'Models fetched successfully',
        type: 'success',
      })
    } catch (e) {
      toast.add({
        title: 'Error',
        description: 'Failed to fetch models',
        type: 'error',
      })
    }
  }

  const compatibles = [{ label: 'OpenAI Compatible', value: 'OpenAI' }]

  return (
    <div className="w-full max-w-md">
      <form onSubmit={(e) => e.preventDefault()}>
        <FieldGroup>
          <FieldSet>
            <FieldGroup>
              <Field>
                <FieldLabel htmlFor="provider-name">Provider name:</FieldLabel>
                <Input
                  id="provider-name"
                  placeholder="OpenAI"
                  required
                  defaultValue={provider.name}
                  onChange={(e) =>
                    setProvider({ ...provider, name: e.target.value })
                  }
                />
              </Field>
              <Field>
                <FieldLabel htmlFor="comp-base-url">Compatible:</FieldLabel>
                <Select
                  items={compatibles}
                  defaultValue={provider.compatible}
                  onValueChange={(value) =>
                    setProvider({ ...provider, compatible: value! })
                  }
                >
                  <SelectTrigger className="w-45">
                    <SelectValue placeholder="Compatible" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectGroup>
                      {compatibles.map((item) => (
                        <SelectItem key={item.value} value={item.value}>
                          {item.label}
                        </SelectItem>
                      ))}
                    </SelectGroup>
                  </SelectContent>
                </Select>
              </Field>
              <Field>
                <FieldLabel htmlFor="provider-base-url">BaseURL:</FieldLabel>
                <Input
                  id="provider-base-url"
                  placeholder="https://example.com/openai/v1"
                  required
                  defaultValue={provider.base_url}
                  onChange={(e) =>
                    setProvider({ ...provider, base_url: e.target.value })
                  }
                />
              </Field>
              <Field>
                <FieldLabel htmlFor="provider-api-key">API Key:</FieldLabel>
                <Textarea
                  id="provider-api-key"
                  placeholder="Your API Key"
                  className="resize-none"
                  rows={3}
                  defaultValue={provider.api_key}
                  onChange={(e) =>
                    setProvider({ ...provider, api_key: e.target.value })
                  }
                />
              </Field>

              <Field>
                <FieldLabel>Available models:</FieldLabel>
                <div className="space-y-2 mt-1">
                  {models.map((model, index) => (
                    <div key={index} className="flex items-center gap-2">
                      <Input
                        defaultValue={model}
                        onChange={(e) =>
                          handleUpdateModel(e.target.value, index)
                        }
                        placeholder="Model name..."
                      />
                      <Button
                        type="button"
                        variant="destructive"
                        size="sm"
                        onClick={() => handleRemoveModel(index)}
                      >
                        Remove
                      </Button>
                    </div>
                  ))}

                  <div className="flex items-center gap-2 pt-2">
                    <Input
                      defaultValue={newModelInput}
                      onChange={(e) => setNewModelInput(e.target.value)}
                      placeholder="Model name"
                      onKeyDown={(e) => {
                        if (e.key === 'Enter') {
                          e.preventDefault()
                          handleAddModel()
                        }
                      }}
                    />
                    <Button
                      type="button"
                      variant="secondary"
                      onClick={handleAddModel}
                    >
                      Add
                    </Button>
                    <Button
                      type="button"
                      variant="secondary"
                      onClick={handleFetchModel}
                    >
                      Fetch
                    </Button>
                  </div>
                </div>
              </Field>

              <div className="flex items-center justify-end gap-3 pt-4">
                <Button type="button" variant="outline" onClick={onCancel}>
                  Cancel
                </Button>
                <Button onClick={handleSubmit}>Save</Button>
              </div>
            </FieldGroup>
          </FieldSet>
        </FieldGroup>
      </form>
    </div>
  )
}
