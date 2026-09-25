import { create } from "zustand";
import type { Session } from "@nagare-app/dtos";
import type { Message } from "@nagare-app/messages";

export type ChatState = {
  currentChatID: string | null;
  chatSession: Session | null;
  isConnected: boolean;
  isProcessing: boolean;
  chats: Session[];
  pendingMessages: Message[];

  setChatSession: (session: Session | null) => void;
  setIsConnected: (isConnected: boolean) => void;
  setIsProcessing: (isProcessing: boolean) => void;
  setPendingMessages: (pendingMessages: Message[]) => void;
  reset: () => void;
};

export const useChat = create<ChatState>((set) => ({
  currentChatID: null,
  chatSession: null,
  chats: [],
  isConnected: false,
  isProcessing: false,
  pendingMessages: [],
  setChatSession: (session: Session | null) => {
    if (!session) {
      set({ currentChatID: null, chatSession: null });
      return;
    }

    set({ currentChatID: session.id, chatSession: session });
  },

  setIsConnected: (isConnected: boolean) => {
    set({ isConnected });
  },

  setIsProcessing: (isProcessing: boolean) => {
    set({ isProcessing });
  },

  setPendingMessages: (messages: Message[]) => {
    set({ pendingMessages: messages });
  },

  reset: () =>
    set({
      currentChatID: null,
      chatSession: null,
      isConnected: false,
      isProcessing: false,
      pendingMessages: [],
    }),
}));
