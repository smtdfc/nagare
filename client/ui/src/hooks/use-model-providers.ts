import { create } from "zustand";
import { ModelProviderService } from "@nagare-app/services";
import type {
  AddModelProviderRequest,
  ModelProvider,
} from "@nagare-app/dtos";

export type ModelProvidersState = {
  providers: ModelProvider[];
  isLoading: boolean;
  isSaving: boolean;

  setProviders: (providers: ModelProvider[]) => void;
  fetchProviders: () => Promise<void>;
  addProvider: (form: AddModelProviderRequest) => Promise<ModelProvider>;
  removeProvider: (id: string) => Promise<void>;
  reset: () => void;
};

export const useModelProviders = create<ModelProvidersState>((set) => ({
  providers: [],
  isLoading: false,
  isSaving: false,

  setProviders: (providers) => {
    set({ providers });
  },

  fetchProviders: async () => {
    set({ isLoading: true });
    try {
      const list = await ModelProviderService.list();
      set({ providers: list ?? [] });
    } finally {
      set({ isLoading: false });
    }
  },

  addProvider: async (form) => {
    set({ isSaving: true });
    try {
      const created = await ModelProviderService.add(form);
      set((state) => ({ providers: [...state.providers, created] }));
      return created;
    } finally {
      set({ isSaving: false });
    }
  },

  removeProvider: async (id) => {
    await ModelProviderService.delete(id);
    set((state) => ({
      providers: state.providers.filter((p) => p.id !== id),
    }));
  },

  reset: () =>
    set({
      providers: [],
      isLoading: false,
      isSaving: false,
    }),
}));
