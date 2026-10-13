// Deprecated: use ./model-provider instead. Kept for backward compat with UI.
export {
  ModelProviderService as LLMProviderService,
  type AddModelProviderPayload as AddLLMProviderPayload,
} from "./model-provider";
export type {
  AddModelProviderRequest as AddLLMProviderRequest,
  AddModelProviderResponse as AddLLMProviderResponse,
  DeleteModelProviderRequest as DeleteLLMProviderRequest,
  GetListModelProviderResponse as GetListLLMProviderResponse,
  GetModelProviderDetailsResponse as GetLLMProviderDetailsResponse,
  GetModelProviderModelsRequest as GetLLMProviderModelsRequest,
  GetModelProviderModelsResponse as GetLLMProviderModelsResponse,
  ModelProvider as LLMProvider,
} from "@nagare-app/dtos";
export {
  AddModelProviderEndpoint as AddLLMProviderEndpoint,
  DeleteModelProviderEndpoint as DeleteLLMProviderEndpoint,
  ListModelProvidersEndpoint as ListLLMProvidersEndpoint,
  GetModelProviderDetailsEndpoint as GetLLMProviderDetailsEndpoint,
  GetModelProviderModelsEndpoint as GetLLMProviderModelsEndpoint,
} from "@nagare-app/dtos";
