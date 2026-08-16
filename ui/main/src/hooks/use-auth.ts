import type { Profile } from "@nagare-agent/dto";
import { create } from "zustand";

type AuthStore = {
  auth: Profile | null;
  setAuthState: (auth: Profile) => void;
};

export const useAuth = create<AuthStore>((set) => ({
  auth: null,
  setAuthState: (auth: Profile) => set({ auth }),
}));
