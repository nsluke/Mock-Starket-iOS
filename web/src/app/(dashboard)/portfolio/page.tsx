'use client';

import { useState, useMemo } from 'react';
import Link from 'next/link';
import dynamic from 'next/dynamic';
import { usePortfolio, usePortfolioHistory, useTradeHistory } from '@/hooks/use-portfolio';
import { PageTransition } from '@/components/ui/PageTransition';
import { CardSkeleton } from '@/components/ui/CardSkeleton';
import { PieChart, COLORS } from '@/components/charts/PieChart';
import { OptionsPositionsTable } from '@/components/options/OptionsPositionsTable';
import { formatCurrency, formatPercent, priceChangeColor } from '@/lib/formatters';

const PortfolioChart = dynamic(() => import('@/components/charts/PortfolioChart'), { ssr: false });

export default function PortfolioPage() {
  const { data, isLoading } = usePortfolio();
  const { data: trades = [] } = useTradeHistory(20, 0);
  const { data: portfolioHistory = [] } = usePortfolioHistory(100);
  const [tab, setTab] = useState<'holdings' | 'options' | 'history'>('holdings');

  if (isLoading) {
    return (
      <div className="p-6 max-w-6xl mx-auto space-y-6">
        <h1 className="text-2xl font-bold">Portfolio</h1>
        <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
          <CardSkeleton /><CardSkeleton /><CardSkeleton /><CardSkeleton />
        </div>
      </div>
    );
  }

  if (!data) {
    return <div className="p-6 text-center text-[#8B949E]">Failed to load portfolio</div>;
  }

  const netWorth = parseFloat(data.net_worth);
  const positions = data.positions || [];
  const cash = parseFloat(data.portfolio?.cash || '0');
  const invested = parseFloat(data.invested || '0');
  const totalPnl = netWorth - 100000;
  const totalPnlPct = (totalPnl / 100000) * 100;

  return (
    <PageTransition>
    <div className="p-6 max-w-6xl mx-auto space-y-6">
      <h1 className="text-2xl font-bold">Portfolio</h1>

      {/* Summary Cards */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
        <div className="rounded-xl bg-[#161B22] border border-[#30363D] p-4">
          <p className="text-xs text-[#6E7681] mb-1">Net Worth</p>
          <p className="text-xl font-bold">{formatCurrency(netWorth)}</p>
          <p className={`text-xs mt-1 ${priceChangeColor(totalPnl)}`}>
            {formatPercent(totalPnlPct)} all time
          </p>
        </div>
        <div className="rounded-xl bg-[#161B22] border border-[#30363D] p-4">
          <p className="text-xs text-[#6E7681] mb-1">Cash</p>
          <p className="text-xl font-bold">{formatCurrency(cash)}</p>
          <p className="text-xs text-[#6E7681] mt-1">
            {((cash / netWorth) * 100).toFixed(1)}% of portfolio
          </p>
        </div>
        <div className="rounded-xl bg-[#161B22] border border-[#30363D] p-4">
          <p className="text-xs text-[#6E7681] mb-1">Invested</p>
          <p className="text-xl font-bold">{formatCurrency(invested)}</p>
          <p className="text-xs text-[#6E7681] mt-1">
            {positions.length} position{positions.length !== 1 ? 's' : ''}
          </p>
        </div>
        <div className="rounded-xl bg-[#161B22] border border-[#30363D] p-4">
          <p className="text-xs text-[#6E7681] mb-1">Total P&L</p>
          <p className={`text-xl font-bold ${priceChangeColor(totalPnl)}`}>
            {totalPnl >= 0 ? '+' : ''}{formatCurrency(totalPnl)}
          </p>
          <p className={`text-xs mt-1 ${priceChangeColor(totalPnl)}`}>
            {formatPercent(totalPnlPct)}
          </p>
        </div>
      </div>

      {/* Portfolio Performance Chart */}
      {portfolioHistory.length > 1 && (
        <div className="rounded-xl bg-[#161B22] border border-[#30363D] p-6">
          <h2 className="font-semibold mb-4">Performance</h2>
          <PortfolioChart
            data={portfolioHistory.map((h: any) => ({
              time: h.recorded_at,
              value: parseFloat(h.net_worth),
            }))}
          />
        </div>
      )}

      {/* Portfolio Breakdown */}
      {positions.length > 0 && (
        <div className="rounded-xl bg-[#161B22] border border-[#30363D] p-6">
          <h2 className="font-semibold mb-4">Portfolio Breakdown</h2>
          <PieChart
            data={[
              ...positions.map((pos, i) => ({
                label: pos.ticker,
                value: parseFloat(pos.market_value),
                color: COLORS[i % COLORS.length],
              })),
              { label: 'Cash', value: cash, color: '#6E7681' },
            ]}
            size={160}
          />
        </div>
      )}

      {/* Tab Toggle */}
      <div className="flex border-b border-[#30363D]">
        <button
          onClick={() => setTab('holdings')}
          className={`px-4 py-2.5 text-sm font-medium border-b-2 transition-colors ${
            tab === 'holdings'
              ? 'border-[#50E3C2] text-white'
              : 'border-transparent text-[#8B949E] hover:text-white'
          }`}
        >
          Holdings
        </button>
        <button
          onClick={() => setTab('options')}
          className={`px-4 py-2.5 text-sm font-medium border-b-2 transition-colors ${
            tab === 'options'
              ? 'border-[#50E3C2] text-white'
              : 'border-transparent text-[#8B949E] hover:text-white'
          }`}
        >
          Options
        </button>
        <button
          onClick={() => setTab('history')}
          className={`px-4 py-2.5 text-sm font-medium border-b-2 transition-colors ${
            tab === 'history'
              ? 'border-[#50E3C2] text-white'
              : 'border-transparent text-[#8B949E] hover:text-white'
          }`}
        >
          Trade History
        </button>
      </div>

      {/* Holdings Table */}
      {tab === 'holdings' && (
        <div className="rounded-xl bg-[#161B22] border border-[#30363D] overflow-hidden">
          {positions.length === 0 ? (
            <div className="p-8 text-center text-[#8B949E]">
              <p>No positions yet.</p>
              <Link href="/market" className="text-[#50E3C2] hover:underline mt-1 inline-block">
                Browse the market to start trading
              </Link>
            </div>
          ) : (
            <table className="w-full">
              <thead>
                <tr className="border-b border-[#30363D] text-xs text-[#8B949E] uppercase">
                  <th className="text-left px-4 py-3">Stock</th>
                  <th className="text-right px-4 py-3 hidden sm:table-cell">Shares</th>
                  <th className="text-right px-4 py-3 hidden sm:table-cell">Avg Cost</th>
                  <th className="text-right px-4 py-3">Market Value</th>
                  <th className="text-right px-4 py-3">P&L</th>
                </tr>
              </thead>
              <tbody>
                {positions.map((pos) => (
                  <tr
                    key={pos.id}
                    className="border-b border-[#21262D] hover:bg-[#21262D] transition-colors"
                  >
                    <td className="px-4 py-3">
                      <Link href={`/stock/${pos.ticker}`} className="flex items-center gap-2">
                        <span className="rounded bg-[#50E3C2]/10 px-2 py-0.5 text-xs font-mono font-bold text-[#50E3C2]">
                          {pos.ticker}
                        </span>
                        <span className="text-sm text-[#8B949E] sm:hidden">{pos.shares} shares</span>
                      </Link>
                    </td>
                    <td className="px-4 py-3 text-right text-sm hidden sm:table-cell">{pos.shares}</td>
                    <td className="px-4 py-3 text-right text-sm text-[#8B949E] hidden sm:table-cell">
                      {formatCurrency(pos.avg_cost)}
                    </td>
                    <td className="px-4 py-3 text-right text-sm font-semibold">
                      {formatCurrency(pos.market_value)}
                    </td>
                    <td className="px-4 py-3 text-right">
                      <div className={`text-sm font-semibold ${priceChangeColor(pos.pnl)}`}>
                        {parseFloat(pos.pnl) >= 0 ? '+' : ''}{formatCurrency(pos.pnl)}
                      </div>
                      <div className={`text-xs ${priceChangeColor(pos.pnl_pct)}`}>
                        {formatPercent(pos.pnl_pct)}
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      )}

      {/* Options Positions */}
      {tab === 'options' && (
        <div className="rounded-xl bg-[#161B22] border border-[#30363D] overflow-hidden">
          <OptionsPositionsTable />
        </div>
      )}

      {/* Trade History */}
      {tab === 'history' && (
        <div className="rounded-xl bg-[#161B22] border border-[#30363D] overflow-hidden">
          {trades.length === 0 ? (
            <div className="p-8 text-center text-[#8B949E]">No trades yet.</div>
          ) : (
            <table className="w-full">
              <thead>
                <tr className="border-b border-[#30363D] text-xs text-[#8B949E] uppercase">
                  <th className="text-left px-4 py-3">Stock</th>
                  <th className="text-left px-4 py-3">Side</th>
                  <th className="text-right px-4 py-3 hidden sm:table-cell">Shares</th>
                  <th className="text-right px-4 py-3">Price</th>
                  <th className="text-right px-4 py-3">Total</th>
                  <th className="text-right px-4 py-3 hidden sm:table-cell">Time</th>
                </tr>
              </thead>
              <tbody>
                {trades.map((trade) => (
                  <tr key={trade.id} className="border-b border-[#21262D]">
                    <td className="px-4 py-3">
                      <span className="rounded bg-[#50E3C2]/10 px-2 py-0.5 text-xs font-mono font-bold text-[#50E3C2]">
                        {trade.ticker}
                      </span>
                    </td>
                    <td className="px-4 py-3">
                      <span
                        className={`text-xs font-semibold uppercase ${
                          trade.side === 'buy' ? 'text-emerald-400' : 'text-red-400'
                        }`}
                      >
                        {trade.side}
                      </span>
                    </td>
                    <td className="px-4 py-3 text-right text-sm hidden sm:table-cell">{trade.shares}</td>
                    <td className="px-4 py-3 text-right text-sm">{formatCurrency(trade.price)}</td>
                    <td className="px-4 py-3 text-right text-sm font-semibold">{formatCurrency(trade.total)}</td>
                    <td className="px-4 py-3 text-right text-xs text-[#8B949E] hidden sm:table-cell">
                      {new Date(trade.created_at).toLocaleString()}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      )}
    </div>
    </PageTransition>
  );
}
