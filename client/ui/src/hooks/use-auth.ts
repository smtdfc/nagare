import { create } from "zustand";

export type AuthState = {
  isAuthenticated: boolean;
  setIsAuthenticated: (isAuthenticated: boolean) => void;
};
export const useAuth = create<AuthState>((set) => ({
  isAuthenticated: false,
  setIsAuthenticated(isAuthenticated) {
    set({ isAuthenticated });
  },
}));
