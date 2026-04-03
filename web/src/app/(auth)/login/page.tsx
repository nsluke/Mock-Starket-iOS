'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { TrendingUp } from 'lucide-react';
import { apiClient } from '@/lib/api-client';
import { useAuthStore } from '@/stores/auth-store';

export default function LoginPage() {
  const router = useRouter();
  const { setToken, setUser } = useAuthStore();
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleGuestLogin = async () => {
    setIsLoading(true);
    setError(null);

    try {
      const guestUID = `guest_${crypto.randomUUID().slice(0, 8)}`;
      apiClient.setToken(guestUID);

      const user = await apiClient.createGuest();
      setToken(guestUID);
      setUser(user);
      localStorage.setItem('mockstarket_token', guestUID);
      document.cookie = `mockstarket_token=${guestUID}; path=/; max-age=${60 * 60 * 24 * 30}`;
      router.replace('/market');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Login failed');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-[#0D1117] px-4">
      <div className="w-full max-w-sm space-y-10">
        {/* Logo */}
        <div className="text-center space-y-4">
          <TrendingUp className="w-16 h-16 text-[#50E3C2] mx-auto" />
          <h1 className="text-4xl font-bold tracking-tight text-white">
            Mock Starket
          </h1>
          <p className="text-[#8B949E]">Learn to trade. Risk nothing.</p>
        </div>

        {/* Actions */}
        <div className="space-y-4">
          <button
            onClick={handleGuestLogin}
            disabled={isLoading}
            className="w-full rounded-xl bg-[#50E3C2] px-6 py-4 text-lg font-semibold text-black transition-transform hover:scale-[1.02] active:scale-[0.98] disabled:opacity-50"
          >
            {isLoading ? 'Loading...' : 'Continue as Guest'}
          </button>

          <button
            onClick={handleGuestLogin}
            className="w-full rounded-xl bg-[#21262D] px-6 py-4 text-lg font-semibold text-white transition-colors hover:bg-[#30363D]"
          >
            Sign in with Email
          </button>

          {error && (
            <p className="text-center text-sm text-red-400">{error}</p>
          )}
        </div>

        <p className="text-center text-xs text-[#6E7681]">
          Start with $100,000 in virtual cash
        </p>
      </div>
    </div>
  );
}
