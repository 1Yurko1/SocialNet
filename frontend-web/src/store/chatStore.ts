import { create } from 'zustand';
import type {WSMessage} from '../types';

interface ChatState {
    messages: Record<string, WSMessage[]>;
    activeChatId: string | null;
    activeUserId: string | null;
    setActiveChat: (chatId: string, userId: string) => void;
    addMessage: (chatId: string, message: WSMessage) => void;
}

export const useChatStore = create<ChatState>((set) => ({
    messages: {},
    activeChatId: null,
    activeUserId: null,
    setActiveChat: (chatId, userId) => set({ activeChatId: chatId, activeUserId: userId }),
    addMessage: (chatId, message) => set((state) => ({
        messages: {
            ...state.messages,
            [chatId]: [...(state.messages[chatId] || []), message],
        },
    })),
}));
