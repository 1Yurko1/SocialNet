import React, { useEffect, useState, useRef, useCallback } from 'react';
import api from '../api/axios';
import { type Post } from '../types';
import {Heart, MessageCircle, Send, Image as ImageIcon, MessageSquare} from 'lucide-react';
import {useNavigate} from "react-router-dom";

export default function FeedPage() {
    const navigate = useNavigate();
    const [posts, setPosts] = useState<Post[]>([]);
    const [content, setContent] = useState<string>('');
    const [loading, setLoading] = useState<boolean>(true);
    const [offset, setOffset] = useState<number>(0);
    const [hasMore, setHasMore] = useState<boolean>(true);

    const [, setCommentingPostId] = useState<string | null>(null);
    const [commentText, setCommentText] = useState<string>('');

    const [comments, setComments] = useState<any[]>([]);
    const [showComments, setShowComments] = useState<string | null>(null);

    // Реф для отслеживания последнего элемента (бесконечный скролл)
    const observer = useRef<IntersectionObserver | null>(null);

    // Функция загрузки постов с поддержкой смещения (offset)
    const fetchPosts = useCallback(async (currentOffset: number) => {
        try {
            const res = await api.get<Post[]>(`/feed?limit=10&offset=${currentOffset}`);
            const newPosts = res.data || [];
            if (newPosts.length < 10) setHasMore(false);
            setPosts(prev => (currentOffset === 0 ? newPosts : [...prev, ...newPosts]));
        } catch (err) {
            console.error("Error fetching feed:", err);
        } finally {
            setLoading(false);
        }
    }, []);

    // Первичная загрузка
    useEffect(() => {
        // eslint-disable-next-line react-hooks/set-state-in-effect
        fetchPosts(0);
    }, [fetchPosts]);

    const loadComments = async (postId: string) => {
        try {
            const res = await api.get(`/posts/${postId}/comments`);
            setComments(res.data || []);
            setShowComments(postId);
        } catch (err) {
            console.error("Error loading comments", err);
            // Даже при ошибке открываем блок, чтобы показать "No comments yet"
            setComments([]);
            setShowComments(postId);
        }
    };

    // Функция лайка
    const toggleLike = async (postId: string, currentlyLiked: boolean) => {
        try {
            if (currentlyLiked) {
                await api.delete(`/posts/${postId}/like`);
            } else {
                await api.post(`/posts/${postId}/like`);
            }

            // Оптимистичное обновление интерфейса
            setPosts(prev => prev.map(p =>
                p.id === postId
                    ? {
                        ...p,
                        is_liked: !currentlyLiked,
                        likes_count: (p.likes_count || 0) + (currentlyLiked ? -1 : 1)
                    }
                    : p
            ));
            // eslint-disable-next-line @typescript-eslint/no-unused-vars
        } catch (err) {
            alert("Could not update like");
        }
    };

    // Функция комментариев
    const handleAddComment = async (e: React.FormEvent, postId: string) => {
        e.preventDefault();
        if (!commentText.trim()) return;
        try {
            await api.post(`/posts/${postId}/comment`, { content: commentText });
            setCommentText('');
            setCommentingPostId(null);
            // Можно добавить обновление счетчика комментариев здесь
            // eslint-disable-next-line @typescript-eslint/no-unused-vars
        } catch (err) {
            alert("Could not add comment");
        }
    };

    // Функция для создания поста
    const createPost = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!content.trim()) return;
        try {
            await api.post('/posts', { content });
            setContent('');
            // Оптимистичное обновление: добавляем пост в начало списка сразу
            const tempPost: Post = {
                id: 'temp-' + Date.now(),
                author_id: 'me',
                author_name: 'Me',
                content: content,
                created_at: new Date().toISOString(),
                updated_at: new Date().toISOString(),
            };
            setPosts(prev => [tempPost, ...prev]);
            // Обновляем данные с сервера для получения реального ID
            fetchPosts(0);
            // eslint-disable-next-line @typescript-eslint/no-unused-vars
        } catch (err) {
            alert("Error creating post");
        }
    };

    // Callback для IntersectionObserver (триггер подгрузки)
    const lastPostRef = useCallback((node: HTMLDivElement | null) => {
        if (loading) return;
        if (observer.current) observer.current.disconnect();

        observer.current = new IntersectionObserver(entries => {
            if (entries[0].isIntersecting && hasMore) {
                const nextOffset = offset + 10;
                setOffset(nextOffset);
                fetchPosts(nextOffset);
            }
        });
        if (node) observer.current.observe(node);
    }, [loading, hasMore, fetchPosts, offset]);

    return (
        <div className="min-h-screen bg-slate-50 font-sans text-slate-900">
            <nav className="sticky top-0 z-50 bg-white/80 backdrop-blur-md border-b border-slate-200 p-4">
                <div className="max-w-2xl mx-auto flex justify-between items-center">
                    <h1 className="text-2xl font-black text-blue-600 tracking-tight cursor-pointer" onClick={() => navigate('/feed')}>
                        SocialNet
                    </h1>
                    <div className="flex items-center gap-4">
                        {/* КНОПКА ПЕРЕХОДА В ЧАТЫ */}
                        <button
                            onClick={() => navigate('/chat')}
                            className="p-2 text-slate-500 hover:text-blue-600 hover:bg-blue-50 rounded-full transition-all"
                            title="Messages"
                        >
                            <MessageSquare size={22} />
                        </button>

                        <button
                            onClick={() => { localStorage.removeItem('token'); window.location.reload(); }}
                            className="text-sm font-medium text-slate-500 hover:text-red-500 transition-colors"
                        >
                            Logout
                        </button>
                    </div>
                </div>
            </nav>

            <main className="max-w-2xl mx-auto py-8 px-4">
                <div className="mb-10 p-6 bg-white rounded-3xl shadow-sm border border-slate-200">
                    <form onSubmit={createPost} className="space-y-4">
                        <div className="flex gap-4">
                            <div className="w-12 h-12 bg-gradient-to-br from-blue-400 to-blue-600 rounded-full flex-shrink-0" />
                            <textarea
                                className="flex-1 p-3 bg-slate-50 rounded-2xl outline-none focus:ring-2 focus:ring-blue-500 transition-all resize-none"
                                placeholder="What's happening?"
                                rows={3}
                                value={content}
                                onChange={e => setContent(e.target.value)}
                            />
                        </div>
                        <div className="flex justify-between items-center pt-2">
                            <button type="button" className="p-2 text-slate-400 hover:text-blue-500 transition-colors">
                                <ImageIcon size={22} />
                            </button>
                            <button className="bg-blue-600 text-white px-6 py-2 rounded-full font-bold hover:bg-blue-700 transition-all flex items-center gap-2 shadow-lg shadow-blue-200">
                                <Send size={18} /> Post
                            </button>
                        </div>
                    </form>
                </div>

                <div className="space-y-6">
                    {loading && posts.length === 0 ? (
                        <div className="space-y-4">
                            {[1, 2, 3].map(i => (
                                <div key={i} className="p-6 bg-white rounded-3xl shadow-sm border border-slate-200 animate-pulse">
                                    <div className="flex items-center gap-3 mb-4">
                                        <div className="w-10 h-10 bg-slate-200 rounded-full" />
                                        <div className="h-4 w-32 bg-slate-200 rounded" />
                                    </div>
                                    <div className="h-4 w-full bg-slate-200 rounded mb-2" />
                                </div>
                            ))}
                        </div>
                    ) : posts.length === 0 ? (
                        <div className="text-center py-20 text-slate-400">No posts found.</div>
                    ) : (
                        <>
                            {posts.map((post, index) => (
                                <div
                                    key={post.id}
                                    ref={index === posts.length - 1 ? lastPostRef : null}
                                    className="p-6 bg-white rounded-3xl shadow-sm border border-slate-200 hover:border-blue-300 transition-all"
                                >
                                    <div className="flex items-center justify-between mb-4">
                                        <div className="flex items-center gap-3">
                                            <div className="w-10 h-10 bg-slate-200 rounded-full" />
                                            <div>
                                                <p className="font-bold text-slate-800 text-sm">{post.author_name || 'User'}</p>
                                                <p className="text-xs text-slate-400">Recently</p>
                                            </div>
                                        </div>
                                    </div>

                                    <p className="text-slate-700 leading-relaxed whitespace-pre-wrap mb-4">{post.content}</p>

                                    <div className="flex items-center gap-6 pt-4 border-t border-slate-50">
                                        <button
                                            onClick={() => toggleLike(post.id, post.is_liked || false)}
                                            className={`flex items-center gap-2 transition-colors group ${post.is_liked ? 'text-red-500' : 'text-slate-400 hover:text-red-500'}`}
                                        >
                                            <Heart size={20} className={post.is_liked ? 'fill-red-500' : 'group-hover:fill-red-500 transition-all'} />
                                            <span className="text-xs font-medium">{post.likes_count || 0}</span>
                                        </button>

                                        <button
                                            onClick={() => {
                                                if (showComments === post.id) {
                                                    setShowComments(null);
                                                } else {
                                                    // Сначала грузим данные, потом открываем блок
                                                    loadComments(post.id);
                                                    setShowComments(post.id);
                                                }
                                            }}
                                            className="flex items-center gap-2 text-slate-400 hover:text-blue-500 transition-colors group"
                                        >
                                            <MessageCircle size={20} className="group-hover:fill-blue-500 transition-all" />
                                            <span className="text-xs font-medium">Comments</span>
                                        </button>
                                    </div>

                                    {showComments === post.id && (
                                        <div className="mt-4 p-4 bg-slate-50 rounded-2xl border border-slate-100 space-y-4">
                                            {/* Заголовок секции */}
                                            <p className="text-xs font-bold text-slate-400 uppercase tracking-wider">
                                                Comments
                                            </p>

                                            {/* Список комментариев */}
                                            <div className="space-y-3">
                                                {(!comments || comments.length === 0) ? (
                                                    <p className="text-xs text-slate-400 italic">No comments yet. Be the first!</p>
                                                ) : (
                                                    comments?.map(com => ( // Используем ?. (optional chaining)
                                                        <div key={com.id} className="flex gap-2 text-sm leading-tight">
                                                            <span className="font-bold text-slate-700 whitespace-nowrap">
                                                                {com.author_name}:
                                                            </span>
                                                            <span className="text-slate-600">{com.content}</span>
                                                        </div>
                                                    ))
                                                )}
                                            </div>

                                            {/* Поле ввода комментария (теперь внутри того же блока) */}
                                            <form
                                                onSubmit={(e) => handleAddComment(e, post.id)}
                                                className="flex gap-2 pt-2 border-t border-slate-200"
                                            >
                                                <input
                                                    autoFocus
                                                    className="flex-1 p-2 bg-white border border-slate-200 rounded-xl text-sm outline-none focus:ring-2 focus:ring-blue-500 transition-all"
                                                    placeholder="Write a comment..."
                                                    value={commentText}
                                                    onChange={e => setCommentText(e.target.value)}
                                                />
                                                <button className="bg-blue-600 text-white px-4 py-2 rounded-xl text-sm font-bold hover:bg-blue-700 transition-colors">
                                                    Send
                                                </button>
                                            </form>
                                        </div>
                                    )}

                                </div>
                            ))}
                        </>
                    )}
                </div>
            </main>
        </div>
    );
}
