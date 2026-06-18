import React, {useEffect, useRef, useState} from 'react';
import { useChatStore } from '../store/chatStore';
import api from '../api/axios';
import {ArrowLeft, MessageCircle, Search, Send, User as UserIcon} from 'lucide-react';
import type {WSMessage, User} from '../types';
import {useNavigate} from "react-router-dom";

export default function ChatPage() {
    const navigate = useNavigate()
    const { messages, activeChatId, activeUserId, setActiveChat, addMessage } = useChatStore();
    const [searchQuery, setSearchQuery] = useState('');
    const [users, setUsers] = useState<User[]>([]);
    const [messageText, setMessageText] = useState('');
    const [socket, setSocket] = useState<WebSocket | null>(null);
    const scrollRef = useRef<HTMLDivElement>(null);

    // 1. Подключение к WebSocket
    useEffect(() => {
        const token = localStorage.getItem('token');
        if (!token) {
            console.error("No token found");
            return;
        }

        // Передаем токен как параметр запроса
        const ws = new WebSocket(`ws://localhost:8080/api/v1/ws?token=${token}`);

        ws.onopen = () => {
            console.log("WebSocket Connected!");
        };

        ws.onmessage = (event) => {
            const data: WSMessage = JSON.parse(event.data);
            addMessage(data.chat_id, data);
        };

        ws.onerror = (error) => {
            console.error("WebSocket Error:", error);
        };

        ws.onclose = () => {
            console.log("WebSocket Disconnected");
        };

        // eslint-disable-next-line react-hooks/set-state-in-effect
        setSocket(ws);
        return () => ws.close();
    }, [addMessage]);

    // Авто-скролл вниз при новых сообщениях
    useEffect(() => {
        if (scrollRef.current) {
            scrollRef.current.scrollTop = scrollRef.current.scrollHeight;
        }
    }, [messages, activeChatId]);

    // 2. Поиск пользователей
    const searchUsers = async () => {
        if (!searchQuery.trim()) return;
        try {
            const res = await api.get(`/users/search?q=${searchQuery}`);
            setUsers(res.data);
        } catch (err) {
            console.error("Search failed", err);
        }
    };

    // 3. Выбор пользователя и инициализация чата
    const handleSelectUser = async (user: User) => {
        try {
            const res = await api.post('/chats/private', { target_user_id: user.id });
            const chatId = res.data.chat_id;

            setActiveChat(chatId, user.id);

            const historyRes = await api.get(`/chats/${chatId}/messages`);
            const history = historyRes.data;

            if (Array.isArray(history)) {
                history.forEach((msg: WSMessage) => addMessage(chatId, msg));
            } else {
                console.warn("Received invalid history format:", history);
            }

        } catch (err) {
            console.error("Chat init failed", err);
            alert("Could not start chat");
        }
    };

    // 4. Отправка сообщения
    const sendMessage = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!messageText.trim() || !activeChatId || !activeUserId || !socket) {
            console.error("Missing required data for sending message");
            return;
        }

        console.log("Sending to UserID:", activeChatId);

        const payload: WSMessage = {
            type: 'private_msg',
            chat_id: activeChatId,
            text: messageText,
            sender_id: 'me',
            to_user_id: activeUserId,
            created_at: new Date().toISOString(),
        };

        socket.send(JSON.stringify(payload));

        // Оптимистично добавляем своё сообщение
        addMessage(activeChatId, {
            ...payload,
            sender_id: 'me' // чтобы стилизовать как своё
        });
        setMessageText('');
    };

    return (
        <div className="flex h-screen bg-white overflow-hidden font-sans">
            {/* ЛЕВАЯ ПАНЕЛЬ */}
            <div className="w-80 border-r border-slate-200 flex flex-col bg-slate-50">
                <div className="p-4 border-b bg-white">
                    <div className="flex items-center gap-2 mb-4">
                        <button
                            onClick={() => navigate('/')}
                            className="p-1 hover:bg-slate-100 rounded-full text-slate-500 transition-colors"
                            title="Back to Feed"
                        >
                            <ArrowLeft size={20} />
                        </button>
                        <h2 className="text-xl font-bold text-slate-800">Messages</h2>
                    </div>
                    <div className="relative">
                        <input
                            className="w-full p-2 pl-10 bg-slate-100 rounded-xl outline-none text-sm focus:ring-2 focus:ring-blue-500 transition-all"
                            placeholder="Search users..."
                            value={searchQuery}
                            onChange={e => setSearchQuery(e.target.value)}
                            onKeyDown={e => e.key === 'Enter' && searchUsers()}
                        />
                        <Search className="absolute left-3 top-2.5 text-slate-400" size={16} />
                    </div>
                </div>

                <div className="flex-1 overflow-y-auto p-2 space-y-2">
                    {users.map(user => (
                        <div
                            key={user.id}
                            onClick={() => handleSelectUser(user)}
                            className={`p-3 rounded-2xl cursor-pointer transition-all flex items-center gap-3 ${activeChatId === user.id ? 'bg-blue-100 text-blue-600' : 'hover:bg-slate-200 text-slate-700'}`}
                        >
                            <div className="w-10 h-10 bg-slate-300 rounded-full flex items-center justify-center text-slate-500 overflow-hidden">
                                {user.avatar_url ? <img src={user.avatar_url} alt="" /> : <UserIcon size={20} />}
                            </div>
                            <span className="font-medium">{user.username}</span>
                        </div>
                    ))}
                </div>
            </div>

            {/* ПРАВАЯ ПАНЕЛЬ */}
            <div className="flex-1 flex flex-col bg-white">
                {activeChatId ? (
                    <>
                        <div className="p-4 border-b flex items-center gap-3 bg-white">
                            <div className="w-10 h-10 bg-blue-500 rounded-full" />
                            <span className="font-bold text-slate-800">Chat Room</span>
                        </div>

                        <div
                            ref={scrollRef}
                            className="flex-1 overflow-y-auto p-4 space-y-4 bg-slate-50"
                        >
                            {(messages[activeChatId] || []).map((msg, i) => (
                                <div key={i} className={`flex ${msg.sender_id === 'me' ? 'justify-end' : 'justify-start'}`}>
                                    <div className={`max-w-xs p-3 rounded-2xl text-sm shadow-sm ${
                                        msg.sender_id === 'me'
                                            ? 'bg-blue-600 text-white rounded-tr-none'
                                            : 'bg-white border border-slate-200 text-slate-800 rounded-tl-none'
                                    }`}>
                                        <p>{msg.text}</p>
                                        <p className={`text-[10px] mt-1 text-right ${msg.sender_id === 'me' ? 'text-blue-200' : 'text-slate-400'}`}>
                                            {msg.created_at ? new Date(msg.created_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }) : 'just now'}
                                        </p>
                                    </div>
                                </div>
                            ))}
                        </div>

                        <form onSubmit={sendMessage} className="p-4 border-t flex gap-2 bg-white">
                            <input
                                className="flex-1 p-3 bg-slate-100 rounded-2xl outline-none focus:ring-2 focus:ring-blue-500 transition-all"
                                placeholder="Type a message..."
                                value={messageText}
                                onChange={e => setMessageText(e.target.value)}
                            />
                            <button className="bg-blue-600 text-white p-3 rounded-2xl hover:bg-blue-700 transition-all shadow-lg shadow-blue-200">
                                <Send size={20} />
                            </button>
                        </form>
                    </>
                ) : (
                    <div className="flex-1 flex items-center justify-center text-slate-400 flex-col gap-2">
                        <MessageCircle size={64} className="opacity-10" />
                        <p className="text-lg font-medium">Select a user to start chatting</p>
                    </div>
                )}
            </div>
        </div>
    );
}
