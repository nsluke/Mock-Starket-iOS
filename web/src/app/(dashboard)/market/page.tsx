'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { apiClient } from '@/lib/api-client';
import { useMarketStore } from '@/stores/market-store';
import { formatCurrency, formatPercent, priceChangeColor, priceChangeBg } from '@/lib/formatters';
import type { Stock, MarketSummary } from '@/types/stock';

export default function MarketPage() {
  const { stocks, summary, isLoading, searchQuery, setStocks, setSummary, setLoading, setSearchQuery, filteredStocks } = useMarketStore();

  useEffect(() => {
    async function loadData() {
      setLoading(true);
      try {
        const [stocksData, summaryData] = await Promise.all([
          apiClient.getStocks(),
          apiClient.getMarketSummary(),
        ]);
        setStocks(stocksData);
        setSummary(summaryData);
      } catch (err) {
        console.error('Failed to load market data:', err);
      } finally {
        setLoading(false);
      }
    }
    loadData();
  }, [setStocks, setSummary, setLoading]);

  const displayed = filteredStocks();

  return (
    <div className="p-6 max-w-6xl mx-auto space-y-6">
      <h1 className="text-2xl font-bold">Market</h1>

      {/* Market Summary */}
      {summary && (
        <div className="rounded-xl bg-[#161B22] border border-[#30363D] p-6">
          <div className="flex items-baseline justify-between">
            <div>
              <p className="text-sm text-[#8B949E] mb-1">Market Index</p>
              <p className="text-3xl font-bold">{formatCurrency(summary.index_value)}</p>
            </div>
            <div className={`text-lg font-semibold ${priceChangeColor(summary.index_change_pct)}`}>
              {formatPercent(summary.index_change_pct)}
            </div>
          </div>
          <div className="mt-3 flex gap-4 text-sm text-[#8B949E]">
            <span className="text-emerald-400">{summary.gainers} gainers</span>
            <span className="text-red-400">{summary.losers} losers</span>
            <span>{summary.total_stocks} stocks</span>
          </div>
        </div>
      )}

      {/* Search */}
      <input
        type="text"
        placeholder="Search stocks..."
        value={searchQuery}
        onChange={(e) => setSearchQuery(e.target.value)}
        className="w-full rounded-lg bg-[#161B22] border border-[#30363D] px-4 py-3 text-sm text-white placeholder-[#6E7681] focus:outline-none focus:border-[#50E3C2]"
      />

      {/* Stock Table */}
      <div className="rounded-xl bg-[#161B22] border border-[#30363D] overflow-hidden">
        <table className="w-full">
          <thead>
            <tr className="border-b border-[#30363D] text-xs text-[#8B949E] uppercase">
              <th className="text-left px-4 py-3">Stock</th>
              <th className="text-left px-4 py-3 hidden sm:table-cell">Sector</th>
              <th className="text-right px-4 py-3">Price</th>
              <th className="text-right px-4 py-3">Change</th>
            </tr>
          </thead>
          <tbody>
            {displayed.map((stock) => {
              const change = parseFloat(stock.current_price) - parseFloat(stock.day_open);
              const changePct = parseFloat(stock.day_open) !== 0 ? (change / parseFloat(stock.day_open)) * 100 : 0;

              return (
                <tr
                  key={stock.ticker}
                  className="border-b border-[#21262D] hover:bg-[#21262D] transition-colors cursor-pointer"
                >
                  <td className="px-4 py-3">
                    <Link href={`/stock/${stock.ticker}`} className="flex items-center gap-3">
                      <span className="rounded bg-[#50E3C2]/10 px-2 py-1 text-xs font-mono font-bold text-[#50E3C2]">
                        {stock.ticker}
                      </span>
                      <span className="text-sm font-medium text-white">{stock.name}</span>
                    </Link>
                  </td>
                  <td className="px-4 py-3 text-sm text-[#8B949E] hidden sm:table-cell">
                    {stock.sector}
                  </td>
                  <td className="px-4 py-3 text-right text-sm font-semibold text-white">
                    {formatCurrency(stock.current_price)}
                  </td>
                  <td className="px-4 py-3 text-right">
                    <span className={`inline-block rounded-md px-2 py-1 text-xs font-semibold ${priceChangeBg(changePct)}`}>
                      {formatPercent(changePct)}
                    </span>
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>

        {isLoading && (
          <div className="p-8 text-center text-[#8B949E]">Loading stocks...</div>
        )}
      </div>
    </div>
  );
}
