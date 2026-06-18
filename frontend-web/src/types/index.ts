export interface User {
    id: string;
    username: string;
    email: string;
    avatar_url?: string;
    bio?: string;
}

export interface Post {
    id: string;
    author_id: string;
    author_name?: string;
    content: string;
    media_url?: string;
    created_at: string;
    updated_at: string;
    likes_count?: number;
    is_liked?: boolean;
    comments_count?: number;
}

export interface WSMessage {
    type: 'private_msg' | 'group_msg' | 'status';
    chat_id: string;
    text: string;
    to_user_id?: string;
    sender_id?: string;
    created_at?: string;
}

export interface AuthResponse {
    token: string;
}
