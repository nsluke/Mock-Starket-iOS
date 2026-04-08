'use client';

import { useState } from 'react';
import Link from 'next/link';
import { PageTransition } from '@/components/ui/PageTransition';
import { useWatchlist, useAddToWatchlist, useRemoveFromWatchlist } from '@/hooks/use-watchlist';
import { useMarketStore } from '@/stores/market-store';
import { formatCurrency, formatPercent, priceChangeBg } from '@/lib/formatters';

export default function WatchlistPage() {
  const { data: watchedTickers = [], isLoading } = useWatchlist();
  const { stocks } = useMarketStore();
  const addToWatchlist = useAddToWatchlist();
  const removeFromWatchlist = useRemoveFromWatchlist();

  const [showAdd, setShowAdd] = useState(false);
  const [search, setSearch] = useState('');

  const allStocks = stocks || [];
  const watchedStocks = watchedTickers
    .map((ticker) => allStocks.find((s) => s.ticker === ticker))
    .filter(Boolean);

  const searchResults = search.length > 0
    ? allStocks.filter(
        (s) =>
          !watchedTickers.includes(s.ticker) &&
          (s.ticker.toLowerCase().includes(search.toLowerCase()) ||
            s.name.toLowerCase().includes(search.toLowerCase()))
      ).slice(0, 8)
    : [];

  return (
    <PageTransition>
    <div className="p-6 max-w-4xl mx-auto space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">Watchlist</h1>
        <button
          onClick={() => setShowAdd(!showAdd)}
          className="px-4 py-2 rounded-lg bg-[#50E3C2] text-[#0D1117] text-sm font-semibold hover:bg-[#3BC4A7] transition-colors"
        >
          {showAdd ? 'Done' : 'Add Stock'}
        </button>
      </div>

      {/* Add to watchlist search */}
      {showAdd && (
        <div className="rounded-xl bg-[#161B22] border border-[#30363D] p-4 space-y-3">
          <input
            type="text"
            placeholder="Search stocks to add..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="w-full rounded-lg bg-[#0D1117] border border-[#30363D] px-4 py-2.5 text-sm text-white placeholder-[#6E7681] focus:outline-none focus:border-[#50E3C2]"
            autoFocus
          />
          {searchResults.length > 0 && (
            <div className="space-y-1">
              {searchResults.map((stock) => (
                <div
                  key={stock.ticker}
                  className="flex items-center justify-between px-3 py-2 rounded-lg hover:bg-[#21262D] transition-colors"
                >
                  <div className="flex items-center gap-3">
                    <span className="rounded bg-[#50E3C2]/10 px-2 py-0.5 text-xs font-mono font-bold text-[#50E3C2]">
                      {stock.ticker}
                    </span>
                    <span className="text-sm text-[#8B949E]">{stock.name}</span>
                  </div>
                  <button
                    onClick={() => {
                      addToWatchlist.mutate(stock.ticker);
                      setSearch('');
                    }}
                    className="text-xs text-[#50E3C2] hover:text-[#3BC4A7] font-medium transition-colors"
                  >
                    + Add
                  </button>
                </div>
              ))}
            </div>
          )}
          {search.length > 0 && searchResults.length === 0 && (
            <p className="text-sm text-[#6E7681] px-3">No stocks found</p>
          )}
        </div>
      )}

      {/* Watchlist table */}
      <div className="rounded-xl bg-[#161B22] border border-[#30363D] overflow-hidden">
        {isLoading ? (
          <div className="p-8 text-center text-[#8B949E]">Loading watchlist...</div>
        ) : watchedStocks.length === 0 ? (
          <div className="p-8 text-center text-[#8B949E]">
            <p>Your watchlist is empty.</p>
            <p className="mt-1 text-sm">Add stocks you want to track.</p>
          </div>
        ) : (
          <table className="w-full">
            <thead>
              <tr className="border-b border-[#30363D] text-xs text-[#8B949E] uppercase">
                <th className="text-left px-4 py-3">Stock</th>
                <th className="text-right px-4 py-3">Price</th>
                <th className="text-right px-4 py-3">Change</th>
                <th className="text-right px-4 py-3 w-20">Action</th>
              </tr>
            </thead>
            <tbody>
              {watchedStocks.map((stock) => {
                if (!stock) return null;
                const change = parseFloat(stock.current_price) - parseFloat(stock.day_open);
                const changePct = parseFloat(stock.day_open) !== 0 ? (change / parseFloat(stock.day_open)) * 100 : 0;

                return (
                  <tr key={stock.ticker} className="border-b border-[#21262D] hover:bg-[#21262D] transition-colors">
                    <td className="px-4 py-3">
                      <Link href={`/stock/${stock.ticker}`} className="flex items-center gap-3">
                        <span className="rounded bg-[#50E3C2]/10 px-2 py-1 text-xs font-mono font-bold text-[#50E3C2]">
                          {stock.ticker}
                        </span>
                        <span className="text-sm font-medium text-white">{stock.name}</span>
                      </Link>
                    </td>
                    <td className="px-4 py-3 text-right text-sm font-semibold text-white">
                      {formatCurrency(stock.current_price)}
                    </td>
                    <td className="px-4 py-3 text-right">
                      <span className={`inline-block rounded-md px-2 py-1 text-xs font-semibold ${priceChangeBg(changePct)}`}>
                        {formatPercent(changePct)}
                      </span>
                    </td>
                    <td className="px-4 py-3 text-right">
                      <button
                        onClick={() => removeFromWatchlist.mutate(stock.ticker)}
                        className="text-xs text-red-400 hover:text-red-300 transition-colors"
                      >
                        Remove
                      </button>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        )}
      </div>
    </div>
    </PageTransition>
  );
}
