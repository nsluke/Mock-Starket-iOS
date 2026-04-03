'use client';

import { useEffect } from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { useQueryClient } from '@tanstack/react-query';
import { BarChart3, Briefcase, ClipboardList, Trophy, Star, Bell, Target, Award, Settings } from 'lucide-react';
import { apiClient } from '@/lib/api-client';
import { wsClient } from '@/lib/websocket-client';
import { toast } from '@/lib/toast';
import { useAuthStore } from '@/stores/auth-store';
import { useMarketStore } from '@/stores/market-store';
import { ConnectionStatus } from '@/components/ui/ConnectionStatus';
import { formatCurrency } from '@/lib/formatters';

const navItems = [
  { href: '/market', label: 'Market', icon: BarChart3 },
  { href: '/portfolio', label: 'Portfolio', icon: Briefcase },
  { href: '/orders', label: 'Orders', icon: ClipboardList },
  { href: '/leaderboard', label: 'Leaderboard', icon: Trophy },
];

const sidebarExtras = [
  { href: '/watchlist', label: 'Watchlist', icon: Star },
  { href: '/alerts', label: 'Alerts', icon: Bell },
  { href: '/challenges', label: 'Challenges', icon: Target },
  { href: '/achievements', label: 'Achievements', icon: Award },
];

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const pathname = usePathname();
  const { token } = useAuthStore();
  const { updatePrices } = useMarketStore();
  const queryClient = useQueryClient();

  useEffect(() => {
    // Restore token from localStorage
    const savedToken = localStorage.getItem('mockstarket_token');
    if (savedToken) {
      apiClient.setToken(savedToken);
      useAuthStore.getState().setToken(savedToken);
      // Ensure cookie is set for middleware
      document.cookie = `mockstarket_token=${savedToken}; path=/; max-age=${60 * 60 * 24 * 30}`;
    }
  }, []);

  useEffect(() => {
    if (!token) return;

    wsClient.connect(token);

    // Price updates
    wsClient.on('price_batch', (data) => {
      updatePrices(data);
    });

    // Trade executed
    wsClient.on('trade_executed', (data) => {
      const action = data.side === 'buy' ? 'Bought' : 'Sold';
      toast.info(`${action} ${data.shares} shares of ${data.ticker} at ${formatCurrency(data.price)}`);
      queryClient.invalidateQueries({ queryKey: ['portfolio'] });
      queryClient.invalidateQueries({ queryKey: ['trades'] });
    });

    // Order matched/filled
    wsClient.on('order_matched', (data) => {
      toast.success(`Order filled: ${data.shares} shares of ${data.ticker} at ${formatCurrency(data.price)}`);
      queryClient.invalidateQueries({ queryKey: ['orders'] });
      queryClient.invalidateQueries({ queryKey: ['portfolio'] });
      queryClient.invalidateQueries({ queryKey: ['trades'] });
    });

    // Alert triggered
    wsClient.on('alert_triggered', (data) => {
      toast.info(`Price alert: ${data.ticker} ${data.condition} ${formatCurrency(data.target_price)}`);
      queryClient.invalidateQueries({ queryKey: ['alerts'] });
    });

    // Achievement earned
    wsClient.on('achievement_earned', (data) => {
      toast.success(`Achievement unlocked: ${data.name || 'New achievement!'}`);
      queryClient.invalidateQueries({ queryKey: ['achievements'] });
    });

    // Challenge completed
    wsClient.on('challenge_completed', () => {
      toast.success('Daily challenge completed! Claim your reward.');
      queryClient.invalidateQueries({ queryKey: ['challenge'] });
    });

    return () => wsClient.disconnect();
  }, [token, updatePrices, queryClient]);

  return (
    <div className="flex h-screen bg-[#0D1117]">
      {/* Sidebar */}
      <aside className="hidden md:flex w-64 flex-col border-r border-[#30363D] bg-[#161B22]">
        <div className="p-6">
          <h1 className="text-xl font-bold text-white">Mock Starket</h1>
          <p className="text-xs text-[#6E7681] mt-1">Paper Trading Simulator</p>
        </div>

        <nav className="flex-1 px-3 space-y-1">
          {navItems.map((item) => {
            const Icon = item.icon;
            return (
              <Link
                key={item.href}
                href={item.href}
                className={`flex items-center gap-3 rounded-lg px-3 py-2.5 text-sm font-medium transition-colors ${
                  pathname === item.href
                    ? 'bg-[#21262D] text-white'
                    : 'text-[#8B949E] hover:bg-[#21262D] hover:text-white'
                }`}
              >
                <Icon className="w-4 h-4" />
                {item.label}
              </Link>
            );
          })}
        </nav>

        <div className="px-3 pt-4 border-t border-[#30363D] space-y-1">
          {sidebarExtras.map((item) => {
            const Icon = item.icon;
            return (
              <Link
                key={item.href}
                href={item.href}
                className={`flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors ${
                  pathname === item.href
                    ? 'bg-[#21262D] text-white'
                    : 'text-[#8B949E] hover:bg-[#21262D] hover:text-white'
                }`}
              >
                <Icon className="w-4 h-4" />
                {item.label}
              </Link>
            );
          })}
        </div>

        <div className="p-4 border-t border-[#30363D] space-y-3">
          <Link
            href="/settings"
            className="flex items-center gap-3 rounded-lg px-3 py-2.5 text-sm text-[#8B949E] hover:text-white transition-colors"
          >
            <Settings className="w-4 h-4" />
            Settings
          </Link>
          <div className="px-3">
            <ConnectionStatus />
          </div>
        </div>
      </aside>

      {/* Main content */}
      <main className="flex-1 overflow-y-auto">
        {children}
      </main>

      {/* Mobile bottom nav */}
      <nav className="md:hidden fixed bottom-0 left-0 right-0 flex border-t border-[#30363D] bg-[#161B22]">
        {navItems.map((item) => {
          const Icon = item.icon;
          return (
            <Link
              key={item.href}
              href={item.href}
              className={`flex-1 flex flex-col items-center py-2 text-xs ${
                pathname === item.href ? 'text-[#50E3C2]' : 'text-[#6E7681]'
              }`}
            >
              <Icon className="w-5 h-5 mb-0.5" />
              {item.label}
            </Link>
          );
        })}
      </nav>
    </div>
  );
}
