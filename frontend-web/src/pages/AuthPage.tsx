import React, { useState } from 'react';
import api from '../api/axios';
import { useNavigate } from 'react-router-dom';
import type {AuthResponse} from '../types';

export default function AuthPage() {
    const [isLogin, setIsLogin] = useState<boolean>(true);
    const [form, setForm] = useState({ username: '', email: '', password: '' });
    const navigate = useNavigate();

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            if (isLogin) {
                const res = await api.post<AuthResponse>('/auth/login', {
                    username: form.username,
                    password: form.password
                });
                localStorage.setItem('token', res.data.token);
                navigate('/feed');
            } else {
                await api.post('/auth/register', form);
                alert('Registration successful! Please login.');
                setIsLogin(true);
            }
        } catch (err: any) {
            alert(err.response?.data?.error || 'Something went wrong');
        }
    };

    return (
        <div className="min-h-screen flex items-center justify-center bg-slate-100">
            <div className="bg-white p-8 rounded-2xl shadow-xl w-full max-w-md border border-slate-200">
                <h2 className="text-2xl font-bold mb-6 text-center text-slate-800">
                    {isLogin ? 'Welcome Back' : 'Create Account'}
                </h2>
                <form onSubmit={handleSubmit} className="space-y-4">
                    {!isLogin && (
                        <input
                            className="w-full p-3 border rounded-xl outline-none focus:ring-2 focus:ring-blue-500"
                            placeholder="Email"
                            required
                            onChange={e => setForm({...form, email: e.target.value})}
                        />
                    )}
                    <input
                        className="w-full p-3 border rounded-xl outline-none focus:ring-2 focus:ring-blue-500"
                        placeholder="Username"
                        required
                        onChange={e => setForm({...form, username: e.target.value})}
                    />
                    <input
                        className="w-full p-3 border rounded-xl outline-none focus:ring-2 focus:ring-blue-500"
                        type="password"
                        placeholder="Password"
                        required
                        onChange={e => setForm({...form, password: e.target.value})}
                    />
                    <button className="w-full bg-blue-600 text-white p-3 rounded-xl font-bold hover:bg-blue-700 transition-colors">
                        {isLogin ? 'Login' : 'Register'}
                    </button>
                </form>
                <p
                    className="mt-4 text-center text-sm text-blue-600 cursor-pointer hover:underline"
                    onClick={() => setIsLogin(!isLogin)}
                >
                    {isLogin ? "Don't have an account? Register" : "Already have one? Login"}
                </p>
            </div>
        </div>
    );
}
