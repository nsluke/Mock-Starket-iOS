'use client';

import { useState } from 'react';
import { PageTransition } from '@/components/ui/PageTransition';
import { TableSkeleton } from '@/components/ui/TableSkeleton';
import { useLeaderboard } from '@/hooks/use-leaderboard';
import { formatCurrency, formatPercent, priceChangeColor } from '@/lib/formatters';

const periods = [
  { value: 'alltime', label: 'All Time' },
  { value: 'weekly', label: 'Weekly' },
  { value: 'daily', label: 'Daily' },
];

export default function LeaderboardPage() {
  const [period, setPeriod] = useState('alltime');
  const { data: entries = [], isLoading } = useLeaderboard(period);

  return (
    <PageTransition>
    <div className="p-6 max-w-4xl mx-auto space-y-6">
      <h1 className="text-2xl font-bold">Leaderboard</h1>

      {/* Period Toggle */}
      <div className="flex gap-1 bg-[#161B22] rounded-lg p-1 border border-[#30363D] w-fit">
        {periods.map((p) => (
          <button
            key={p.value}
            onClick={() => setPeriod(p.value)}
            className={`px-4 py-2 rounded-md text-sm font-medium transition-colors ${
              period === p.value
                ? 'bg-[#21262D] text-white'
                : 'text-[#8B949E] hover:text-white'
            }`}
          >
            {p.label}
          </button>
        ))}
      </div>

      {/* Leaderboard Table */}
      <div className="rounded-xl bg-[#161B22] border border-[#30363D] overflow-hidden">
        {isLoading ? (
          <TableSkeleton rows={6} columns={4} />
        ) : entries.length === 0 ? (
          <div className="p-8 text-center text-[#8B949E]">
            No rankings yet. Start trading to appear on the leaderboard!
          </div>
        ) : (
          <table className="w-full">
            <thead>
              <tr className="border-b border-[#30363D] text-xs text-[#8B949E] uppercase">
                <th className="text-left px-4 py-3 w-16">Rank</th>
                <th className="text-left px-4 py-3">Trader</th>
                <th className="text-right px-4 py-3">Net Worth</th>
                <th className="text-right px-4 py-3">Return</th>
              </tr>
            </thead>
            <tbody>
              {entries.map((entry) => (
                <tr
                  key={entry.id}
                  className="border-b border-[#21262D] hover:bg-[#21262D] transition-colors"
                >
                  <td className="px-4 py-3">
                    <span
                      className={`inline-flex items-center justify-center w-8 h-8 rounded-full text-sm font-bold ${
                        entry.rank === 1
                          ? 'bg-yellow-500/20 text-yellow-400'
                          : entry.rank === 2
                          ? 'bg-gray-400/20 text-gray-300'
                          : entry.rank === 3
                          ? 'bg-amber-700/20 text-amber-500'
                          : 'text-[#8B949E]'
                      }`}
                    >
                      {entry.rank}
                    </span>
                  </td>
                  <td className="px-4 py-3">
                    <span className="font-medium">{entry.display_name}</span>
                  </td>
                  <td className="px-4 py-3 text-right font-semibold">
                    {formatCurrency(entry.net_worth)}
                  </td>
                  <td className="px-4 py-3 text-right">
                    <span className={`font-semibold ${priceChangeColor(entry.total_return)}`}>
                      {formatPercent(entry.total_return)}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
    </PageTransition>
  );
}
