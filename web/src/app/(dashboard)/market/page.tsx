'use client';

import { useState, useMemo, useEffect } from 'react';
import Link from 'next/link';
import { useStocks, useMarketSummary, useMarketStatus } from '@/hooks/use-stocks';
import { useMarketStore } from '@/stores/market-store';
import { useDebounce } from '@/hooks/use-debounce';
import { useSort } from '@/hooks/use-sort';
import { PageTransition } from '@/components/ui/PageTransition';
import { TableSkeleton } from '@/components/ui/TableSkeleton';
import { PieChart, COLORS } from '@/components/charts/PieChart';
import { formatCurrency, formatPercent, priceChangeColor, priceChangeBg } from '@/lib/formatters';
import type { Stock } from '@/types/stock';
import { displayTicker } from '@/types/stock';

const assetFilters = [
  { value: 'all', label: 'All' },
  { value: 'stock', label: 'Stocks' },
  { value: 'etf', label: 'ETFs' },
  { value: 'crypto', label: 'Crypto' },
];

export default function MarketPage() {
  const stocks = useMarketStore((s) => s.stocks);
  const searchQuery = useMarketStore((s) => s.searchQuery);
  const setSearchQuery = useMarketStore((s) => s.setSearchQuery);
  const [localSearch, setLocalSearch] = useState('');
  const [assetFilter, setAssetFilter] = useState('all');
  const [sectorFilter, setSectorFilter] = useState('all');
  const debouncedSearch = useDebounce(localSearch, 300);

  const { isLoading: stocksLoading } = useStocks();
  const { data: summary } = useMarketSummary();
  const { data: marketStatus } = useMarketStatus();

  // Sync debounced value to store
  useEffect(() => {
    setSearchQuery(debouncedSearch);
  }, [debouncedSearch, setSearchQuery]);

  // Available sectors based on current asset type filter
  const availableSectors = useMemo(() => {
    const pool = assetFilter === 'all' ? stocks : stocks.filter((s) => s.asset_type === assetFilter);
    const unique = [...new Set(pool.map((s) => s.sector))].sort();
    return unique;
  }, [stocks, assetFilter]);

  // Reset sector when asset type changes
  useEffect(() => {
    setSectorFilter('all');
  }, [assetFilter]);

  const filtered = useMemo(() => {
    let result = stocks;
    if (searchQuery) {
      const q = searchQuery.toLowerCase();
      result = result.filter((s) => s.ticker.toLowerCase().includes(q) || s.name.toLowerCase().includes(q));
    }
    if (assetFilter !== 'all') {
      result = result.filter((s: Stock) => s.asset_type === assetFilter);
    }
    if (sectorFilter !== 'all') {
      result = result.filter((s: Stock) => s.sector === sectorFilter);
    }
    return result.map((s) => {
      const change = parseFloat(s.current_price) - parseFloat(s.day_open);
      const changePct = parseFloat(s.day_open) !== 0 ? (change / parseFloat(s.day_open)) * 100 : 0;
      return { ...s, _changePct: changePct };
    });
  }, [stocks, searchQuery, assetFilter, sectorFilter]);

  const { sorted: displayed, sortKey, sortDirection, onSort } = useSort(filtered, 'ticker', 'asc');

  // Sector breakdown for pie chart (by count of stocks, excluding ETFs)
  const sectorData = useMemo(() => {
    const counts: Record<string, number> = {};
    stocks.forEach((s) => {
      if (s.asset_type === 'etf') return;
      const key = s.sector || s.asset_type;
      counts[key] = (counts[key] || 0) + 1;
    });
    return Object.entries(counts)
      .map(([label, value], i) => ({ label, value, color: COLORS[i % COLORS.length] }))
      .sort((a, b) => b.value - a.value);
  }, [stocks]);

  return (
    <PageTransition>
    <div className="p-6 max-w-6xl mx-auto space-y-6">
      <div className="flex items-center gap-3">
        <h1 className="text-2xl font-bold">Market</h1>
        {marketStatus && (
          <span className={`inline-flex items-center gap-1.5 rounded-full px-3 py-1 text-xs font-medium ${
            marketStatus.is_open
              ? 'bg-emerald-400/10 text-emerald-400'
              : 'bg-[#8B949E]/10 text-[#8B949E]'
          }`}>
            <span className={`w-1.5 h-1.5 rounded-full ${marketStatus.is_open ? 'bg-emerald-400 animate-pulse' : 'bg-[#6E7681]'}`} />
            {marketStatus.is_open ? 'Market Open' : marketStatus.session === 'pre_market' ? 'Pre-Market' : marketStatus.session === 'after_hours' ? 'After Hours' : 'Market Closed'}
          </span>
        )}
      </div>

      {/* Market Summary */}
      {summary && (
        <div className="rounded-xl bg-[#161B22] border border-[#30363D] p-6">
          <div className="flex flex-col md:flex-row md:items-center gap-6">
            <div className="flex-1">
              <div className="flex items-baseline justify-between mb-3">
                <div>
                  <p className="text-sm text-[#8B949E] mb-1">Market Index</p>
                  <p className="text-3xl font-bold">{formatCurrency(summary.index_value)}</p>
                </div>
                <div className={`text-lg font-semibold ${priceChangeColor(summary.index_change_pct)}`}>
                  {formatPercent(summary.index_change_pct)}
                </div>
              </div>
              <div className="flex gap-4 text-sm text-[#8B949E]">
                <span className="text-emerald-400">{summary.gainers} gainers</span>
                <span className="text-red-400">{summary.losers} losers</span>
                <span>{summary.total_stocks} assets</span>
              </div>
            </div>
            {sectorData.length > 0 && (
              <div>
                <p className="text-xs text-[#6E7681] mb-2">Composition by Sector</p>
                <PieChart data={sectorData} size={130} />
              </div>
            )}
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

      {/* Sector Filter */}
      {availableSectors.length > 1 && (
        <div className="flex gap-1.5 overflow-x-auto pb-1">
          <button
            onClick={() => setSectorFilter('all')}
            className={`px-3 py-1.5 rounded-full text-xs font-medium whitespace-nowrap transition-colors ${
              sectorFilter === 'all'
                ? 'bg-[#50E3C2]/15 text-[#50E3C2]'
                : 'bg-[#21262D] text-[#8B949E] hover:text-white'
            }`}
          >
            All Sectors
          </button>
          {availableSectors.map((sector) => (
            <button
              key={sector}
              onClick={() => setSectorFilter(sector)}
              className={`px-3 py-1.5 rounded-full text-xs font-medium whitespace-nowrap transition-colors ${
                sectorFilter === sector
                  ? 'bg-[#50E3C2]/15 text-[#50E3C2]'
                  : 'bg-[#21262D] text-[#8B949E] hover:text-white'
              }`}
            >
              {sector}
            </button>
          ))}
        </div>
      )}

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
              <th className="text-right px-4 py-3 cursor-pointer hover:text-white select-none" onClick={() => onSort('_changePct' as any)}>
                Change {sortKey === ('_changePct' as any) && (sortDirection === 'asc' ? '↑' : '↓')}
              </th>
            </tr>
          </thead>
          <tbody>
            {displayed.map((stock) => {
              const changePct = stock._changePct;

              return (
                <tr
                  key={stock.ticker}
                  className="border-b border-[#21262D] hover:bg-[#21262D] transition-colors cursor-pointer"
                >
                  <td className="px-4 py-3">
                    <Link href={`/stock/${stock.ticker}`} className="flex items-center gap-3">
                      {stock.logo_url ? (
                        <img src={stock.logo_url} alt="" className="w-8 h-8 rounded-full bg-[#21262D] object-contain" />
                      ) : (
                        <span className="w-8 h-8 rounded-full bg-[#50E3C2]/10 flex items-center justify-center text-xs font-bold text-[#50E3C2]">
                          {displayTicker(stock.ticker).slice(0, 2)}
                        </span>
                      )}
                      <div className="min-w-0">
                        <span className="text-xs font-mono font-bold text-[#50E3C2]">
                          {displayTicker(stock.ticker)}
                        </span>
                        <p className="text-sm font-medium text-white truncate">{stock.name}</p>
                      </div>
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
