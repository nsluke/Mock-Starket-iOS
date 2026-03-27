'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { apiClient } from '@/lib/api-client';
import { formatCurrency, formatPercent, priceChangeColor, priceChangeBg } from '@/lib/formatters';
import type { Trade } from '@/types/portfolio';

interface Position {
  id: string;
  ticker: string;
  shares: number;
  avg_cost: string;
  current_price: string;
  market_value: string;
  pnl: string;
  pnl_pct: string;
}

interface PortfolioData {
  portfolio: { id: string; cash: string; net_worth: string };
  positions: Position[];
  net_worth: string;
  invested: string;
}

export default function PortfolioPage() {
  const [data, setData] = useState<PortfolioData | null>(null);
  const [trades, setTrades] = useState<Trade[]>([]);
  const [loading, setLoading] = useState(true);
  const [tab, setTab] = useState<'holdings' | 'history'>('holdings');

  useEffect(() => {
    async function load() {
      setLoading(true);
      try {
        const [portfolio, tradeHistory] = await Promise.all([
          apiClient.getPortfolio(),
          apiClient.getTradeHistory(20, 0),
        ]);
        setData(portfolio);
        setTrades(tradeHistory || []);
      } catch (err) {
        console.error('Failed to load portfolio:', err);
      } finally {
        setLoading(false);
      }
    }
    load();
  }, []);

  if (loading) {
    return <div className="p-6 text-center text-[#8B949E]">Loading portfolio...</div>;
  }

  if (!data) {
    return <div className="p-6 text-center text-[#8B949E]">Failed to load portfolio</div>;
  }

  const netWorth = parseFloat(data.net_worth);
  const cash = parseFloat(data.portfolio.cash);
  const invested = parseFloat(data.invested);
  const totalPnl = netWorth - 100000; // starting cash
  const totalPnlPct = (totalPnl / 100000) * 100;

  return (
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
            {data.positions.length} position{data.positions.length !== 1 ? 's' : ''}
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
          {data.positions.length === 0 ? (
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
                {data.positions.map((pos) => (
                  <tr
                    key={pos.id}
                    className="border-b border-[#21262D] hover:bg-[#21262D] transition-colors"
                  >
                    <td className="px-4 py-3">
                      <Link
                        href={`/stock/${pos.ticker}`}
                        className="flex items-center gap-2"
                      >
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
  );
}
