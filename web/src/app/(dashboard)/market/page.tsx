'use client';

import { useState, useMemo } from 'react';
import Link from 'next/link';
import { useStocks, useMarketSummary } from '@/hooks/use-stocks';
import { useMarketStore } from '@/stores/market-store';
import { useDebounce } from '@/hooks/use-debounce';
import { useSort } from '@/hooks/use-sort';
import { PageTransition } from '@/components/ui/PageTransition';
import { TableSkeleton } from '@/components/ui/TableSkeleton';
import { CardSkeleton } from '@/components/ui/CardSkeleton';
import { formatCurrency, formatPercent, priceChangeColor, priceChangeBg } from '@/lib/formatters';
import type { Stock } from '@/types/stock';

const assetFilters = [
  { value: 'all', label: 'All' },
  { value: 'stock', label: 'Stocks' },
  { value: 'etf', label: 'ETFs' },
  { value: 'crypto', label: 'Crypto' },
  { value: 'commodity', label: 'Commodities' },
];

export default function MarketPage() {
  const { setSearchQuery, filteredStocks } = useMarketStore();
  const [localSearch, setLocalSearch] = useState('');
  const [assetFilter, setAssetFilter] = useState('all');
  const debouncedSearch = useDebounce(localSearch, 300);

  const { isLoading: stocksLoading } = useStocks();
  const { data: summary } = useMarketSummary();

  // Sync debounced value to store
  useMemo(() => {
    setSearchQuery(debouncedSearch);
  }, [debouncedSearch, setSearchQuery]);

  const filtered = useMemo(() => {
    const searched = filteredStocks();
    if (assetFilter === 'all') return searched;
    return searched.filter((s: Stock) => s.asset_type === assetFilter);
  }, [filteredStocks, assetFilter]);

  const { sorted: displayed, sortKey, sortDirection, onSort } = useSort(filtered, 'ticker', 'asc');

  return (
    <PageTransition>
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
        value={localSearch}
        onChange={(e) => setLocalSearch(e.target.value)}
        className="w-full rounded-lg bg-[#161B22] border border-[#30363D] px-4 py-3 text-sm text-white placeholder-[#6E7681] focus:outline-none focus:border-[#50E3C2]"
      />

      {/* Asset Type Filter */}
      <div className="flex gap-1 overflow-x-auto">
        {assetFilters.map((f) => (
          <button
            key={f.value}
            onClick={() => setAssetFilter(f.value)}
            className={`px-4 py-2 rounded-lg text-sm font-medium whitespace-nowrap transition-colors ${
              assetFilter === f.value
                ? 'bg-[#21262D] text-white'
                : 'text-[#8B949E] hover:text-white'
            }`}
          >
            {f.label}
          </button>
        ))}
      </div>

      {/* Stock Table */}
      <div className="rounded-xl bg-[#161B22] border border-[#30363D] overflow-hidden">
        <table className="w-full">
          <thead>
            <tr className="border-b border-[#30363D] text-xs text-[#8B949E] uppercase">
              <th className="text-left px-4 py-3 cursor-pointer hover:text-white select-none" onClick={() => onSort('ticker')}>
                Stock {sortKey === 'ticker' && (sortDirection === 'asc' ? '↑' : '↓')}
              </th>
              <th className="text-left px-4 py-3 hidden sm:table-cell">Type</th>
              <th className="text-right px-4 py-3 cursor-pointer hover:text-white select-none" onClick={() => onSort('current_price')}>
                Price {sortKey === 'current_price' && (sortDirection === 'asc' ? '↑' : '↓')}
              </th>
              <th className="text-right px-4 py-3 cursor-pointer hover:text-white select-none" onClick={() => onSort('name')}>
                Change {sortKey === 'name' && (sortDirection === 'asc' ? '↑' : '↓')}
              </th>
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
                  <td className="px-4 py-3 hidden sm:table-cell">
                    <span className={`text-xs font-medium px-2 py-0.5 rounded ${
                      stock.asset_type === 'crypto' ? 'bg-orange-400/10 text-orange-400' :
                      stock.asset_type === 'commodity' ? 'bg-yellow-400/10 text-yellow-400' :
                      stock.asset_type === 'etf' ? 'bg-purple-400/10 text-purple-400' :
                      'bg-[#8B949E]/10 text-[#8B949E]'
                    }`}>
                      {stock.asset_type?.toUpperCase() || 'STOCK'}
                    </span>
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

        {stocksLoading && <TableSkeleton rows={8} columns={4} />}
      </div>
    </div>
    </PageTransition>
  );
}
