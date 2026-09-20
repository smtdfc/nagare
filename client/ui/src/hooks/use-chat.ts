import { create } from 'zustand'


export type ChatState = {
    currentChatID: string | null
    chats:any[]

    setCurrentChatID: (newChatID: string) => void
}
export const useChat = create<ChatState>((set) => ({
    currentChatID: null,
    chats: [],
    setCurrentChatID: (newChatID: string) => {set({currentChatID: newChatID})},
}))