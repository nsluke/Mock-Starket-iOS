'use client';

import { useEffect } from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { apiClient } from '@/lib/api-client';
import { wsClient } from '@/lib/websocket-client';
import { useAuthStore } from '@/stores/auth-store';
import { useMarketStore } from '@/stores/market-store';

const navItems = [
  { href: '/market', label: 'Market', icon: '📊' },
  { href: '/portfolio', label: 'Portfolio', icon: '💼' },
  { href: '/orders', label: 'Orders', icon: '📋' },
  { href: '/leaderboard', label: 'Leaderboard', icon: '🏆' },
];

const sidebarExtras = [
  { href: '/alerts', label: 'Alerts', icon: '🔔' },
  { href: '/challenges', label: 'Challenges', icon: '🎯' },
  { href: '/achievements', label: 'Achievements', icon: '⭐' },
];

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const pathname = usePathname();
  const { token } = useAuthStore();
  const { updatePrices } = useMarketStore();

  useEffect(() => {
    // Restore token from localStorage
    const savedToken = localStorage.getItem('mockstarket_token');
    if (savedToken) {
      apiClient.setToken(savedToken);
      useAuthStore.getState().setToken(savedToken);
    }
  }, []);

  useEffect(() => {
    if (!token) return;

    wsClient.connect(token);
    wsClient.on('price_batch', (data) => {
      updatePrices(data);
    });

    return () => wsClient.disconnect();
  }, [token, updatePrices]);

  return (
    <div className="flex h-screen bg-[#0D1117]">
      {/* Sidebar */}
      <aside className="hidden md:flex w-64 flex-col border-r border-[#30363D] bg-[#161B22]">
        <div className="p-6">
          <h1 className="text-xl font-bold text-white">Mock Starket</h1>
          <p className="text-xs text-[#6E7681] mt-1">Paper Trading Simulator</p>
        </div>

        <nav className="flex-1 px-3 space-y-1">
          {navItems.map((item) => (
            <Link
              key={item.href}
              href={item.href}
              className={`flex items-center gap-3 rounded-lg px-3 py-2.5 text-sm font-medium transition-colors ${
                pathname === item.href
                  ? 'bg-[#21262D] text-white'
                  : 'text-[#8B949E] hover:bg-[#21262D] hover:text-white'
              }`}
            >
              <span>{item.icon}</span>
              {item.label}
            </Link>
          ))}
        </nav>

        <div className="px-3 pt-4 border-t border-[#30363D] space-y-1">
          {sidebarExtras.map((item) => (
            <Link
              key={item.href}
              href={item.href}
              className={`flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors ${
                pathname === item.href
                  ? 'bg-[#21262D] text-white'
                  : 'text-[#8B949E] hover:bg-[#21262D] hover:text-white'
              }`}
            >
              <span>{item.icon}</span>
              {item.label}
            </Link>
          ))}
        </div>

        <div className="p-4 border-t border-[#30363D]">
          <Link
            href="/settings"
            className="flex items-center gap-3 rounded-lg px-3 py-2.5 text-sm text-[#8B949E] hover:text-white transition-colors"
          >
            ⚙️ Settings
          </Link>
        </div>
      </aside>

      {/* Main content */}
      <main className="flex-1 overflow-y-auto">
        {children}
      </main>

      {/* Mobile bottom nav */}
      <nav className="md:hidden fixed bottom-0 left-0 right-0 flex border-t border-[#30363D] bg-[#161B22]">
        {navItems.map((item) => (
          <Link
            key={item.href}
            href={item.href}
            className={`flex-1 flex flex-col items-center py-2 text-xs ${
              pathname === item.href ? 'text-[#50E3C2]' : 'text-[#6E7681]'
            }`}
          >
            <span className="text-lg">{item.icon}</span>
            {item.label}
          </Link>
        ))}
      </nav>
    </div>
  );
}
