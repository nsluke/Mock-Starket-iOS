import { create } from 'zustand';
import type { User } from '@/types/user';

interface AuthState {
  user: User | null;
  token: string | null;
  isLoading: boolean;
  setUser: (user: User | null) => void;
  setToken: (token: string | null) => void;
  setLoading: (loading: boolean) => void;
  signOut: () => void;
}

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  token: null,
  isLoading: true,

  setUser: (user) => set({ user }),
  setToken: (token) => set({ token }),
  setLoading: (isLoading) => set({ isLoading }),

  signOut: () => {
    if (typeof window !== 'undefined') {
      localStorage.removeItem('mockstarket_token');
      document.cookie = 'mockstarket_token=; path=/; max-age=0';
    }
    // Fire-and-forget Firebase sign out
    import('@/lib/auth-service').then(({ signOut }) => signOut()).catch(() => {});
    set({ user: null, token: null });
  },
}));
